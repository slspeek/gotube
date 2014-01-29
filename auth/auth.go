package auth

import (
	"encoding/json"
	auth "github.com/abbot/go-http-auth"
	"github.com/slspeek/go-restful"
	"labix.org/v2/mgo/bson"
	"net/http"
)


func NewAuthenticator(passwdFile string) Authenticator {
	provider := auth.HtpasswdFileProvider(passwdFile)
	return Authenticator{auth.NewBasicAuthenticator("gotube.org", provider)}
}

type Authenticator struct {
	*auth.BasicAuth
}

func (auth *Authenticator) AuthService(w http.ResponseWriter, r *http.Request) {
	userinfo := make(bson.M)
	username := auth.CheckAuth(r)
	if username != "" {
		userinfo["username"] = username
		b, err := json.Marshal(userinfo)
		if err != nil {
			http.Error(w, "Unable to marshal", http.StatusInternalServerError)
			return
		}
		w.Write(b)
	} else {
		if r.Header.Get("Do-Not-Challenge") != "True" {
			w.Header().Add("WWW-Authenticate", "Basic realm=gotube.org")
		}
		http.Error(w, "Not authenticated", http.StatusUnauthorized)
	}
}

func (self *Authenticator) Filter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	username := self.CheckAuth(req.Request)
	if username != "" {
		req.SetAttribute("username", username)
		chain.ProcessFilter(req, resp)
	} else {
		if req.Request.Header.Get("Do-Not-Challenge") != "True" {
			resp.AddHeader("WWW-Authenticate", "Basic realm=gotube.org")
		}
		resp.WriteErrorString(http.StatusUnauthorized, "No username provided")
	}
}
