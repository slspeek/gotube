package common

import (
	"labix.org/v2/mgo/bson"
)

type Video struct {
  Id bson.ObjectId `_id,omitempty`
	Owner, Name, Desc, BlobId string
}
