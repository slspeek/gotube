package mongo

import (
	"github.com/slspeek/gotube/common"
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

func createVideo(t *testing.T, v Video) (vout Video, id string) {
	dao := dao(t)
	id, err := dao.Create(v)
	if err != nil {
		t.Fatal(err)
	}
  vout = v
	return
}

func createNovencento(t *testing.T) (v Video, id string) {
	v = Video{"", "steven", "Novecento", "", ""}
	return createVideo(t, v)
}

func TestDao(t *testing.T) {
	dao := dao(t)
	_, id := createNovencento(t)
	reloaded := new(Video)
	err := dao.Get(id, &reloaded)
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
	_, id := createNovencento(t)
	reloaded := new(Video)
	err := dao.Get(id, &reloaded)
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
	v1, id := createNovencento(t)
	v1.Name = "Novecento II"
	err := dao.Update(id, v1)
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
	_, id := createVideo(t, Video{"", "steven", "Novecento II", "", ""})
	createVideo(t,Video{"", "mike", "Novecento III", "", ""})
	reloaded := make([]Video, 1)
  err := dao.Find(bson.M{"owner": "steven"}, &reloaded, []string{})
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
func TestFindOrder(t *testing.T) {
	dao := dao(t)
	dao.DeleteAll()
	createVideo(t, Video{"", "steven", "Novecento II", "", ""})
  _, id := createVideo(t,Video{"", "mike", "Novecento III", "", ""})
	reloaded := make([]Video, 2)
  err := dao.Find(bson.M{}, &reloaded, []string{"owner"})
	if err != nil {
		t.Fatal(err)
	}
	if len(reloaded) != 2 {
		t.Fatal("Expected 2 got ", len(reloaded))
	}
	if reloaded[0].Name != "Novecento III" {
		t.Fatal("Expected Novecento")
	}
	if reloaded[0].Id.Hex() != id {
		t.Fatal("Expected ", id, " but was ", reloaded[0].Id.Hex())
	}
	dao.DeleteAll()
}

func TestVideoDaoPatch(t *testing.T) {
	_, id := createNovencento(t)
	dao := dao(t)
	vdao := VideoDao{dao}
	vInput := common.CVideo{Id: id, Name: "NV", Thumbs: []string{"foo", "bar"}}
	err := vdao.Patch(id, vInput)
	if err != nil {
		t.Fatal(err)
	}
	readBack := new(Video)
	err = vdao.Get(id, readBack)
	if err != nil {
		t.Fatal(err)
	}
}
