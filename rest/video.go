package rest

import (
	"fmt"
	"github.com/slspeek/flowgo"
	"github.com/slspeek/go-restful"
	"github.com/slspeek/goblob"
	"github.com/slspeek/gotube/auth"
	"github.com/slspeek/gotube/common"
	"github.com/slspeek/gotube/mongo"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
	"net/http"
)


type VideoResource struct {
	db          string
	sess        *mgo.Session
	videos      *mongo.Dao
	auth        *auth.Auth
	flowService *flow.UploadHandler
}

func (v *VideoResource) writeBackBlobId(r *http.Request, blobId string) {
	id := r.URL.Path[12:36]
	vid := new(common.Video)
	err := v.videos.Get(id, vid)
	if err != nil {
		log.Println("Could not get video : ", id)
	}
	vid.BlobId = blobId
	err = v.videos.Update(id, vid)
	if err != nil {
		log.Println("Could not update video : ", id)
	}
	log.Println("Updated video: ", id, " with BlobId: ", blobId)
}

func NewVideoResource(sess *mgo.Session, db string, collecion string, auth *auth.Auth) *VideoResource {
	v := new(VideoResource)
	v.sess = sess.Copy()
	v.db = db
	v.auth = auth
	v.videos = mongo.NewDao(sess.Copy(), db, collecion)
  blobs := goblob.NewBlobService(sess.Copy(), db, "flowfs")
	v.flowService = flow.NewUploadHandler(blobs, v.writeBackBlobId)
	return v
}

func ObjectIdFilter(param string) restful.FilterFunction {
	return func(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
		id := req.PathParameter(param)
		if !bson.IsObjectIdHex(id) {
			resp.WriteErrorString(http.StatusBadRequest, "No valid object id passed")
			return
		}
		chain.ProcessFilter(req, resp)
	}
}

func (v VideoResource) IsOwnerFilter(param string) restful.FilterFunction {
	return func(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
		id := req.PathParameter(param)
		user := req.Attribute("username")
		vid := new(common.Video)
		err := v.videos.Get(id, &vid)
		if err != nil {
			resp.AddHeader("Content-Type", "text/plain")
			resp.WriteErrorString(http.StatusNotFound, "Video could not be found.")
			return
		}
		if user != vid.Owner {
			resp.AddHeader("Content-Type", "text/plain")
			resp.WriteErrorString(http.StatusForbidden, "Not the owner")
			return
		}
		req.SetAttribute("video-object", vid)
		chain.ProcessFilter(req, resp)
	}
}

func (v VideoResource) Register(container *restful.Container) {
	videoIdFilter := ObjectIdFilter("video-id")
	ownerFilter := v.IsOwnerFilter("video-id")
	ws := new(restful.WebService)
	ws2 := new(restful.WebService)
	ws2.
		Path("/content/videos") /*.Produces("application/octet-stream")*/
	ws.
		Path("/api/videos").
		Consumes(restful.MIME_XML, restful.MIME_JSON, "multipart/form-data").
		//Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML)

	ws.Route(ws.GET("/playlist").To(v.playlist).
		// docs
		Doc("get a playlist for all your videos").
		Writes(common.Video{})) // on the response
	ws.Route(ws.GET("").To(v.findAllVideos).
		// docs
		Doc("get all your videos").
		Writes(common.Video{})) // on the response

	ws.Route(ws.GET("/{video-id}").Filter(videoIdFilter).Filter(ownerFilter).To(v.findVideo).
		// docs
		Doc("get a video").
		Param(ws.PathParameter("video-id", "identifier of the video").DataType("string")).
		Writes(common.Video{})) // on the response

	ws2.Route(ws.GET("/{video-id}").Filter(videoIdFilter).Filter(ownerFilter).To(v.serveVideo).
		// docs
		Doc("download a video").
		Param(ws.PathParameter("video-id", "identifier of the video").DataType("string")))

	ws.Route(ws.POST("").To(v.createVideo).
		// docs
		Doc("create a video").
		Reads(common.Video{})) // from the request

	ws.Route(ws.POST("/{video-id}/upload").To(v.uploadVideo).Doc("upload video content"))

	ws.Route(ws.PUT("/{video-id}").Filter(videoIdFilter).Filter(ownerFilter).To(v.updateVideo).
		// docs
		Doc("update a video").
		Param(ws.PathParameter("video-id", "identifier of the video").DataType("string")).
		Reads(common.Video{})) // from the request

	ws.Route(ws.DELETE("/{video-id}").Filter(videoIdFilter).Filter(ownerFilter).To(v.removeVideo).
		// docs
		Doc("delete a video").
		Param(ws.PathParameter("video-id", "identifier of the video").DataType("string")))

	ws.Filter(v.auth.Filter)
	ws2.Filter(v.auth.Filter)

	container.Add(ws)
	container.Add(ws2)
}

