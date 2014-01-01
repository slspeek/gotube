package main

import (
	"flag"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/slspeek/go-restful"
	"github.com/slspeek/gotube/auth"
	"github.com/slspeek/gotube/rest"
	"labix.org/v2/mgo"
	"log"
	"net/http"
	"os"
	"syscall"
)

var (
	sess *mgo.Session
)

func check(err error) {
	if err != nil {
		panic(err)
	}
}

var noSsl = flag.Bool("no-ssl", false, "run without SSL")
var dbName = flag.String("db", "gotube", "name of the mongodb to use")
var pwFile = flag.String("pw", "htpasswd", "path of the html password file")

func main() {
	flag.Parse()
	var err error
	sess, err = mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	f, err := os.Create("gotube.pid")
	if err != nil {
		panic(err)
	}
	fmt.Fprintln(f, syscall.Getpid())
  f.Close()

	r := mux.NewRouter()

	authenticator := auth.Authenticator(*pwFile)
	videoService := rest.NewVideoResource(sess.Copy(), *dbName, "Video", &auth.Auth{authenticator.CheckAuth})
	container := restful.NewContainer()
	videoService.Register(container)
	container.Filter(container.OPTIONSFilter)

	r.HandleFunc("/auth", authenticator.Wrap(auth.AuthService))

	r.PathPrefix("/api/videos").Handler(container)
	r.PathPrefix("/content/videos").Handler(container)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir(flag.Arg(0))))
	if !*noSsl {
		log.Fatal(http.ListenAndServeTLS(":8080", "cert.pem", "key.pem", r))
	} else {
		log.Fatal(http.ListenAndServe(":8080", r))
	}
}
