package main

import (
	"github.com/gorilla/mux"
	"github.com/slspeek/go-restful"
	"github.com/slspeek/gotube/auth"
	"github.com/slspeek/gotube/rest"
	"labix.org/v2/mgo"
	"log"
	"net/http"
	"os"
)

var (
	sess *mgo.Session
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	var err error
	sess, err = mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()

	authenticator := auth.Authenticator("htpasswd")
	videoService := rest.NewVideoResource(sess.Copy(), "gotube", "Video", &auth.Auth{authenticator.CheckAuth})
	container := restful.NewContainer()
	videoService.Register(container)
	container.Filter(container.OPTIONSFilter)

	r.HandleFunc("/auth", authenticator.Wrap(auth.AuthService))

	r.PathPrefix("/api/videos").Handler(container)
	r.PathPrefix("/content/videos").Handler(container)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir(os.Args[1])))
	//log.Fatal(http.ListenAndServeTLS(":8080", "cert.pem", "key.pem", r))
	log.Fatal(http.ListenAndServe(":8080", r))
}
