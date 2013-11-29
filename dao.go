package main

import "labix.org/v2/mgo"
import "labix.org/v2/mgo/bson"

type MongoDao struct {
	session *mgo.Session
	db      string
	kind    string
}

func NewMongoDao(s *mgo.Session, db string, kind string) *MongoDao {
	return &MongoDao{s, db, kind}
}

func (self *MongoDao) Create(entity interface{}) (id string, err error) {
	oid := bson.NewObjectId()
	info, err := self.collection().UpsertId(oid, entity)
	return info.UpsertedId.(bson.ObjectId).Hex(), err
}

func (self *MongoDao) Get(id string, result interface{}) (err error) {
	oid := bson.ObjectIdHex(id)
	err = self.collection().FindId(oid).One(result)
	return
}

func (self *MongoDao) GetAll(result interface{}) (err error) {
	err = self.collection().Find(nil).All(result)
	return
}

func (self *MongoDao) Find(constraint bson.M, result interface{}) (err error) {
	err = self.collection().Find(constraint).All(result)
	return
}

func (self *MongoDao) Update(id string, value interface{}) (err error) {
	oid := bson.ObjectIdHex(id)
	err = self.collection().UpdateId(oid, value)
	return
}

func (self *MongoDao) Delete(id string) (err error) {
	oid := bson.ObjectIdHex(id)
	err = self.collection().RemoveId(oid)
	return
}

func (self *MongoDao) DeleteAll() (err error) {
	_, err = self.collection().RemoveAll(nil)
	return
}

func (self *MongoDao) collection() *mgo.Collection {
	return self.session.DB(self.db).C(self.kind)
}
