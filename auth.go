package main

import (
	"encoding/json"
	auth "github.com/abbot/go-http-auth"
	"labix.org/v2/mgo/bson"
	"log"
	"net/http"
)

func authenticator() *auth.BasicAuth {
	provider := auth.HtpasswdFileProvider("htpasswd")
  log.Println("password file: ", provider)
	return auth.NewBasicAuthenticator("gotube.org", provider)
}

func authService(w http.ResponseWriter, r *auth.AuthenticatedRequest) {
	userinfo := make(bson.M)
	userinfo["username"] = r.Username
	b, err := json.Marshal(userinfo)
	if err != nil {
		http.Error(w, "Unable to marshal", http.StatusInternalServerError)
		return
	}
	w.Write(b)
}
