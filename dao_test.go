package main

import (
  "labix.org/v2/mgo"
  "testing"
)



func TestMongoDao(t *testing.T) {
  sess, err := mgo.Dial("localhost")
  if err != nil {
    t.Fatal(err)
  }
  dao := NewMongoDao(sess, "test", "Video")
  v1 := Video{"", "steven", "Novocento", "", ""}
  id, err := dao.Create(v1)
  if err != nil {
    t.Fatal(err)
  }
  t.Log("Video id: ", id)
  //videos := sess.DB("test").C("Video")
  //reloaded := bson.M{}
  reloaded := new(Video)
  //oid := bson.ObjectIdHex(id)
  //err = videos.FindId(oid).One(reloaded)
  if err != nil {
    t.Fatal(err)
  }
  t.Logf("reloaded: %v", reloaded)
  err = dao.Get(id, &reloaded)
  if err != nil {
    t.Fatal(err)
  }
  t.Logf("reloaded: %v", reloaded)
  if reloaded.Name != "Novocento" {
    t.Fatal("Expected Novocento")
  }
}
