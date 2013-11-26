package main

import (
	"github.com/slspeek/go-restful"
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"log"
	"net/http"
)

//CREATE AUTH, user becomes owner
//UPDATE AUTH user has to be owner
//GET_ONE AUTH user has to be owner
//GET_ALL ?PUBLIC later
//DELETE_ONE AUTH user has to be owner
//DELETE_ALL?

//GET_ALL_OWN AUTH

type Video struct {
  Id   bson.ObjectId   "_id,omitempty"
	Owner, Name, Descr, BlobId string
}

func (v Video) SetId(id string) {
  v.Id = bson.ObjectIdHex(id)
}

type VideoResource struct {
	// normally one would use DAO (data access object)
	videos *MongoDao
}

func NewVideoResource(sess *mgo.Session, db string, collecion string) *VideoResource {
	dao := NewMongoDao(sess, db, collecion)
	return &VideoResource{dao}
}

func (v VideoResource) Register(container *restful.Container) {
	ws := new(restful.WebService)
	ws.
		Path("/videos").
		Consumes(restful.MIME_XML, restful.MIME_JSON).
		Produces(restful.MIME_JSON, restful.MIME_XML) // you can specify this per route as well

	ws.Route(ws.GET("/{video-id}").To(v.findVideo).
		// docs
		Doc("get a video").
		Param(ws.PathParameter("video-id", "identifier of the video").DataType("string")).
		Writes(Video{})) // on the response

	ws.Route(ws.POST("").To(v.updateVideo).
		// docs
		Doc("create a video").
		Reads(Video{})) // from the request

	ws.Route(ws.PUT("/{video-id}").To(v.createVideo).
		// docs
		Doc("update a video").
		Param(ws.PathParameter("video-id", "identifier of the video").DataType("string")).
		Reads(Video{})) // from the request

	ws.Route(ws.DELETE("/{video-id}").To(v.removeVideo).
		// docs
		Doc("delete a video").
		Param(ws.PathParameter("video-id", "identifier of the video").DataType("string")))

	container.Add(ws)
}

func (v VideoResource) findVideo(request *restful.Request, response *restful.Response) {
	log.Print("find")
	id := request.PathParameter("video-id")
	user := request.Attribute("username")
  log.Print("User: ", user)
	vid := new(Video)
  err := v.videos.Get(id, &vid)
	if err !=nil || len(vid.Id) == 0 {
    log.Print("Error: ", err)
    log.Printf("Video: %#v", vid)
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

func (v *VideoResource) updateVideo(request *restful.Request, response *restful.Response) {
	vid := new(Video)
	err := request.ReadEntity(&vid)
	user := request.Attribute("username")
	if err == nil {
		if user == vid.Owner {
			v.videos.Update(vid.Id.Hex(), vid)
			response.WriteEntity(vid)
		} else {
			response.WriteErrorString(http.StatusForbidden, "Not the owner")
		}
	} else {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
	}
}

func (v *VideoResource) createVideo(request *restful.Request, response *restful.Response) {
	vid := Video{Id: bson.ObjectIdHex(request.PathParameter("video-id"))}
	user := request.Attribute("username")
  log.Print("User: ", user)
	err := request.ReadEntity(&vid)
	if err == nil {
		if user == vid.Owner {
			id, err := v.videos.Create(vid)
			if err != nil {
				response.AddHeader("Content-Type", "text/plain")
				response.WriteErrorString(http.StatusInternalServerError, err.Error())
			}
			vid.Id = bson.ObjectIdHex(id)
			response.WriteHeader(http.StatusCreated)
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
