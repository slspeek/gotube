package rest

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"github.com/slspeek/flowgo"
	"github.com/slspeek/goblob"
	"io"
	"labix.org/v2/mgo"
	"math"
	"math/rand"
	"mime/multipart"
	"net/http"
	"testing"
	"time"
)

const cid = 42
const hex = "ABCD"

var f = flow.Flow{"flow_id", "test.txt", 2, 1024}

func check(t *testing.T, err error) {
	if err != nil {
		t.Fatal(err)
	}
}

func blobService() (bs *goblob.BlobService) {
	s, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	bs = goblob.NewBlobService(s, "test", "testfs")
	return
}

func makeRequest(url string, body io.Reader, f flow.Flow, chunkNumber int) *http.Request {
	buf := new(bytes.Buffer)
	writer := multipart.NewWriter(buf)
	writer.WriteField("flowChunkSize", fmt.Sprintf("%v", f.ChunkSize))
	writer.WriteField("flowChunkNumber", fmt.Sprintf("%v", chunkNumber))
	writer.WriteField("flowTotalChunks", fmt.Sprintf("%v", f.TotalChunks))
	writer.WriteField("flowIdentifier", f.Identifier)
	writer.WriteField("flowFilename", f.Filename)
	fileWriter, _ := writer.CreateFormFile("file", f.Filename)
	io.Copy(fileWriter, body)
	writer.Close()
	req, _ := http.NewRequest("POST", url, buf)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	return req
}

func testBytes(n int) (buf *bytes.Buffer, md5sum string) {
	buf = new(bytes.Buffer)
	h := md5.New()
	for i := 0; i < n; i++ {
		c := byte(rand.Int())
		buf.WriteByte(c)
		h.Write([]byte{c})
	}
	md5sum = fmt.Sprintf("%x", h.Sum(nil))
	return
}

func PrepareRequests(url string, fn string, fs int64, chunkSize int64) (f flow.Flow, reqs []*http.Request, md5sum string) {
	identifier := fmt.Sprintf("%d-%s", fs, fn)
	totalChunks := int(math.Ceil(float64(fs) / float64(chunkSize)))
	f = flow.Flow{identifier, fn, totalChunks, chunkSize}
	reqs = make([]*http.Request, totalChunks)
	data, md5sum := testBytes(100*1024 - 2)
	buf := new(bytes.Buffer)
	for i := 1; i <= totalChunks; i++ {
		io.CopyN(buf, data, chunkSize)
		r := makeRequest(url, buf, f, i)
		reqs[i-1] = r
	}
	return
}
func TestUpload(t *testing.T) {
	dao := dao(t)
	v1 := Video{"", "", "Novocento", "", ""}
	id, err := dao.Create(v1)
	if err != nil {
		t.Fatal(err)
	}
	bs := blobService()
	defer bs.Close()
	ts := videoTestServer(t)
	defer ts.Close()
	url := ts.URL + "/api/videos/" + id + "/upload"
	fl, requests, md5sum := PrepareRequests(url, "foo.data", 100*1024-4, 1024)
	for i := 0; i < fl.TotalChunks; i++ {
		resp, err := http.DefaultClient.Do(requests[i])
		if err != nil {
			t.Fatal("DefaultClient recieves error: ", err)
		}
		if resp.StatusCode != 200 {
			t.Fatal("StatusCode should be 200, but was ", resp.StatusCode)
		}
	}

	time.Sleep(time.Duration(6) * time.Second)
	v2 := new(Video)
	err = dao.Get(id, v2)
	fid := v2.BlobId
	file, err := bs.Open(fid)
	if err != nil {
		t.Fatal("Open file went south: ", err)

	}
	if file.MD5() != md5sum {
		t.Fatal("Checksum of uploaded file mismatched")
	}
	bs.Remove(fid)
}
