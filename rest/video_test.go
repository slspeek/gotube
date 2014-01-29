package rest

import (
	"fmt"
	"github.com/slspeek/go-restful"
	"github.com/slspeek/goblob"
	"github.com/slspeek/gotube/auth"
	"github.com/slspeek/gotube/common"
	"github.com/slspeek/gotube/mongo"
	"io/ioutil"
	"labix.org/v2/mgo"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var (
	NOVECENTO      = common.Video{Owner: "steven", Name: "Novecento", Desc: "Italian classic"}
	PUBLIC_FICTION = common.Video{Owner: "steven", Public: true, Name: "Public Fiction", Desc: "Unreal classic"}
	allwaysSteven  = func(*http.Request) string { return "steven" }
	allwaysNobody  = func(*http.Request) string { return "" }
)

func blobService() (bs *goblob.BlobService) {
	s, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	bs = goblob.NewBlobService(s, "test", "flowfs")
	return
}

func insertLittleBlob() (id string, err error) {
	bs := blobService()
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
func createPublicFiction(t *testing.T) (id string) {
	dao := dao(t)
	id, err := dao.Create(PUBLIC_FICTION)
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

//Generic test functions
func verifyRequestWillNotAuthorize(req *http.Request, t *testing.T) {
	container := loadedContainer(t, allwaysNobody)
	rw := httptest.NewRecorder()
	container.ServeHTTP(rw, req)
	if rw.Code != http.StatusUnauthorized {
		t.Fatal("Should have been unauthorized, but was: ", rw.Code)
	}
}

func getVideoIdOwnedByRob(t *testing.T) string {
	dao := dao(t)
	v1 := common.Video{Owner: "rob", Name: "Novocento"}
	id, err := dao.Create(v1)
	if err != nil {
		t.Fatal(err)
	}
	return id
}

func verifyRequestGetsForbiddenForUserSteven(req *http.Request, t *testing.T) *httptest.ResponseRecorder {
	return verifyRequestReturns(req, http.StatusForbidden, t)
}

func verifyRequestReturnsBadRequest(req *http.Request, t *testing.T) *httptest.ResponseRecorder {
	return verifyRequestReturns(req, http.StatusBadRequest, t)
}
func verifyRequestReturnsInternalError(req *http.Request, t *testing.T) *httptest.ResponseRecorder {
	return verifyRequestReturns(req, http.StatusInternalServerError, t)
}
func verifyRequestReturnsOK(req *http.Request, t *testing.T) *httptest.ResponseRecorder {
	return verifyRequestReturns(req, http.StatusOK, t)
}

func verifyRequestReturns(req *http.Request, status int, t *testing.T) (rw *httptest.ResponseRecorder) {
	container := loadedContainer(t, allwaysSteven)
	rw = httptest.NewRecorder()
	container.ServeHTTP(rw, req)
	if rw.Code != status {
		t.Fatal("Should have been: ", status, " but was: ", rw.Code)
	}
	return rw
}

//Create method tests
func TestNotAuthorizedCreate(t *testing.T) {
	req, _ := http.NewRequest("POST", "/api/videos", nil)
	verifyRequestWillNotAuthorize(req, t)
}

func TestInvalidDataCreate(t *testing.T) {
	req, _ := http.NewRequest("POST", "/api/videos", nil)
	verifyRequestReturnsInternalError(req, t)
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

//Findall method tests
func TestNotAuthorizedFindAll(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/videos", nil)
	verifyRequestWillNotAuthorize(req, t)
}

func TestSuccessFindAll(t *testing.T) {
	createNovecento(t)
	req, _ := http.NewRequest("GET", "/api/videos", nil)
	verifyRequestReturnsOK(req, t)
}

//Playlist method tests
func TestNotAuthorizedPlaylist(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/videos/playlist", nil)
	verifyRequestWillNotAuthorize(req, t)
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

//Find method tests
func TestInvalidIdFind(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/videos/999", nil)
	verifyRequestReturnsBadRequest(req, t)
}

func TestNotOwnerFind(t *testing.T) {
	id := getVideoIdOwnedByRob(t)
	req, _ := http.NewRequest("GET", "/api/videos/"+id, nil)
	verifyRequestGetsForbiddenForUserSteven(req, t)
}

func TestNotAuthorizedFind(t *testing.T) {
	id := createNovecento(t)
	req, _ := http.NewRequest("GET", "/api/videos/"+id, nil)
	verifyRequestWillNotAuthorize(req, t)
}

func TestSuccessFind(t *testing.T) {
	id := createNovecento(t)
	req, _ := http.NewRequest("GET", "/api/videos/"+id, nil)
	verifyRequestReturnsOK(req, t)
}

//Update method tests
func TestInvalidIdUpdate(t *testing.T) {
	createNovecento(t)
	req, _ := http.NewRequest("PUT", "/api/videos/888", strings.NewReader(`{"Owner":"steven","Name":"As it is in heaven"}`))
	req.Header["Content-Type"] = []string{"application/json"}
	verifyRequestReturnsBadRequest(req, t)
}

func TestNotOwnerUpdate(t *testing.T) {
	id := getVideoIdOwnedByRob(t)
	req, _ := http.NewRequest("PUT", "/api/videos/"+id, strings.NewReader(`{"Owner":"steven","Name":"As it is in heaven"}`))
	req.Header["Content-Type"] = []string{"application/json"}
	verifyRequestGetsForbiddenForUserSteven(req, t)
}

func TestNotAuthorizedUpdate(t *testing.T) {
	id := createNovecento(t)
	req, _ := http.NewRequest("PUT", "/api/videos/"+id, nil)
	verifyRequestWillNotAuthorize(req, t)
}

func TestSuccessUpdate(t *testing.T) {
	id := createNovecento(t)
  req, _ := http.NewRequest("PUT", "/api/videos/"+id, strings.NewReader(`{"Owner":"steven","Public":true,"Name":"As it is in heaven"}`))
	req.Header["Content-Type"] = []string{"application/json"}
	verifyRequestReturnsOK(req, t)
  reloaded := new(common.Video)
  err := dao(t).Get(id, reloaded)
  if err != nil {
    t.Fatal(err)
  }
  if reloaded.Name != "As it is in heaven" {
    t.Fatal()
  }
  if reloaded.Public != true {
    t.Fatal()
  }
}

//Remove methods tests
func TestInvalidIdRemove(t *testing.T) {
	createNovecento(t)
	req, _ := http.NewRequest("DELETE", "/api/videos/888", nil)
	verifyRequestReturnsBadRequest(req, t)
}

func TestNotOwnerRemove(t *testing.T) {
	id := getVideoIdOwnedByRob(t)
	req, _ := http.NewRequest("DELETE", "/api/videos/"+id, nil)
	verifyRequestGetsForbiddenForUserSteven(req, t)
}

func TestNotAuthorizedRemove(t *testing.T) {
	id := createNovecento(t)
	req, _ := http.NewRequest("DELETE", "/api/videos/"+id, nil)
	verifyRequestWillNotAuthorize(req, t)
}

func TestSuccessRemoveWithoutBlob(t *testing.T) {
	id := createNovecento(t)
	req, _ := http.NewRequest("DELETE", "/api/videos/"+id, nil)
	verifyRequestReturnsOK(req, t)
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
	req, _ := http.NewRequest("DELETE", "/api/videos/"+id, nil)
	verifyRequestReturnsOK(req, t)
	err = dao.Get(id, new(common.Video))
	if err == nil {
		t.Fatal("Should have been removed")
	}
	bs := blobService()
	_, err = bs.Open(blobid)
	if err == nil {
		t.Fatal("Blob Should have been removed")
	}
}

func TestSuccessRemoveWithThumb(t *testing.T) {
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
	vid.Thumbs = []string{blobid}
	err = dao.Update(id, vid)
	if err != nil {
		t.Fatal("Could not save video")
	}
	req, _ := http.NewRequest("DELETE", "/api/videos/"+id, nil)
	verifyRequestReturnsOK(req, t)
	err = dao.Get(id, new(common.Video))
	if err == nil {
		t.Fatal("Should have been removed")
	}

	bs := blobService()
	_, err = bs.Open(blobid)
	if err == nil {
		t.Fatal("Thumb Should have been removed")
	}
}

//Download method tests
func TestInvalidIdDownload(t *testing.T) {
	req, _ := http.NewRequest("GET", "/content/videos/999/download", nil)
	verifyRequestReturnsBadRequest(req, t)
}

func TestNotOwnerDownload(t *testing.T) {
	id := getVideoIdOwnedByRob(t)
	req, _ := http.NewRequest("GET", "/content/videos/"+id+"/download", nil)
	verifyRequestGetsForbiddenForUserSteven(req, t)
}

func TestNotAuthorizedDownload(t *testing.T) {
	id := createNovecento(t)
	req, _ := http.NewRequest("GET", "/content/videos/"+id+"/download", nil)
	verifyRequestWillNotAuthorize(req, t)
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
	req, _ := http.NewRequest("GET", "/content/videos/"+id+"/download", nil)
	rw := verifyRequestReturnsOK(req, t)
	contentType := rw.Header().Get("Content-Type")
	if contentType != "application/octet-stream" {
		t.Fatal("Expected octet-stream, got: ", contentType)
	}

	disposition := rw.Header().Get("Content-Disposition")
	if disposition != `attachment; filename="foo"` {
		t.Fatal("Expected attachment got: ", disposition)
	}
}

func TestPublicFindAll(t *testing.T) {
	dao(t).DeleteAll()
	createNovecento(t)
	req, _ := http.NewRequest("GET", "/public/api/videos", nil)
	rw := verifyRequestReturnsOK(req, t)

	b, err := ioutil.ReadAll(rw.Body)
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != "[]" {
		t.Fatal()
	}

	createPublicFiction(t)
	req, _ = http.NewRequest("GET", "/public/api/videos", nil)
	rw = verifyRequestReturnsOK(req, t)
	b, err = ioutil.ReadAll(rw.Body)
	if err != nil {
		t.Fatal(err)
	}
	output := string(b)
	if !strings.Contains(output, "Public Fiction") {
		t.Fatal(output)
	}
	if strings.Contains(output, "Novecento") {
		t.Fatal(output)
	}
}
