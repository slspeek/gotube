package rest

import (
	"fmt"
	"github.com/slspeek/go-restful"
	"github.com/slspeek/goblob"
	"github.com/slspeek/gotube/auth"
	"github.com/slspeek/gotube/common"
	"github.com/slspeek/gotube/mongo"
	"io"
	"labix.org/v2/mgo"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

var (
	NOVECENTO     = common.Video{Owner: "steven", Name: "Novecento", Desc: "Italian classic"}
	allwaysSteven = func(*http.Request) string { return "steven" }
	allwaysNobody = func(*http.Request) string { return "" }
)

func insertLittleBlob() (id string, err error) {
	sess, err := mgo.Dial("localhost")
	if err != nil {
		return
	}
	bs := goblob.NewBlobService(sess, "test", "flowfs")
	file, err := bs.Create("foo")
	if err != nil {
		return
	}
	fmt.Fprint(file, "Hello")
	err = file.Close()
	id = file.StringId()
	return
}
func dao(t *testing.T) *mongo.Dao {
	sess, err := mgo.Dial("localhost")
	if err != nil {
		t.Fatal(err)
	}
	return mongo.NewDao(sess, "test", "Video")
}

func createNovecento(t *testing.T) (id string) {
	dao := dao(t)
	id, err := dao.Create(NOVECENTO)
	if err != nil {
		t.Fatal(err)
	}
	return
}

func videoResource(t *testing.T, checker auth.CheckAuthFunc) *VideoResource {
	sess, err := mgo.Dial("localhost")
	if err != nil {
		t.Fatal(err)
	}
	return NewVideoResource(sess, "test", "Video", &auth.Auth{checker})
}

func videoTestServer(t *testing.T) *httptest.Server {
	return httptest.NewServer(loadedContainer(t, allwaysSteven))
}

func loadedContainer(t *testing.T, checker auth.CheckAuthFunc) *restful.Container {
	vr := videoResource(t, checker)
	container := restful.NewContainer()
	vr.Register(container)
	return container
}

func TestVideoTag(t *testing.T) {
	tag := videoTag(NOVECENTO, "host")
	if tag != `<video src="host/content/videos/" />` {
		t.Log(tag)
		t.Fatal("Wrong video tag")
	}

}

// Test Matrix
//   Properties:
//   IsOwner
//   InvalidId
//   Not authorized
//
//   methods of VideoResource:
//   findVideo, updateVideo, removeVideo
//
//   createVideo, findAllVideo serapately variing Not authorized
//   createVideo with invalid data

func TestNotAuthorizedCreate(t *testing.T) {
	container := loadedContainer(t, allwaysNobody)
	req, _ := http.NewRequest("POST", "/api/videos", nil)
	rw := httptest.NewRecorder()
	container.ServeHTTP(rw, req)
	if rw.Code != http.StatusUnauthorized {
		t.Fatal("Should have been unauthorized, but was: ", rw.Code)
	}
}

func TestInvalidDataCreate(t *testing.T) {
	container := loadedContainer(t, allwaysSteven)
	req, _ := http.NewRequest("POST", "/api/videos", nil)
	rw := httptest.NewRecorder()
	container.ServeHTTP(rw, req)
	if rw.Code != http.StatusInternalServerError {
		t.Fatal("Should have been internalerror, but was: ", rw.Code)
	}
}

func TestSuccessCreate(t *testing.T) {
	container := loadedContainer(t, allwaysSteven)
	req, _ := http.NewRequest("POST", "/api/videos", strings.NewReader(`{"Owner":"steven","Name":"As it is in heaven"}`))
	req.Header["Content-Type"] = []string{"application/json"}
	rw := httptest.NewRecorder()
	container.ServeHTTP(rw, req)
	if rw.Code != http.StatusCreated {
		t.Fatal("Should have been created, but was: ", rw.Code)
	}
}

func TestNotAuthorizedFindAll(t *testing.T) {
	container := loadedContainer(t, allwaysNobody)
	req, _ := http.NewRequest("GET", "/api/videos", nil)
	rw := httptest.NewRecorder()
	container.ServeHTTP(rw, req)
	if rw.Code != http.StatusUnauthorized {
		t.Fatal("Should have been unauthorized, but was: ", rw.Code)
	}
}

func TestSuccessFindAll(t *testing.T) {
	createNovecento(t)
	container := loadedContainer(t, allwaysSteven)
	req, _ := http.NewRequest("GET", "/api/videos", nil)
	rw := httptest.NewRecorder()
	container.ServeHTTP(rw, req)
	if rw.Code != http.StatusOK {
		t.Fatal("Should have been OK, but was: ", rw.Code)
	}
}
func TestNotAuthorizedPlaylist(t *testing.T) {
	container := loadedContainer(t, allwaysNobody)
	req, _ := http.NewRequest("GET", "/api/videos/playlist", nil)
	rw := httptest.NewRecorder()
	container.ServeHTTP(rw, req)
	if rw.Code != http.StatusUnauthorized {
		t.Fatal("Should have been unauthorized, but was: ", rw.Code)
	}
}

func TestSuccessPlaylist(t *testing.T) {
	dao := dao(t)
	dao.DeleteAll()
	container := loadedContainer(t, allwaysSteven)
	req, _ := http.NewRequest("GET", "/api/videos/playlist", nil)
	rw := httptest.NewRecorder()
	container.ServeHTTP(rw, req)
	if rw.Code != http.StatusOK {
		t.Fatal("Should have been OK, but was: ", rw.Code)
	}
	if body := rw.Body.String(); body != `<smil><body><seq></seq></body></smil>` {
		t.Fatal("Expected ", body)
	}
}

func TestInvalidIdFind(t *testing.T) {
	container := loadedContainer(t, allwaysSteven)
	req, _ := http.NewRequest("GET", "/api/videos/999", nil)
	rw := httptest.NewRecorder()
	container.ServeHTTP(rw, req)
	if rw.Code != http.StatusBadRequest {
		t.Fatal("Should have been StatusBadRequest, but was: ", rw.Code)
	}
}

func TestNotOwnerFind(t *testing.T) {
	dao := dao(t)
  v1 := common.Video{Owner:"rob", Name:"Novocento"}
	id, err := dao.Create(v1)
	if err != nil {
		t.Fatal(err)
	}
	container := loadedContainer(t, allwaysSteven)
	req, _ := http.NewRequest("GET", "/api/videos/"+id, nil)
	rw := httptest.NewRecorder()
	container.ServeHTTP(rw, req)
	if rw.Code != http.StatusForbidden {
		t.Fatal("Should have been StatusForbidden, but was: ", rw.Code)
	}
}

func TestNotAuthorizedFind(t *testing.T) {
	id := createNovecento(t)
	container := loadedContainer(t, allwaysNobody)
	req, _ := http.NewRequest("GET", "/api/videos/"+id, nil)
	rw := httptest.NewRecorder()
	container.ServeHTTP(rw, req)
	if rw.Code != http.StatusUnauthorized {
		t.Fatal("Should have been unauthorized, but was: ", rw.Code)
	}

}

func TestSuccessFind(t *testing.T) {
	id := createNovecento(t)
	container := loadedContainer(t, allwaysSteven)
	req, _ := http.NewRequest("GET", "/api/videos/"+id, nil)
	rw := httptest.NewRecorder()
	container.ServeHTTP(rw, req)
	if rw.Code != http.StatusOK {
		t.Fatal("Should have been Ok, but was: ", rw.Code)
	}
}

func TestInvalidIdUpdate(t *testing.T) {
	createNovecento(t)
	container := loadedContainer(t, allwaysSteven)
	req, _ := http.NewRequest("PUT", "/api/videos/888", strings.NewReader(`{"Owner":"steven","Name":"As it is in heaven"}`))
	req.Header["Content-Type"] = []string{"application/json"}
	rw := httptest.NewRecorder()
	container.ServeHTTP(rw, req)
	if rw.Code != http.StatusBadRequest {
		t.Fatal("Should have been bad request, but was: ", rw.Code)
	}
}

func TestNotOwnerUpdate(t *testing.T) {
	dao := dao(t)
  v1 := common.Video{Owner:"rob", Name:"Novocento"}
	id, err := dao.Create(v1)
	if err != nil {
		t.Fatal(err)
	}
	container := loadedContainer(t, allwaysSteven)
	req, _ := http.NewRequest("PUT", "/api/videos/"+id, strings.NewReader(`{"Owner":"steven","Name":"As it is in heaven"}`))
	req.Header["Content-Type"] = []string{"application/json"}
	rw := httptest.NewRecorder()
	container.ServeHTTP(rw, req)
	if rw.Code != http.StatusForbidden {
		t.Fatal("Should have been StatusForbidden, but was: ", rw.Code)
	}
}

func TestNotAuthorizedUpdate(t *testing.T) {
	id := createNovecento(t)
	container := loadedContainer(t, allwaysNobody)
	req, _ := http.NewRequest("PUT", "/api/videos/"+id, nil)
	rw := httptest.NewRecorder()
	container.ServeHTTP(rw, req)
	if rw.Code != http.StatusUnauthorized {
		t.Fatal("Should have been unauthorized, but was: ", rw.Code)
	}
}

func TestSuccessUpdate(t *testing.T) {
	id := createNovecento(t)
	container := loadedContainer(t, allwaysSteven)
	req, _ := http.NewRequest("PUT", "/api/videos/"+id, strings.NewReader(`{"Owner":"steven","Name":"As it is in heaven"}`))
	req.Header["Content-Type"] = []string{"application/json"}
	rw := httptest.NewRecorder()
	container.ServeHTTP(rw, req)
	if rw.Code != http.StatusOK {
		t.Fatal("Should have been OK, but was: ", rw.Code)
	}
}

func TestInvalidIdRemove(t *testing.T) {
	createNovecento(t)
	container := loadedContainer(t, allwaysSteven)
	req, _ := http.NewRequest("DELETE", "/api/videos/888", nil)
	rw := httptest.NewRecorder()
	container.ServeHTTP(rw, req)
	if rw.Code != http.StatusBadRequest {
		t.Fatal("Should have been bad request, but was: ", rw.Code)
	}

}

func TestNotOwnerRemove(t *testing.T) {
	dao := dao(t)
  v1 := common.Video{Owner:"rob", Name:"Novocento"}
	id, err := dao.Create(v1)
	if err != nil {
		t.Fatal(err)
	}

	container := loadedContainer(t, allwaysSteven)
	req, _ := http.NewRequest("DELETE", "/api/videos/"+id, nil)
	rw := httptest.NewRecorder()
	container.ServeHTTP(rw, req)
	if rw.Code != http.StatusForbidden {
		t.Fatal("Should have been forbidden, but was: ", rw.Code)
	}
}

func TestNotAuthorizedRemove(t *testing.T) {
	id := createNovecento(t)
	container := loadedContainer(t, allwaysNobody)
	req, _ := http.NewRequest("DELETE", "/api/videos/"+id, nil)
	rw := httptest.NewRecorder()
	container.ServeHTTP(rw, req)
	if rw.Code != http.StatusUnauthorized {
		t.Fatal("Should have been unauthorized, but was: ", rw.Code)
	}
}

func TestSuccessRemoveWithoutBlob(t *testing.T) {
	id := createNovecento(t)
	container := loadedContainer(t, allwaysSteven)
	req, _ := http.NewRequest("DELETE", "/api/videos/"+id, nil)
	rw := httptest.NewRecorder()
	container.ServeHTTP(rw, req)
	if rw.Code != http.StatusOK {
		t.Fatal("Should have been ok, but was: ", rw.Code)
	}
	dao := dao(t)
	err := dao.Get(id, new(common.Video))
	if err == nil {
		t.Fatal("Should have been removed")
	}
}

func TestSuccessRemoveWithBlob(t *testing.T) {
	dao := dao(t)
	id := createNovecento(t)
	blobid, err := insertLittleBlob()
	if err != nil {
		t.Fatal("Could not create little file")
	}
	vid := new(common.Video)
	err = dao.Get(id, vid)
	if err != nil {
		t.Fatal("Could not reload Video")
	}
	vid.BlobId = blobid
	err = dao.Update(id, vid)
	if err != nil {
		t.Fatal("Could not save video")
	}
	container := loadedContainer(t, allwaysSteven)
	req, _ := http.NewRequest("DELETE", "/api/videos/"+id, nil)
	rw := httptest.NewRecorder()
	container.ServeHTTP(rw, req)
	if rw.Code != http.StatusOK {
		t.Fatal("Should have been ok, but was: ", rw.Code)
	}
	err = dao.Get(id, new(common.Video))
	if err == nil {
		t.Fatal("Should have been removed")
	}

	sess, err := mgo.Dial("localhost")
	if err != nil {
		return
	}
	bs := goblob.NewBlobService(sess, "test", "flowfs")
	_, err = bs.Open(blobid)
	if err == nil {
		t.Fatal("Blob Should have been removed")
	}
}

func TestVideoResource(t *testing.T) {
	ts := videoTestServer(t)
	defer ts.Close()

	resp, _ := http.Post(ts.URL+"/api/videos", "application/json", strings.NewReader(`{"Owner":"","Name":"As it is in heaven"}`))
	io.Copy(os.Stderr, resp.Body)

	if resp.StatusCode != http.StatusForbidden {
		t.Fatal("StatusCode: ", resp.StatusCode)
	}
}

func TestGet(t *testing.T) {
	dao := dao(t)
  v1 := common.Video{Name:"Novocento"}
	id, err := dao.Create(v1)
	if err != nil {
		t.Fatal(err)
	}
	ts := videoTestServer(t)
	defer ts.Close()
	url := ts.URL + "/api/videos/" + id
	resp, err := http.Get(url)
	io.Copy(os.Stderr, resp.Body)

	if resp.StatusCode != http.StatusForbidden {
		t.Fatal("StatusCode: ", resp.StatusCode)
	}
}

func TestInvalidIdDownload(t *testing.T) {
	container := loadedContainer(t, allwaysSteven)
	req, _ := http.NewRequest("GET", "/content/videos/999/download", nil)
	rw := httptest.NewRecorder()
	container.ServeHTTP(rw, req)
	if rw.Code != http.StatusBadRequest {
		t.Fatal("Should have been StatusBadRequest, but was: ", rw.Code)
	}
}

func TestNotOwnerDownload(t *testing.T) {
	dao := dao(t)
  v1 := common.Video{Owner:"rob", Name:"Novocento"}
	id, err := dao.Create(v1)
	if err != nil {
		t.Fatal(err)
	}
	container := loadedContainer(t, allwaysSteven)
	req, _ := http.NewRequest("GET", "/content/videos/"+id+"/download", nil)
	rw := httptest.NewRecorder()
	container.ServeHTTP(rw, req)
	if rw.Code != http.StatusForbidden {
		t.Fatal("Should have been StatusForbidden, but was: ", rw.Code)
	}
}

func TestNotAuthorizedDownload(t *testing.T) {
	id := createNovecento(t)
	container := loadedContainer(t, allwaysNobody)
	req, _ := http.NewRequest("GET", "/content/videos/"+id+"/download", nil)
	rw := httptest.NewRecorder()
	container.ServeHTTP(rw, req)
	if rw.Code != http.StatusUnauthorized {
		t.Fatal("Should have been unauthorized, but was: ", rw.Code)
	}

}

func TestSuccessDownload(t *testing.T) {
	dao := dao(t)
	id := createNovecento(t)
	blobid, err := insertLittleBlob()
	if err != nil {
		t.Fatal("Could not create little file")
	}
	vid := new(common.Video)
	err = dao.Get(id, vid)
	if err != nil {
		t.Fatal("Could not reload Video")
	}
	vid.BlobId = blobid
	err = dao.Update(id, vid)
	if err != nil {
		t.Fatal("Could not save video")
	}
	container := loadedContainer(t, allwaysSteven)
	req, _ := http.NewRequest("GET", "/content/videos/"+id+"/download", nil)
	rw := httptest.NewRecorder()
	container.ServeHTTP(rw, req)
	if rw.Code != http.StatusOK {
		t.Fatal("Should have been Ok, but was: ", rw.Code)
	}
	contentType := rw.Header().Get("Content-Type")
	if contentType != "application/octet-stream" {
		t.Fatal("Expected octet-stream, got: ", contentType)
	}

	disposition := rw.Header().Get("Content-Disposition")
	if disposition != `attachment; filename="foo"` {
		t.Fatal("Expected attachment got: ", disposition)
	}

}
