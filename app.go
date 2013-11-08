package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/slspeek/crudapi"
	"github.com/slspeek/crudapi/storage/mongo"
	"github.com/slspeek/go_httpauth"
	"github.com/slspeek/goblob"
	"github.com/slspeek/flowgo"
	"html/template"
	"io"
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

func storage() *mongo.MongoStorage {
	s, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	return mongo.NewMongoStorage(s, "test")
}

func uploadHandler() *flow.UploadHandler {
  bs, err := goblob.NewBlobService("localhost", "test", "flow")
  if err != nil {
    panic("No blobservice could be created")
  }
	return flow.NewUploadHandler(bs, func(w http.ResponseWriter, r *http.Request, blobId string) {
    log.Println("Upload completed: ", blobId)
  })
}

var rootTemplate = template.Must(template.New("root").Parse(rootTemplateHTML))

const rootTemplateHTML = `
<html><body>
<form action="{{.}}" method="POST" enctype="multipart/form-data">
Upload File: <input type="file" name="file"><br>
<input type="submit" name="submit" value="Submit">
</form></body></html>
`

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func handleRoot(w http.ResponseWriter, r *http.Request) {
	uploadURL := "/upload"
	w.Header().Set("Content-Type", "text/html")
	err := rootTemplate.Execute(w, uploadURL)
	if err != nil {
		log.Printf("%v", err)
	}
}

func handleServe(w http.ResponseWriter, r *http.Request) {
	filename := r.FormValue("filename")
	bs, err := goblob.NewBlobService("localhost", "test", "fs")
	check(err)
	gf, err := bs.OpenName(filename)
	io.Copy(w, gf)
	gf.Close()
	bs.Close()
}

func handleUpload(w http.ResponseWriter, r *http.Request) {
	log.Print("Upload ", r)
	//reader, err := r.MultipartReader()
	//if err != nil {
	//fmt.Println(err)
	//http.Error(w, "not a form", http.StatusBadRequest)
	//}
	f, _, err := r.FormFile("file")
	if err != nil {
		fmt.Println(err)
		http.Error(w, "not a form", http.StatusBadRequest)
	}
	err = r.ParseForm()
	if err != nil {
		fmt.Println(err)
		http.Error(w, "not a form", http.StatusBadRequest)
	}

	fmt.Println("Form: ", r.Form)
	defer r.Body.Close()

	bs, err := goblob.NewBlobService("localhost", "test", "fs")
	if err != nil {
		panic(err)
	}
	gf, err := bs.Create(r.URL.Query().Get("flowFilename"))
	if err != nil {
		panic(err)
	}
	_, err = io.Copy(gf, f)
	if err != nil {
		panic(err)
	}
	gf.Close()
	bs.Close()
	//http.Redirect(w, r, "/serve/?filename="+filename, http.StatusFound)
}

func main() {

	// storage
	s := storage()
	// router
	r := mux.NewRouter()

	r.HandleFunc("/form", handleRoot)
	r.HandleFunc("/serve", handleServe)
	r.Handle("/upload", uploadHandler())
	// mounting the API
	crudapi.MountAPI(r.PathPrefix("/api").Subrouter(), s, NewMyGuard())

	r.PathPrefix("/").Handler(http.FileServer(http.Dir(os.Args[1])))
	log.Fatal(http.ListenAndServe(":8080", r))
}
