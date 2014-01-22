package rest

import (
	"fmt"
	"github.com/slspeek/flowgo"
	"github.com/slspeek/go-restful"
	"github.com/slspeek/goblob"
	"github.com/slspeek/gotube/auth"
	"github.com/slspeek/gotube/common"
	"github.com/slspeek/gotube/mongo"
	"github.com/slspeek/gotube/thumb"
	"io"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
	"net/http"
	"strconv"
	"time"
)

type VideoResource struct {
	db          string
	sess        *mgo.Session
	videos      *mongo.Dao
	auth        *auth.Auth
	bs          *goblob.BlobService
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
	v.generateThumbs(*vid, 5, 128)
}

func NewVideoResource(sess *mgo.Session, db string, collecion string, auth *auth.Auth) *VideoResource {
	v := new(VideoResource)
	v.sess = sess.Copy()
	v.db = db
	v.auth = auth
	v.videos = mongo.NewDao(sess.Copy(), db, collecion)
	v.bs = goblob.NewBlobService(sess.Copy(), db, "flowfs")
	v.flowService = flow.NewUploadHandler(v.bs, v.writeBackBlobId)
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

func (v VideoResource) IsPublicFilter(param string) restful.FilterFunction {
	return func(req *restful.Request, resp *restful.Response, chain *restful.FilterChain) {
		id := req.PathParameter(param)
		vid := new(common.Video)
		err := v.videos.Get(id, &vid)
		if err != nil {
			resp.AddHeader("Content-Type", "text/plain")
			resp.WriteErrorString(http.StatusNotFound, "Video could not be found.")
			return
		}
		if ! vid.Public {
			resp.AddHeader("Content-Type", "text/plain")
			resp.WriteErrorString(http.StatusForbidden, "Not public")
			return
		}
		req.SetAttribute("video-object", vid)
		chain.ProcessFilter(req, resp)
	}
}

func (v VideoResource) Register(container *restful.Container) {
	videoIdFilter := ObjectIdFilter("video-id")
	ownerFilter := v.IsOwnerFilter("video-id")
	publicFilter := v.IsPublicFilter("video-id")
	apiWs := new(restful.WebService)
	contentWs := new(restful.WebService)
	publicWs := new(restful.WebService)
	publicContentWs := new(restful.WebService)

	publicWs.Path("/public/api/videos").Produces(restful.MIME_JSON, restful.MIME_XML)

  publicContentWs. 
		Path("/public/content/videos") /*.Produces("application/octet-stream")*/
	contentWs.
		Path("/content/videos") /*.Produces("application/octet-stream")*/
	apiWs.
		Path("/api/videos").
		Consumes(restful.MIME_XML, restful.MIME_JSON, "multipart/form-data").
		//Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML)

	apiWs.Route(apiWs.GET("").To(v.findAllVideos).
		// docs
		Doc("get all your videos").
		Writes(common.Video{})) // on the response

	publicWs.Route(publicWs.GET("").To(v.findPublicVideos).
		// docs
		Doc("get all public videos").
		Writes(common.Video{})) // on the response

	apiWs.Route(apiWs.GET("/playlist").To(v.playlist).
		// docs
		Doc("get a playlist for all your videos").
		Writes(common.Video{})) // on the response

	apiWs.Route(apiWs.GET("/{video-id}").Filter(videoIdFilter).Filter(ownerFilter).To(v.findVideo).
		// docs
		Doc("get a video").
		Param(apiWs.PathParameter("video-id", "identifier of the video").DataType("string")).
		Writes(common.Video{})) // on the response

	publicWs.Route(publicWs.GET("/{video-id}").Filter(videoIdFilter).Filter(publicFilter).To(v.findVideo).
		// docs
		Doc("get a public video").
		Param(publicWs.PathParameter("video-id", "identifier of the video").DataType("string")).
		Writes(common.Video{})) // on the response

	publicContentWs.Route(publicContentWs.GET("/{video-id}/thumbs/{thumb-id}").Filter(videoIdFilter).Filter(publicFilter).To(v.serveThumb).
		// docs
		Doc("get a thumb from video").
		Param(publicContentWs.PathParameter("video-id", "identifier of the video").DataType("string")).
		Param(publicContentWs.PathParameter("thumb-id", "identifier of the video").DataType("integer")))

	publicContentWs.Route(publicContentWs.GET("/{video-id}").Filter(videoIdFilter).Filter(publicFilter).To(v.serveVideo).
		// docs
		Doc("stream a video").
		Param(publicContentWs.PathParameter("video-id", "identifier of the video").DataType("string")))
	publicContentWs.Route(publicContentWs.GET("/{video-id}/download").Filter(videoIdFilter).Filter(publicFilter).To(v.downloadVideo).
		// docs
		Doc("download a video").
		Param(apiWs.PathParameter("video-id", "identifier of the video").DataType("string")))
	contentWs.Route(contentWs.GET("/{video-id}").Filter(videoIdFilter).Filter(ownerFilter).To(v.serveVideo).
		// docs
		Doc("stream a video").
		Param(contentWs.PathParameter("video-id", "identifier of the video").DataType("string")))

	contentWs.Route(contentWs.GET("/{video-id}/download").Filter(videoIdFilter).Filter(ownerFilter).To(v.downloadVideo).
		// docs
		Doc("download a video").
		Param(apiWs.PathParameter("video-id", "identifier of the video").DataType("string")))
	contentWs.Route(contentWs.GET("/{video-id}/thumbs/{thumb-id}").Filter(videoIdFilter).Filter(ownerFilter).To(v.serveThumb).
		// docs
		Doc("get a thumb from video").
		Param(contentWs.PathParameter("video-id", "identifier of the video").DataType("string")).
		Param(contentWs.PathParameter("thumb-id", "identifier of the video").DataType("integer")))

	apiWs.Route(apiWs.POST("").To(v.createVideo).
		// docs
		Doc("create a video").
		Reads(common.Video{})) // from the request

	apiWs.Route(apiWs.POST("/{video-id}/upload").To(v.uploadVideo).Doc("upload video content"))
	apiWs.Route(apiWs.GET("/{video-id}/upload").To(v.uploadVideo).Doc("check if we already have video content"))

	apiWs.Route(apiWs.PUT("/{video-id}").Filter(videoIdFilter).Filter(ownerFilter).To(v.updateVideo).
		// docs
		Doc("update a video").
		Param(apiWs.PathParameter("video-id", "identifier of the video").DataType("string")).
		Reads(common.Video{})) // from the request

	apiWs.Route(apiWs.DELETE("/{video-id}").Filter(videoIdFilter).Filter(ownerFilter).To(v.removeVideo).
		// docs
		Doc("delete a video").
		Param(apiWs.PathParameter("video-id", "identifier of the video").DataType("string")))

	apiWs.Filter(v.auth.Filter)
	contentWs.Filter(v.auth.Filter)

	container.Add(apiWs)
	container.Add(contentWs)
  container.Add(publicWs)
  container.Add(publicContentWs)
}

func (v VideoResource) findVideo(request *restful.Request, response *restful.Response) {
	vid := request.Attribute("video-object")
	response.WriteEntity(vid)
}

func (v VideoResource) findAllVideos(request *restful.Request, response *restful.Response) {
	user := request.Attribute("username")
	vid := make([]common.Video, 10)
	err := v.videos.Find(bson.M{"owner": user}, &vid, []string{"-created"})
	if err != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, "Error retrieving videos")
	} else {
		response.WriteEntity(vid)
	}
}

