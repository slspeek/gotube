package main

import (
	"encoding/json"
	auth "github.com/abbot/go-http-auth"
	"labix.org/v2/mgo/bson"
	"net/http"
)

func secret(user, realm string) string {
	if user == "steven" {
		// password is "gnu"
		return "83191b52807a4c2ccf91492a0944580b"
	}
	return ""
}

func authenticator() *auth.DigestAuth {
	return auth.NewDigestAuthenticator("gotube.org", secret)
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

