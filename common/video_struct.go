package common

import (
	"labix.org/v2/mgo/bson"
	"time"
)

type Video struct {
  Id bson.ObjectId `_id,omitempty`
  Created time.Time
	Owner, Name, Desc, BlobId string
  Thumbs  []string
}
