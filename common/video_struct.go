package common

import (
	"labix.org/v2/mgo/bson"
	"time"
)

type Video struct {
  Id bson.ObjectId `_id,omitempty`
  Created time.Time
	Owner, Name, Desc, BlobId string
  Public bool
  Thumbs  []string
}
// REST Client Video 
type CVideo struct {
  Id string
  Created time.Time
	Owner, Name, Desc string
  Public bool
  Stream  string
  Download string
  Thumbs  []string
}
