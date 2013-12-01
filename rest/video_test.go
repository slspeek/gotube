package rest

import (
	"fmt"
	"github.com/slspeek/go-restful"
	"github.com/slspeek/gotube/auth"
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
	NOVECENTO     = Video{Owner: "steven", Name: "Novecento", Desc: "Italian classic"}
	allwaysSteven = func(*http.Request) string { return "steven" }
	allwaysNobody = func(*http.Request) string { return "" }
)

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
	v1 := Video{"", "rob", "Novocento", "", ""}
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
	v1 := Video{"", "rob", "Novocento", "", ""}
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
	v1 := Video{"", "rob", "Novocento", "", ""}
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

func TestSuccessRemove(t *testing.T) {
	id := createNovecento(t)
	container := loadedContainer(t, allwaysSteven)
	req, _ := http.NewRequest("DELETE", "/api/videos/"+id, nil)
	rw := httptest.NewRecorder()
	container.ServeHTTP(rw, req)
	if rw.Code != http.StatusOK {
		t.Fatal("Should have been ok, but was: ", rw.Code)
	}
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

	if resp.StatusCode != http.StatusForbidden {
		t.Fatal("StatusCode: ", resp.StatusCode)
	}
}

func TestFindUnit(t *testing.T) {
	id := createNovecento(t)
	vr := videoResource(t, allwaysSteven)
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
	id := createNovecento(t)
	vr := videoResource(t, allwaysSteven)
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
	dao := dao(t)
	readBack := new(Video)
	err := dao.Get(id, &readBack)
	if err != nil {
		t.Fatal(err)
	}
	if readBack.Name != "Novecento II" {
		t.Fatal("Expected the sequel")
	}
}

func TestCreate(t *testing.T) {
	vr := videoResource(t, allwaysSteven)
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
	id := createNovecento(t)
	vr := videoResource(t, allwaysSteven)
	req, _ := http.NewRequest("DELETE", "", nil)
	rreq := restful.NewRequest(req, &[]byte{}, map[string]string{"video-id": id}, map[string]interface{}{"username": "steven"})

	rw := httptest.NewRecorder()
	rresp := restful.NewResponse(rw, "application/json", []string{"application/json"})

	vr.removeVideo(rreq, rresp)
	if rw.Code != http.StatusOK {
		t.Fatal("StatusCode: ", rw.Code)
	}
	readBack := make([]Video, 1)
	err := dao.GetAll(&readBack)
	if err != nil {
		t.Fatal(err)
	}
	if len(readBack) != 0 {
		t.Fatal("Expected no elemements")
	}
}

func TestRemoveUnitNotOwner(t *testing.T) {
	dao := dao(t)
	dao.DeleteAll()
	id := createNovecento(t)
	vr := videoResource(t, allwaysSteven)
	req, _ := http.NewRequest("DELETE", "", nil)
	rreq := restful.NewRequest(req, &[]byte{}, map[string]string{"video-id": id}, map[string]interface{}{"username": "rob"})

	rw := httptest.NewRecorder()
	rresp := restful.NewResponse(rw, "application/json", []string{"application/json"})

	vr.removeVideo(rreq, rresp)
	if rw.Code != http.StatusForbidden {
		t.Fatal("StatusCode: ", rw.Code)
	}
	readBack := make([]Video, 1)
	err := dao.GetAll(&readBack)
	if err != nil {
		t.Fatal(err)
	}
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
	vr := videoResource(t, allwaysSteven)
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
	//t.Log("Body: ", rb, len(rb))
	if len(rb) != 255 {
		t.Fatal("Unexpected output from findAllVideos")
	}
	dao.DeleteAll()
}
