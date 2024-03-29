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
var webclientDir = flag.String("www", "src/github.com/slspeek/gotube/web/app", "path of the web client")
var port = flag.Int("port", 8080, "webserver port")
var digest = flag.Bool("digest", false, "run with digest authentication in stead of basic")

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

	var authenticator auth.GeneralAuth
	if *digest {
		authenticator = auth.NewDigestAuthenticator(*pwFile)
	} else {
		authenticator = auth.NewBasicAuthenticator(*pwFile)
	}
	videoService := rest.NewVideoResource(sess.Copy(), *dbName, "Video", auth.FilterFactory(authenticator))
	container := restful.NewContainer()
	videoService.Register(container)
	container.Filter(container.OPTIONSFilter)

	r.HandleFunc("/auth", auth.AuthServiceFactory(authenticator))

	r.PathPrefix("/api/videos").Handler(container)
	r.PathPrefix("/public/api/videos").Handler(container)
	r.PathPrefix("/public/content/videos").Handler(container)
	r.PathPrefix("/content/videos").Handler(container)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir(*webclientDir)))
	listenAddress := fmt.Sprintf(":%d", *port)
	if !*noSsl {
		log.Fatal(http.ListenAndServeTLS(listenAddress, "cert.pem", "key.pem", r))
	} else {
		log.Fatal(http.ListenAndServe(listenAddress, r))
	}
}
