package main

import (
	auth "github.com/abbot/go-http-auth"
	"github.com/gorilla/mux"
	"github.com/slspeek/crudapi"
	"github.com/slspeek/flowgo"
	"github.com/slspeek/go-restful"
	"github.com/slspeek/go_httpauth"
	"github.com/slspeek/goblob"
	"labix.org/v2/mgo"
	"log"
	"net/http"
	"os"
)

type MyGuard struct {
	BasicAuth *go_httpauth.BasicServer
}

func NewMyGuard() MyGuard {
	return MyGuard{go_httpauth.NewBasic("MyRealm", func(user string, realm string) string {
		// Replace this with a real lookup function here
		return "ape"
	})}
}

func (g MyGuard) Authenticate(resp http.ResponseWriter, req *http.Request) (bool, string, string) {
	auth, username := g.BasicAuth.Auth(resp, req)
	log.Println("allowed:", auth, "username:", username)
	return auth, username, ""
}

func (g MyGuard) Authorize(client string, action crudapi.Action, urlVars map[string]string) (ok bool, errorMessage string) {
	log.Println("urlVars:", urlVars, "client:", client, "Action:", action)
	return true, ""
}

func hello(resp http.ResponseWriter, req *http.Request) {
	resp.Write([]byte("Hello there!"))
}

func uploadHandler() *flow.UploadHandler {
	bs, err := goblob.NewBlobService("localhost", "test", "flow")
	if err != nil {
		panic("No blobservice could be created")
	}
	return flow.NewUploadHandler(bs, func(r *http.Request, blobId string) {
		log.Println("Upload completed: ", blobId)
	})
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func handleServe(w http.ResponseWriter, r *http.Request) {
	filename := r.FormValue("filename")
	bs, err := goblob.NewBlobService("localhost", "test", "flow")
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	gf, err := bs.OpenName(filename)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	http.ServeContent(w, r, filename, gf.UploadDate(), gf)
	gf.Close()
	bs.Close()
}

func main() {

	// router

	r := mux.NewRouter()

	sess, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	authenticator := authenticator("htpasswd")
	videoService := NewVideoResource(sess, "test", "Video", &Auth{authenticator.CheckAuth})
	container := restful.NewContainer()
	videoService.Register(container)

	r.HandleFunc("/serve", auth.JustCheck(authenticator, handleServe))
	r.HandleFunc("/auth", authenticator.Wrap(authService))
	uph := uploadHandler()

	r.HandleFunc("/upload", auth.JustCheck(authenticator, uph.ServeHTTP))

	r.PathPrefix("/api/videos").Handler(container)
	r.PathPrefix("/API").HandlerFunc(container.ServeHTTP)
	r.HandleFunc("/apis", container.ServeHTTP)
	log.Println("apis registered")
	r.PathPrefix("/").Handler(http.FileServer(http.Dir(os.Args[1])))
	log.Fatal(http.ListenAndServeTLS(":8080", "cert.pem", "key.pem", r))
}
