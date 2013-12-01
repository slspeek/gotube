package rest

import (
	"github.com/slspeek/go-restful"
	"github.com/slspeek/gotube/auth"
	"github.com/slspeek/gotube/mongo"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
	"net/http"
)

type Video struct {
	Id                        bson.ObjectId "_id,omitempty"
	Owner, Name, Desc, BlobId string
}

type VideoResource struct {
	videos *mongo.Dao
	auth   *auth.Auth
}

func NewVideoResource(sess *mgo.Session, db string, collecion string, auth *auth.Auth) *VideoResource {
	dao := mongo.NewDao(sess, db, collecion)
	return &VideoResource{dao, auth}
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
		vid := new(Video)
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
		chain.ProcessFilter(req, resp)
	}
}

func (v VideoResource) Register(container *restful.Container) {
	videoIdFilter := ObjectIdFilter("video-id")
	ownerFilter := v.IsOwnerFilter("video-id")
	ws := new(restful.WebService)
	ws.
		Path("/api/videos").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML)

	ws.Route(ws.GET("").To(v.findAllVideos).
		// docs
		Doc("get all your videos").
		Writes(Video{})) // on the response

	ws.Route(ws.GET("/{video-id}").Filter(videoIdFilter).Filter(ownerFilter).To(v.findVideo).
		// docs
		Doc("get a video").
		Param(ws.PathParameter("video-id", "identifier of the video").DataType("string")).
		Writes(Video{})) // on the response

	ws.Route(ws.POST("").To(v.createVideo).
		// docs
		Doc("create a video").
		Reads(Video{})) // from the request

	ws.Route(ws.PUT("/{video-id}").Filter(videoIdFilter).Filter(ownerFilter).To(v.updateVideo).
		// docs
		Doc("update a video").
		Param(ws.PathParameter("video-id", "identifier of the video").DataType("string")).
		Reads(Video{})) // from the request

	ws.Route(ws.DELETE("/{video-id}").Filter(videoIdFilter).Filter(ownerFilter).To(v.removeVideo).
		// docs
		Doc("delete a video").
		Param(ws.PathParameter("video-id", "identifier of the video").DataType("string")))

	ws.Filter(v.auth.Filter)
	container.Add(ws)
}

func (v VideoResource) findVideo(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("video-id")
	if !bson.IsObjectIdHex(id) {
		response.WriteErrorString(http.StatusBadRequest, "No valid object id passed")
		return
	}
	user := request.Attribute("username")
	vid := new(Video)
	err := v.videos.Get(id, &vid)
	if err != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotFound, "Video could not be found.")
	} else {
		if user == vid.Owner {
			response.WriteEntity(vid)
		} else {
			response.WriteErrorString(http.StatusForbidden, "Not the owner")
		}
	}
}

func (v VideoResource) findAllVideos(request *restful.Request, response *restful.Response) {
	user := request.Attribute("username")
	vid := make([]Video, 10)
	err := v.videos.Find(bson.M{"owner": user}, &vid)
	if err != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, "Error retrieving videos")
	} else {
		response.WriteEntity(vid)
	}
}

func (v *VideoResource) createVideo(request *restful.Request, response *restful.Response) {
	vid := Video{Id: bson.NewObjectId()}
	err := request.ReadEntity(&vid)
	user := request.Attribute("username")
	log.Print("Create video called")
	if err == nil {
		if user == vid.Owner {
			id, err := v.videos.Create(vid)
			log.Print("After create video called", err)
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
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
	}
}

func (v *VideoResource) updateVideo(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("video-id")
	if !bson.IsObjectIdHex(id) {
		response.WriteErrorString(http.StatusBadRequest, "No valid object id passed")
		return
	}
	vid := &Video{Id: bson.ObjectIdHex(id)}
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
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
	}
}

func (v *VideoResource) removeVideo(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("video-id")
	if !bson.IsObjectIdHex(id) {
		response.WriteErrorString(http.StatusBadRequest, "No valid object id passed")
		return
	}
	user := request.Attribute("username")
	vid := new(Video)
	err := v.videos.Get(id, &vid)
	if err != nil {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
		return
	}
	if user == vid.Owner {
		err := v.videos.Delete(id)
		if err != nil {
			response.AddHeader("Content-Type", "text/plain")
			response.WriteErrorString(http.StatusInternalServerError, err.Error())
		}
	} else {
		response.WriteErrorString(http.StatusForbidden, "Not the owner")
	}
}
