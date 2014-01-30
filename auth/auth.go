package auth

import (
	"encoding/json"
	auth "github.com/abbot/go-http-auth"
	"github.com/slspeek/go-restful"
	"labix.org/v2/mgo/bson"
	"net/http"
)

func NewBasicAuthenticator(passwdFile string) GeneralAuth {
	provider := auth.HtpasswdFileProvider(passwdFile)
	return auth.NewBasicAuthenticator("gotube.org", provider)
}

type DigestAuthenticator struct {
	Digest *auth.DigestAuth
}

func (d DigestAuthenticator) CheckAuth(r *http.Request) (username string) {
	username, _ = d.Digest.CheckAuth(r)
	return
}

func (d DigestAuthenticator) RequireAuth(w http.ResponseWriter, r *http.Request) {
  d.Digest.RequireAuth(w,r)
}

func NewDigestAuthenticator(passwdFile string) GeneralAuth {
	provider := auth.HtdigestFileProvider(passwdFile)
	di := DigestAuthenticator{auth.NewDigestAuthenticator("gotube.org", provider)}
	return di
}

type GeneralAuth interface {
	CheckAuth(*http.Request) string
	RequireAuth(http.ResponseWriter, *http.Request)
}

func ChallengeOptionally(g GeneralAuth, w http.ResponseWriter, r *http.Request) {
  //w.Header().Set("Last-Modified", time.Now().Format(time.RFC1123))
	if r.Header.Get("Do-Not-Challenge") == "True" {
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
	} else {
		g.RequireAuth(w, r)
	}
}

func AuthServiceFactory(g GeneralAuth) (auth func(w http.ResponseWriter, r *http.Request)) {
	auth = func(w http.ResponseWriter, r *http.Request) {
		userinfo := make(bson.M)
		username := g.CheckAuth(r)
		if username != "" {
			userinfo["username"] = username
			b, err := json.Marshal(userinfo)
			if err != nil {
				http.Error(w, "Unable to marshal", http.StatusInternalServerError)
				return
			}
			w.Write(b)
		} else {
			ChallengeOptionally(g, w, r)
		}
	}
	return
}

func FilterFactory(g GeneralAuth) (filter func(req *restful.Request, resp *restful.Response, chain *restful.FilterChain)) {
	filter = func(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
		username := g.CheckAuth(req.Request)
		if username != "" {
			req.SetAttribute("username", username)
			chain.ProcessFilter(req, resp)
		} else {
			ChallengeOptionally(g, resp.ResponseWriter, req.Request)
		}
	}
	return
}

