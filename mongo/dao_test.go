package mongo

import (
	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"
	"testing"
)

type Video struct {
	Id     bson.ObjectId "_id,omitempty"
	Owner  string
	Name   string
	Desc   string
	BlobId string
}

func dao(t *testing.T) *Dao {
	sess, err := mgo.Dial("localhost")
	if err != nil {
		t.Fatal(err)
	}
	return NewDao(sess, "test", "Video")
}

func TestDao(t *testing.T) {
	dao := dao(t)
	v1 := Video{"", "steven", "Novecento", "", ""}
	id, err := dao.Create(v1)
	if err != nil {
		t.Fatal(err)
	}
	reloaded := new(Video)
	err = dao.Get(id, &reloaded)
	if err != nil {
		t.Fatal(err)
	}
	if reloaded.Name != "Novecento" {
		t.Fatal("Expected Novecento")
	}
	dao.Delete(id)
}

func TestDaoId(t *testing.T) {
	dao := dao(t)
	v1 := Video{"", "steven", "Novecento", "", ""}
	id, err := dao.Create(v1)
	if err != nil {
		t.Fatal(err)
	}
	reloaded := new(Video)
	err = dao.Get(id, &reloaded)
	if err != nil {
		t.Fatal(err)
	}
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
	if reloaded.Name != "Novecento II" {
		t.Fatal("Expected Novecento")
	}
	if reloaded.Id.Hex() != id {
		t.Fatal("Expected ", id, " but was ", reloaded.Id)
	}
	dao.Delete(id)
}

func TestGetAll(t *testing.T) {
	dao := dao(t)
	dao.DeleteAll()
	v1 := Video{"", "steven", "Novecento II", "", ""}
	id, err := dao.Create(v1)
	if err != nil {
		t.Fatal(err)
	}
	reloaded := make([]Video, 1)
	err = dao.GetAll(&reloaded)
	if err != nil {
		t.Fatal(err)
	}

	if reloaded[0].Name != "Novecento II" {
		t.Fatal("Expected Novecento")
	}
	if reloaded[0].Id.Hex() != id {
		t.Fatal("Expected ", id, " but was ", reloaded[0].Id)
	}
	dao.DeleteAll()
}

func TestFind(t *testing.T) {
	dao := dao(t)
	dao.DeleteAll()
	v1 := Video{"", "steven", "Novecento II", "", ""}
	id, err := dao.Create(v1)
	if err != nil {
		t.Fatal(err)
	}
	v2 := Video{"", "mike", "Novecento III", "", ""}
	_, err = dao.Create(v2)
	if err != nil {
		t.Fatal(err)
	}
	reloaded := make([]Video, 1)
	err = dao.Find(bson.M{"owner": "steven"}, &reloaded)
	if err != nil {
		t.Fatal(err)
	}

	if len(reloaded) != 1 {
		t.Fatal("Expected 1 got ", len(reloaded))
	}
	if reloaded[0].Name != "Novecento II" {
		t.Fatal("Expected Novecento")
	}
	if reloaded[0].Id.Hex() != id {
		t.Fatal("Expected ", id, " but was ", reloaded[0].Id)
	}
	dao.DeleteAll()
}
