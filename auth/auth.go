package auth

import (
	"encoding/json"
	auth "github.com/abbot/go-http-auth"
	"github.com/slspeek/go-restful"
	"labix.org/v2/mgo/bson"
	"net/http"
)

type CheckAuthFunc func(*http.Request) string

type Auth struct {
	Check CheckAuthFunc
}

func (self *Auth) Filter(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
	username := self.Check(req.Request)
	if username != "" {
		req.SetAttribute("username", username)
		chain.ProcessFilter(req, resp)
	} else {
    resp.AddHeader("WWW-Authenticate", "Basic realm=gotube.org")
		resp.WriteErrorString(http.StatusUnauthorized, "No username provided")
	}
}

func Authenticator(passwdFile string) *auth.BasicAuth {
	provider := auth.HtpasswdFileProvider(passwdFile)
	return auth.NewBasicAuthenticator("gotube.org", provider)
}

func AuthService(w http.ResponseWriter, r *auth.AuthenticatedRequest) {
	userinfo := make(bson.M)
	userinfo["username"] = r.Username
	b, err := json.Marshal(userinfo)
	if err != nil {
		http.Error(w, "Unable to marshal", http.StatusInternalServerError)
		return
	}
	w.Write(b)
}
