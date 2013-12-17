package rest

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"github.com/slspeek/flowgo"
	"github.com/slspeek/go-restful"
	"github.com/slspeek/goblob"
	"github.com/slspeek/gotube/auth"
	"github.com/slspeek/gotube/mongo"
	"io"
	"labix.org/v2/mgo"
	"math"
	"math/rand"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
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
	bs = goblob.NewBlobService(s, "test", "flowfs")
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
	id, err := dao.Create(NOVECENTO)
	check(t, err)
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
	var i int
	for i = 1; i < 200; i++ {
		time.Sleep(time.Duration(100) * time.Millisecond)
		v2 := new(Video)
		v2.BlobId = ""
		err = dao.Get(id, v2)
		fid := v2.BlobId
		if fid == "" {
			continue
		}
		file, err := bs.Open(fid)
		if err != nil {
			continue
		}
		if file.MD5() != md5sum {
			t.Fatal("Checksum of uploaded file mismatched")
		}

		t.Log("Checked")
		bs.Remove(fid)
		break
	}
	if i >= 200 {
		t.Fatal("Timed out")
	}
}

func TestDownloadWhenNotReady(t *testing.T) {
	dao := dao(t)
	id, err := dao.Create(NOVECENTO)
	if err != nil {
		t.Fatal(err)
	}

	ts := videoTestServer(t)
	defer ts.Close()
	req, _ := http.NewRequest("GET", ts.URL+"/content/videos/"+id, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal("Could not get, due to ", err)
	}
	if resp.StatusCode != 204 {
		t.Fatal("Get Returned error: ", resp.StatusCode)
	}
}

func TestDownload(t *testing.T) {
	dao := dao(t)
	id, err := dao.Create(NOVECENTO)
	if err != nil {
		t.Fatal(err)
	}
	bs := blobService()
	defer bs.Close()
	ts := videoTestServer(t)
	defer ts.Close()

	buf, md5sum := testBytes(100*1024 - 4)
	out, err := bs.Create("foo.bar")
	if err != nil {
		t.Fatal("Could not create new grid file")
	}
	io.Copy(out, buf)
	fid := out.Id()
	out.Close()
	v := new(Video)
	err = dao.Get(id, v)
	if err != nil {
		t.Fatal("Could not reload video")
	}
	v.BlobId = fid
	err = dao.Update(id, v)
	if err != nil {
		t.Fatal("Could not update video with blobid")
	}

	req, _ := http.NewRequest("GET", ts.URL+"/content/videos/"+id, nil)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal("Could not get, due to ", err)
	}
	if resp.StatusCode != 200 {
		t.Fatal("Get Returned error: ", resp.StatusCode)
	}
	hash := md5.New()
	io.Copy(hash, resp.Body)

	readBackMD5 := fmt.Sprintf("%x", hash.Sum(nil))
	if md5sum != readBackMD5 {
		t.Fatal("Checksum of downloaded file mismatches: ", readBackMD5, md5sum)
	}
	bs.Remove(fid)
}

func TestConcurrentDownloads(t *testing.T) {
	concurrency := 20
	dao := dao(t)
	id, err := dao.Create(NOVECENTO)
	if err != nil {
		t.Fatal(err)
	}
	bs := blobService()
	defer bs.Close()
	ts := videoTestServer(t)
	defer ts.Close()

	buf, md5sum := testBytes(100*1024 - 4)
	out, err := bs.Create("foo.bar")
	if err != nil {
		t.Fatal("Could not create new grid file")
	}
	io.Copy(out, buf)
	fid := out.Id()
	out.Close()
	v := new(Video)
	err = dao.Get(id, v)
	if err != nil {
		t.Fatal("Could not reload video")
	}
	v.BlobId = fid
	err = dao.Update(id, v)
	if err != nil {
		t.Fatal("Could not update video with blobid")
	}
	endSignal := make(chan bool)
	for i := 0; i < concurrency; i++ {
		go func() {
			req, _ := http.NewRequest("GET", ts.URL+"/content/videos/"+id, nil)
			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				t.Fatal("Could not get, due to ", err)
			}
			if resp.StatusCode != 200 {
				t.Fatal("Get Returned error: ", resp.StatusCode)
			}
			hash := md5.New()
			io.Copy(hash, resp.Body)

			readBackMD5 := fmt.Sprintf("%x", hash.Sum(nil))
			if md5sum != readBackMD5 {
				t.Fatal("Checksum of downloaded file mismatches: ", readBackMD5, md5sum)
			}
			endSignal <- true
		}()
	}
	for i := 0; i < concurrency; i++ {
		<-endSignal
	}
	bs.Remove(fid)
}

func BenchmarkConcurrentDownloads(b *testing.B) {
	b.StopTimer()
	sess, err := mgo.Dial("localhost")
	if err != nil {
		b.Fatal(err)
	}
	dao := mongo.NewDao(sess, "test", "Video")
	concurrency := 100
	id, err := dao.Create(NOVECENTO)
	if err != nil {
		b.Fatal(err)
	}
	bs := blobService()
	defer bs.Close()
	container := restful.NewContainer()
	vr := NewVideoResource(sess, "test", "Video", &auth.Auth{allwaysSteven})
	vr.Register(container)
	ts := httptest.NewServer(container)
	defer ts.Close()

  //test file is 100M
	buf, md5sum := testBytes(100*1024*1024 - 4)
	out, err := bs.Create("foo.bar")
	if err != nil {
		b.Fatal("Could not create new grid file")
	}
	io.Copy(out, buf)
	fid := out.Id()
	out.Close()
	v := new(Video)
	err = dao.Get(id, v)
	if err != nil {
		b.Fatal("Could not reload video")
	}
	v.BlobId = fid
	err = dao.Update(id, v)
	if err != nil {
		b.Fatal("Could not update video with blobid")
	}
	b.StartTimer()
	for n := 0; n < b.N; n++ {
		endSignal := make(chan bool)
		for i := 0; i < concurrency; i++ {
			go func() {
				req, _ := http.NewRequest("GET", ts.URL+"/content/videos/"+id, nil)
				resp, err := http.DefaultClient.Do(req)
				if err != nil {
					b.Fatal("Could not get, due to ", err)
				}
				if resp.StatusCode != 200 {
					b.Fatal("Get Returned error: ", resp.StatusCode)
				}
				hash := md5.New()
				io.Copy(hash, resp.Body)

				readBackMD5 := fmt.Sprintf("%x", hash.Sum(nil))
				if md5sum != readBackMD5 {
					b.Fatal("Checksum of downloaded file mismatches: ", readBackMD5, md5sum)
				}
				endSignal <- true
			}()
		}
		for i := 0; i < concurrency; i++ {
			<-endSignal

		}
	}
	bs.Remove(fid)
}
