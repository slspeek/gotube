package rest

import (
	"labix.org/v2/mgo/bson"
	"testing"
)

func TestVideoPostprocess(t *testing.T) {
	vIn := NOVECENTO
	vIn.BlobId = "foo"
  vIn.Thumbs = []string{"foo", "bar"}
	vIn.Id = bson.NewObjectId()

  out := videoPostProcess(vIn, false)
  expected := "/content/videos/"+vIn.Id.Hex()
	if out.Stream != expected {
    t.Fatal("Private stream: ", out.Stream)
		t.Fail()
	}
	if out.Thumbs[0] != "/content/videos/"+vIn.Id.Hex()+"/thumbs/0" {
    t.Fatal("Private thumb: ", out.Thumbs[0])
		t.Fail()
	}
  if out.Download !=expected + "/download" {
    t.Fatal("Private Download: ", out.Download)
  }

   out = videoPostProcess(vIn, true)
	if out.Stream != "/public" + expected {
		t.Fatal("Stream error", out.Stream)
	}
	if out.Thumbs[0] != "/public/content/videos/"+vIn.Id.Hex()+"/thumbs/0" {
		t.Fail()
	}
  if out.Download != "/public" + expected + "/download" {
    t.Fatal("Private Download: ", out.Download)
  }
}
