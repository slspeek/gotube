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
	return NewVideoResource(sess, "test", "Video", &Auth{func(*http.Request) string { return "steven" }})
}

func videoTestServer(t *testing.T) *httptest.Server {
	vr := videoResource(t)
	container := restful.NewContainer()
	vr.Register(container)
	return httptest.NewServer(container)
}

func TestVideoResource(t *testing.T) {
	ts := videoTestServer(t)
	defer ts.Close()

	resp, _ := http.Post(ts.URL+"/api/videos", "application/json", strings.NewReader(`{"Owner":"","Name":"As it is in heaven"}`))
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
	ts := videoTestServer(t)
	defer ts.Close()
	url := ts.URL + "/api/videos/" + id
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
	vr := videoResource(t)
	req, _ := http.NewRequest("GET", "", nil)
	rreq := restful.NewRequest(req, &[]byte{}, map[string]string{"video-id": id}, map[string]interface{}{"username": "steven"})

	rw := httptest.NewRecorder()
	rresp := restful.NewResponse(rw, "application/json", []string{"application/json"})

	vr.findVideo(rreq, rresp)
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
	vr := videoResource(t)
	reqBody := fmt.Sprintf(`{"Id":"%s","Name":"Novecento II","Owner":"steven"}`, id)
	req, _ := http.NewRequest("PUT", "", strings.NewReader(reqBody))
	req.Header["Content-Type"] = []string{"application/json"}
	rreq := restful.NewRequest(req, &[]byte{}, map[string]string{"video-id": id}, map[string]interface{}{"username": "steven"})

	rw := httptest.NewRecorder()
	rresp := restful.NewResponse(rw, "application/json", []string{"application/json"})

	vr.updateVideo(rreq, rresp)
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
	req, _ := http.NewRequest("PUT", "", strings.NewReader(reqBody))
	req.Header["Content-Type"] = []string{"application/json"}
	rreq := restful.NewRequest(req, &[]byte{}, map[string]string{}, map[string]interface{}{"username": "steven"})

	rw := httptest.NewRecorder()
	rresp := restful.NewResponse(rw, "application/json", []string{"application/json"})

	vr.createVideo(rreq, rresp)
	if rw.Code != http.StatusCreated {
		t.Fatal("StatusCode: ", rw.Code)
	}
}

func TestRemoveUnit(t *testing.T) {
	dao := dao(t)
	dao.DeleteAll()
	v1 := Video{"", "steven", "Novocento", "", ""}
	id, err := dao.Create(v1)
	if err != nil {
		t.Fatal(err)
	}
	vr := videoResource(t)
	req, _ := http.NewRequest("DELETE", "", nil)
	rreq := restful.NewRequest(req, &[]byte{}, map[string]string{"video-id": id}, map[string]interface{}{"username": "steven"})

	rw := httptest.NewRecorder()
	rresp := restful.NewResponse(rw, "application/json", []string{"application/json"})

	vr.removeVideo(rreq, rresp)
	if rw.Code != http.StatusOK {
		t.Fatal("StatusCode: ", rw.Code)
	}
	readBack := make([]Video, 1)
	err = dao.GetAll(&readBack)
	if len(readBack) != 0 {
		t.Fatal("Expected no elemements")
	}
}

func TestRemoveUnitNotOwner(t *testing.T) {
	dao := dao(t)
	dao.DeleteAll()
	v1 := Video{"", "steven", "Novocento", "", ""}
	id, err := dao.Create(v1)
	if err != nil {
		t.Fatal(err)
	}
	vr := videoResource(t)
	req, _ := http.NewRequest("DELETE", "", nil)
	rreq := restful.NewRequest(req, &[]byte{}, map[string]string{"video-id": id}, map[string]interface{}{"username": "rob"})

	rw := httptest.NewRecorder()
	rresp := restful.NewResponse(rw, "application/json", []string{"application/json"})

	vr.removeVideo(rreq, rresp)
	if rw.Code != http.StatusMethodNotAllowed {
		t.Fatal("StatusCode: ", rw.Code)
	}
	readBack := make([]Video, 1)
	err = dao.GetAll(&readBack)
	if len(readBack) != 1 {
		t.Fatal("Expected one element")
	}
}
func TestGetAllUnit(t *testing.T) {
	//t.Skip()
	dao := dao(t)
	dao.DeleteAll()
	v1 := Video{"", "steven", "Novocento", "", ""}
	_, err := dao.Create(v1)
	if err != nil {
		t.Fatal(err)
	}
	v2 := Video{"", "steven", "The spiderwoman", "", ""}
	_, err = dao.Create(v2)
	if err != nil {
		t.Fatal(err)
	}
	v3 := Video{"", "jan", "At the gates", "", ""}
	_, err = dao.Create(v3)
	if err != nil {
		t.Fatal(err)
	}
	vr := videoResource(t)
	req, _ := http.NewRequest("GET", "", nil)
	//req.Header["Content-Type"] = []string{"application/json"}
	rreq := restful.NewRequest(req, &[]byte{}, map[string]string{}, map[string]interface{}{"username": "steven"})

	rw := httptest.NewRecorder()
	rresp := restful.NewResponse(rw, "application/json", []string{"application/json"})

	vr.findAllVideos(rreq, rresp)
	if rw.Code != http.StatusOK {
		t.Fatal("StatusCode: ", rw.Code)
	}
	rb := rw.Body.String()
	//	t.Log("Body: ", rb, len(rb))
	if len(rb) != 257 {
		t.Fatal("Unexpected output from findAllVideos")
	}
	dao.DeleteAll()
}
