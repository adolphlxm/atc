package mgo

import (
	"time"
	"strings"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MgoDB struct {
	session *mgo.Session
}

func NewMgoDB(addrs string) (*MgoDB, error) {
	m := &MgoDB{}
	err := m.open(addrs)
	return m, err
}

func (this *MgoDB) Clone() *MgoDB {
	session := this.session.Clone()
	return &MgoDB{session: session}
}

func (this *MgoDB) Close() {
	this.session.Close()
}

func (this *MgoDB) Insert(name, col string, docs ...interface{}) error {
	session := this.session.Clone()
	defer session.Close()
	return session.DB(name).C(col).Insert(docs...)
}

func (this *MgoDB) Update(name, col string, filter bson.M, update interface{}) error {
	session := this.session.Clone()
	defer session.Close()
	return session.DB(name).C(col).Update(filter, update)
}

func (this *MgoDB) UpdateId(name, col string, id, update interface{}) error {
	session := this.session.Clone()
	defer session.Close()
	return session.DB(name).C(col).UpdateId(id, update)
}

func (this *MgoDB) UpdateAll(name, col string, filter bson.M, update interface{}) (*mgo.ChangeInfo,error) {
	session := this.session.Clone()
	defer session.Close()
	return session.DB(name).C(col).UpdateAll(filter, update)
}

func (this *MgoDB) Find(name, col string, query interface{}) *mgo.Query{
	session := this.session.Clone()
	defer session.Close()
	return session.DB(name).C(col).Find(query)
}

func (this *MgoDB) FindOne(name, col string, filter bson.M, docs interface{}) error {
	session := this.session.Clone()
	defer session.Close()
	return session.DB(strings.TrimSpace(name)).C(strings.TrimSpace(col)).Find(filter).One(docs)
}

func (this *MgoDB) FindAll(name, col string, filter bson.M, docs interface{}) error {
	session := this.session.Clone()
	defer session.Close()
	return session.DB(strings.TrimSpace(name)).C(strings.TrimSpace(col)).Find(filter).All(docs)
}

func (this *MgoDB) RemoveId(name, col string, id interface{}) error {
	session := this.session.Clone()
	defer session.Close()
	return session.DB(strings.TrimSpace(name)).C(strings.TrimSpace(col)).RemoveId(id)
}

func (this *MgoDB) RemoveAll(name, col string, filter bson.M) error {
	session := this.session.Clone()
	defer session.Close()
	_, err := session.DB(strings.TrimSpace(name)).C(strings.TrimSpace(col)).RemoveAll(filter)
	return err
}

func (this *MgoDB) open(addrs string) error {
	dialInfo, err := mgo.ParseURL(addrs)
	if err != nil {
		return nil
	}
	this.session, err = mgo.DialWithInfo(dialInfo)
	if err != nil {
		return nil
	}

	this.session.SetSyncTimeout(dialInfo.Timeout * time.Second)
	this.session.SetSocketTimeout(dialInfo.Timeout * time.Second)

	return nil
}