func (v VideoResource) findPublicVideos(request *restful.Request, response *restful.Response) {
	vid := make([]common.Video, 10)
	err := v.videos.Find(bson.M{"public": true}, &vid, []string{"-created"})
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
	err := v.videos.Find(bson.M{"owner": user}, &vid, []string{"Created-"})
	if err != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, "Error retrieving videos")
	} else {
		playlist := "<smil><body><seq>"
		for i := 0; i < len(vid); i++ {
			playlist += videoTag(vid[i], "http://"+request.Request.Host)
		}
		playlist += "</seq></body></smil>"
		fmt.Fprint(response.ResponseWriter, playlist)
	}
}

func videoTag(vid common.Video, domain string) string {
	return fmt.Sprintf(`<video src="%s/content/videos/%s" />`, domain, vid.Id.Hex())
}

func (v *VideoResource) createVideo(request *restful.Request, response *restful.Response) {
	vid := common.Video{Id: bson.NewObjectId()}
	err := request.ReadEntity(&vid)
	user := request.Attribute("username")
	if err == nil {
		if user == vid.Owner {
			vid.Created = time.Now()
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
		log.Print("Could not read back ", err)
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
	video := request.Attribute("video-object").(*common.Video)
	sess := v.sess.Copy()
	bs := goblob.NewBlobService(sess, v.db, "flowfs")
	defer bs.Close()
	if video.BlobId != "" {
		err := bs.Remove(video.BlobId)
		if err != nil {
			response.AddHeader("Content-Type", "text/plain")
			response.WriteErrorString(http.StatusInternalServerError, err.Error())
			return
		}
	}
	for _, blobId := range video.Thumbs {
		err := bs.Remove(blobId)
		if err != nil {
			response.AddHeader("Content-Type", "text/plain")
			response.WriteErrorString(http.StatusInternalServerError, err.Error())
			return
		}
	}
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

func (v *VideoResource) serveThumb(request *restful.Request, response *restful.Response) {
	video := request.Attribute("video-object").(*common.Video)
	thumbIndex := request.PathParameter("thumb-id")
	index, err := strconv.Atoi(thumbIndex)
	if err != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusBadRequest, "Thumb index could not be parsed")
		return
	}
	if index >= len(video.Thumbs) {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusBadRequest, "Thumb index out of bounds")
		return
	}
	var fid string
	fid = video.Thumbs[index]
	sess := v.sess.Copy()
	bs := goblob.NewBlobService(sess, v.db, "flowfs")
	defer bs.Close()
	blobHandle, err := bs.Open(fid)
	defer bs.Close()
	if err != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, "Thumb could not be opened")
		return
	}

	http.ServeContent(response.ResponseWriter, request.Request, "", blobHandle.UploadDate(), blobHandle)
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
	defer bs.Close()
	blobHandle, err := bs.Open(fid)
	defer bs.Close()
	if err != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, "Video could not be opened")
		return
	}

	http.ServeContent(response.ResponseWriter, request.Request, "", blobHandle.UploadDate(), blobHandle)
}