func (v VideoResource) findVideo(request *restful.Request, response *restful.Response) {
	vid := request.Attribute("video-object")
	response.WriteEntity(vid)
}

func (v VideoResource) findAllVideos(request *restful.Request, response *restful.Response) {
	user := request.Attribute("username")
	vid := make([]common.Video, 10)
	err := v.videos.Find(bson.M{"owner": user}, &vid)
	if err != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, "Error retrieving videos")
	} else {
		response.WriteEntity(vid)
	}
}


func (v VideoResource) playlist(request *restful.Request, response *restful.Response) {
	user := request.Attribute("username")
	vid := make([]common.Video, 10)
	err := v.videos.Find(bson.M{"owner": user}, &vid)
	if err != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, "Error retrieving videos")
	} else {
    log.Printf("URL: %#v ", request.Request.RequestURI)
    log.Printf("Host: %#v ", request.Request.Host)
    playlist := "<smil><body><seq>"
    for i :=0; i< len(vid); i++ {
      playlist += videoTag(vid[i], "http://" + request.Request.Host)
    }



    playlist += "</seq></body></smil>"
    fmt.Fprint(response.ResponseWriter, playlist)
	}
}

func videoTag(vid common.Video, domain string) string {
  return fmt.Sprintf(`<video src="%s/content/videos/%s" />`,domain, vid.Id.Hex())
}

func (v *VideoResource) createVideo(request *restful.Request, response *restful.Response) {
	vid := common.Video{Id: bson.NewObjectId()}
	err := request.ReadEntity(&vid)
	user := request.Attribute("username")
	if err == nil {
		if user == vid.Owner {
			id, err := v.videos.Create(vid)
			if err != nil {
				response.AddHeader("Content-Type", "text/plain")
				response.WriteErrorString(http.StatusInternalServerError, err.Error())
				return
			}
			vid.Id = bson.ObjectIdHex(id)
			response.WriteHeader(http.StatusCreated)
			response.WriteEntity(vid)
		} else {
			response.WriteErrorString(http.StatusForbidden, "Not the owner")
			return
		}
	} else {
		log.Print("Could not read back", err)
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusBadRequest, err.Error())
		return
	}
}

func (v *VideoResource) updateVideo(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("video-id")
	vid := &common.Video{Id: bson.ObjectIdHex(id)}
	user := request.Attribute("username")
	err := request.ReadEntity(&vid)
	if err == nil {
		if user == vid.Owner {
			err := v.videos.Update(vid.Id.Hex(), vid)
			if err != nil {
				response.AddHeader("Content-Type", "text/plain")
				response.WriteErrorString(http.StatusInternalServerError, err.Error())
				return
			}
			response.WriteEntity(vid)
		} else {
			response.WriteErrorString(http.StatusForbidden, "Not the owner")
		}
	} else {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusBadRequest, err.Error())
		return
	}
}

func (v *VideoResource) removeVideo(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("video-id")
	err := v.videos.Delete(id)
	if err != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
		return
	}
}

func (v *VideoResource) uploadVideo(request *restful.Request, response *restful.Response) {
	v.flowService.ServeHTTP(response.ResponseWriter, request.Request)
}

func (v *VideoResource) serveVideo(request *restful.Request, response *restful.Response) {
	video := request.Attribute("video-object").(*common.Video)
	var fid string
	if fid = video.BlobId; fid == "" {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNoContent, "No Video blobid set yet")
		return
	}
	sess := v.sess.Copy()
	bs := goblob.NewBlobService(sess, v.db, "flowfs")
	blobHandle, err := bs.Open(fid)
	if err != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, "Video could not be opened")
		return
	}

	http.ServeContent(response.ResponseWriter, request.Request, "", blobHandle.UploadDate(), blobHandle)
	blobHandle.Close()
	bs.Close()
}
