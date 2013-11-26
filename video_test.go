package main

import (
  "github.com/slspeek/go-restful"
  "io"
  "labix.org/v2/mgo"
  "net/http"
  "net/http/httptest"
  "os"
  "strings"
  "testing"
)

func TestVideoResource(t *testing.T) {
  t.Log("Video test")
  //t.Skip()
  sess, err := mgo.Dial("localhost")
  if err != nil {
    t.Fatal(err)
  }
  vr := NewVideoResource(sess, "test", "Video")
  container := restful.NewContainer()
  vr.Register(container)
  ts := httptest.NewServer(container)
  defer ts.Close()

  resp, err := http.Post(ts.URL + "/videos", "application/json", strings.NewReader(`{"Owner":"","Name":"As it is in heaven"}`) )
  io.Copy(os.Stderr, resp.Body)


 if resp.StatusCode != http.StatusCreated {
   t.Fatal("StatusCode: ", resp.StatusCode)
  }
}

func TestGet(t *testing.T) {
  sess, err := mgo.Dial("localhost")
  if err != nil {
    t.Fatal(err)
  }
  dao := NewMongoDao(sess, "test", "Video")
  v1 := Video{"", "", "Novocento", "", ""}
  id, err := dao.Create(v1)
  if err != nil {
    t.Fatal(err)
  }
  t.Log("Video id: ", id)

  vr := NewVideoResource(sess, "test", "Video")
  container := restful.NewContainer()
  vr.Register(container)
  ts := httptest.NewServer(container)
  defer ts.Close()
  url := ts.URL + "/videos/" + id 
  t.Log("URL: ", url)
  resp, err := http.Get(url)
  io.Copy(os.Stderr, resp.Body)


 if resp.StatusCode != http.StatusOK {
   t.Fatal("StatusCode: ", resp.StatusCode)
  }
}