func (v *VideoResource) downloadVideo(request *restful.Request, response *restful.Response) {
	video := request.Attribute("video-object").(*common.Video)
	var fid string
	if fid = video.BlobId; fid == "" {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNoContent, "No Video blobid set yet")
		return
	}
	sess := v.sess.Copy()
	bs := goblob.NewBlobService(sess, v.db, "flowfs")
	defer bs.Close()
	blobHandle, err := bs.Open(fid)
	defer blobHandle.Close()
	if err != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, "Video could not be opened")
		return
	}

	response.AddHeader("Content-Type", "application/octet-stream")
	response.AddHeader("Content-Disposition", fmt.Sprintf(`attachment; filename="%s"`, blobHandle.Name()))
	_, err = io.Copy(response.ResponseWriter, blobHandle)
	if err != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, "Video could not be copied to the response")
		return
	}
}

func (v *VideoResource) generateThumbs(vid common.Video, count int, size int) (err error) {
	outputChan, err := thumb.CreateThumbs(v.bs, vid.BlobId, count, size)
	if err != nil {
		return
	}
	allThumbs := make(map[int]thumb.ThumbResult)
	for i := 0; i < count; i++ {
		result := <-outputChan
		allThumbs[result.Index] = result
	}
	for i := 0; i < count; i++ {
		if err = allThumbs[i].Err; err != nil {
			log.Printf("thumbnailing failed: %v", err)
			continue
		}
		vid.Thumbs = append(vid.Thumbs, allThumbs[i].BlobId)
	}
	err = v.videos.Update(vid.Id.Hex(), vid)
	if err != nil {
		log.Printf("Could not save thumbids in Video: %#v", err)
	}
	return
}
