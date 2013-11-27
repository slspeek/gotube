package main

import (
	"fmt"
	"github.com/slspeek/go-restful"
	"io"
	"labix.org/v2/mgo"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func videoResource(t *testing.T) *VideoResource {
	sess, err := mgo.Dial("localhost")
	if err != nil {
		t.Fatal(err)
	}
	return NewVideoResource(sess, "test", "Video")
}


func videoTestServer(t *testing.T) *httptest.Server {
	vr := videoResource(t)
	container := restful.NewContainer()
	vr.Register(container)
	return httptest.NewServer(container)
}

func TestVideoResource(t *testing.T) {
	t.Log("Video test")
	ts := videoTestServer(t)
	defer ts.Close()

	resp, _ := http.Post(ts.URL+"/videos", "application/json", strings.NewReader(`{"Owner":"","Name":"As it is in heaven"}`))
	io.Copy(os.Stderr, resp.Body)

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Fatal("StatusCode: ", resp.StatusCode)
	}
}

func TestGet(t *testing.T) {
	dao := dao(t)
	v1 := Video{"", "", "Novocento", "", ""}
	id, err := dao.Create(v1)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Video id: ", id)
	ts := videoTestServer(t)
	defer ts.Close()
	url := ts.URL + "/videos/" + id
	t.Log("URL: ", url)
	resp, err := http.Get(url)
	io.Copy(os.Stderr, resp.Body)

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Fatal("StatusCode: ", resp.StatusCode)
	}
}

func TestFindUnit(t *testing.T) {
	dao := dao(t)
	v1 := Video{"", "steven", "Novocento", "", ""}
	id, err := dao.Create(v1)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Video id: ", id)
	vr := videoResource(t)
	req, _ := http.NewRequest("GET", "", nil)
	rreq := restful.NewRequest(req, &[]byte{}, map[string]string{"video-id": id}, map[string]interface{}{"username": "steven"})

	rw := httptest.NewRecorder()
	rresp := restful.NewResponse(rw, "application/json", []string{"application/json"})

	vr.findVideo(rreq, rresp)
	t.Log("Response: ", rresp, rw.Body.String())
	if rw.Code != http.StatusOK {
		t.Fatal("StatusCode: ", rw.Code)
	}
}
func TestUpdateUnit(t *testing.T) {
	//t.Skip()
	dao := dao(t)
	v1 := Video{"", "steven", "Novocento", "", ""}
	id, err := dao.Create(v1)
	if err != nil {
		t.Fatal(err)
	}
	t.Log("Video id: ", id)
	vr := videoResource(t)
	reqBody := fmt.Sprintf(`{"Id":"%s","Name":"Novecento II","Owner":"steven"}`, id)
	t.Log("Body: ", reqBody)
	req, _ := http.NewRequest("PUT", "", strings.NewReader(reqBody))
  req.Header["Content-Type"] = []string {"application/json"}
	rreq := restful.NewRequest(req, &[]byte{}, map[string]string{"video-id": id}, map[string]interface{}{"username": "steven"})

	rw := httptest.NewRecorder()
	rresp := restful.NewResponse(rw, "application/json", []string{"application/json"})

	vr.updateVideo(rreq, rresp)
	t.Log("Response: ", rresp, rw.Body.String())
	if rw.Code != http.StatusOK {
		t.Fatal("StatusCode: ", rw.Code)
	}
	readBack := new(Video)
	err = dao.Get(id, &readBack)
	if readBack.Name != "Novecento II" {
		t.Fatal("Expected the sequel")
	}
}

func TestCreate(t *testing.T) {
	vr := videoResource(t)
	reqBody := fmt.Sprintf(`{"Name":"Prison break","Owner":"steven"}`)
	t.Log("Body: ", reqBody)
	req, _ := http.NewRequest("PUT", "", strings.NewReader(reqBody))
  req.Header["Content-Type"] = []string {"application/json"}
	rreq := restful.NewRequest(req, &[]byte{}, map[string]string{}, map[string]interface{}{"username": "steven"})

	rw := httptest.NewRecorder()
	rresp := restful.NewResponse(rw, "application/json", []string{"application/json"})

	vr.createVideo(rreq, rresp)
	t.Log("Response: ", rresp, rw.Body.String())
	if rw.Code != http.StatusCreated {
		t.Fatal("StatusCode: ", rw.Code)
	}
}
