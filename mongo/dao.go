package mongo

import "labix.org/v2/mgo"
import "labix.org/v2/mgo/bson"

type Dao struct {
	session *mgo.Session
	db      string
	kind    string
}

func NewDao(s *mgo.Session, db string, kind string) *Dao {
	return &Dao{s, db, kind}
}

func (self *Dao) Create(entity interface{}) (id string, err error) {
	oid := bson.NewObjectId()
	info, err := self.collection().UpsertId(oid, entity)
	return info.UpsertedId.(bson.ObjectId).Hex(), err
}

func (self *Dao) Get(id string, result interface{}) (err error) {
	oid := bson.ObjectIdHex(id)
	err = self.collection().FindId(oid).One(result)
	return
}

func (self *Dao) GetAll(result interface{}) (err error) {
	err = self.collection().Find(nil).All(result)
	return
}

func (self *Dao) Find(constraint bson.M, result interface{}) (err error) {
	err = self.collection().Find(constraint).All(result)
	return
}

func (self *Dao) Update(id string, value interface{}) (err error) {
	oid := bson.ObjectIdHex(id)
	err = self.collection().UpdateId(oid, value)
	return
}

func (self *Dao) Delete(id string) (err error) {
	oid := bson.ObjectIdHex(id)
	err = self.collection().RemoveId(oid)
	return
}

func (self *Dao) DeleteAll() (err error) {
	_, err = self.collection().RemoveAll(nil)
	return
}

func (self *Dao) collection() *mgo.Collection {
	return self.session.DB(self.db).C(self.kind)
}
