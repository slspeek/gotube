package main

import (
  "labix.org/v2/mgo"
  "testing"
)

func dao(t *testing.T) *MongoDao {

  sess, err := mgo.Dial("localhost")
  if err != nil {
    t.Fatal(err)
  }
  return NewMongoDao(sess, "test", "Video")
}

func TestMongoDao(t *testing.T) {
  dao := dao(t)
  v1 := Video{"", "steven", "Novecento", "", ""}
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
  if reloaded.Name != "Novecento" {
    t.Fatal("Expected Novecento")
  }
  dao.Delete(id)
}

func TestMongoDaoId(t *testing.T) {
  dao := dao(t)
  v1 := Video{"", "steven", "Novecento", "", ""}
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
  err = dao.Get(id, &reloaded)
  if err != nil {
    t.Fatal(err)
  }
  t.Logf("reloaded: %#v", reloaded)
  if reloaded.Name != "Novecento" {
    t.Fatal("Expected Novecento")
  }
  if reloaded.Id.Hex() != id {
    t.Fatal("Expected ", id, " b ut was ", reloaded.Id)
  }
  dao.Delete(id)
}

func TestUpdate(t *testing.T) {
  dao := dao(t)
  v1 := Video{"", "steven", "Novecento", "", ""}
  id, err := dao.Create(v1)
  if err != nil {
    t.Fatal(err)
  }
  t.Log("Video id: ", id)
  v1.Name = "Novecento II"
  err = dao.Update(id, v1)
  if err != nil {
    t.Fatal(err)
  }

  reloaded := new(Video)
  err = dao.Get(id, &reloaded)
  if err != nil {
    t.Fatal(err)
  }
  t.Logf("reloaded: %#v", reloaded)
  if reloaded.Name != "Novecento II" {
    t.Fatal("Expected Novecento")
  }
  if reloaded.Id.Hex() != id {
    t.Fatal("Expected ", id, " b ut was ", reloaded.Id)
  }
  dao.Delete(id)
}
