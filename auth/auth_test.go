package auth

import (
	"github.com/slspeek/go-restful"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestAuthFilterChallenges(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/videos", nil)
	recorder := httptest.NewRecorder()

	auth := NewBasicAuthenticator("../test_htpasswd")

	ws := new(restful.WebService)

	ws.Path("/api/videos")

	dummyTarget := func(req *restful.Request, resp *restful.Response) {}

	ws.Route(ws.GET("").Filter(FilterFactory(auth)).To(dummyTarget))

	container := restful.NewContainer()

	container.Add(ws)

	container.ServeHTTP(recorder, req)

	if recorder.Header().Get("WWW-Authenticate") != `Basic realm="gotube.org"` {
	  t.Log(recorder.Header().Get("WWW-Authenticate"))
		t.Fatal("Should challenge")
	}

}
func TestAuthFilterDoesNotChallenge(t *testing.T) {
	req, _ := http.NewRequest("GET", "/api/videos", nil)
	req.Header.Add("Do-Not-Challenge", "True")
	recorder := httptest.NewRecorder()

	auth := NewBasicAuthenticator("../test_htpasswd")

	ws := new(restful.WebService)

	ws.Path("/api/videos")

	dummyTarget := func(req *restful.Request, resp *restful.Response) {}

	ws.Route(ws.GET("").Filter(FilterFactory(auth)).To(dummyTarget))

	container := restful.NewContainer()

	container.Add(ws)

	container.ServeHTTP(recorder, req)

	if recorder.Header().Get("WWW-Authenticate") != "" {
		t.Fatal("Should not challenge")
	}

}
