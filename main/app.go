package main

import (
	ba "github.com/abbot/go-http-auth"
	"github.com/gorilla/mux"
	"github.com/slspeek/flowgo"
	"github.com/slspeek/go-restful"
	"github.com/slspeek/goblob"
	"github.com/slspeek/gotube/auth"
	"github.com/slspeek/gotube/mongo"
	"github.com/slspeek/gotube/rest"
	"labix.org/v2/mgo"
	"log"
	"net/http"
	"os"
)

var (
	sess *mgo.Session
)

func uploadHandler() *flow.UploadHandler {
	bs := goblob.NewBlobService(sess, "test", "flow")
	return flow.NewUploadHandler(bs, func(r *http.Request, blobId string) {
		log.Println("Upload completed: ", blobId)
		id := r.URL.Path[8:]
		vid := new(rest.Video)
		dao := mongo.NewDao(sess, "test", "Video")
		err := dao.Get(id, vid)
		if err != nil {
			log.Println("Could not get video : ", id)
		}
		vid.BlobId = blobId
		err = dao.Update(id, vid)
		if err != nil {
			log.Println("Could not update video : ", id)
		}

	})
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func handleServe(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Path[7:]
	sess := sess.Copy()
	bs := goblob.NewBlobService(sess, "test", "flow")
	gf, err := bs.Open(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}
	http.ServeContent(w, r, "", gf.UploadDate(), gf)
	gf.Close()
	bs.Close()
	sess.Close()
}

func main() {
	var err error
	sess, err = mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}

	r := mux.NewRouter()

	authenticator := auth.Authenticator("htpasswd")
	videoService := rest.NewVideoResource(sess.Copy(), "test", "Video", &auth.Auth{authenticator.CheckAuth})
	container := restful.NewContainer()
	videoService.Register(container)
	container.Filter(container.OPTIONSFilter)

	r.HandleFunc("/serve/{video}", ba.JustCheck(authenticator, handleServe))
	r.HandleFunc("/auth", authenticator.Wrap(auth.AuthService))

	uph := uploadHandler()
	r.Handle("/upload/{video}", uph)

	r.PathPrefix("/api/videos").Handler(container)
	r.PathPrefix("/").Handler(http.FileServer(http.Dir(os.Args[1])))
	//log.Fatal(http.ListenAndServeTLS(":8080", "cert.pem", "key.pem", r))
  log.Fatal(http.ListenAndServe(":8080", r))
}
