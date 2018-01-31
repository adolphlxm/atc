package mgo

import (
	"time"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"strings"
)

type MgoDB struct {
	session *mgo.Session
}

func NewMgoDB(addrs string) (*MgoDB, error) {
	m := &MgoDB{}
	err := m.open(addrs)
	return m, err
}

/************************************/
/**********   Collection  ***********/
/************************************/
func (this *MgoDB) Insert(name, col string, docs ...interface{}) error {
	session := this.session.Clone()
	defer session.Close()
	return session.DB(strings.TrimSpace(name)).C(strings.TrimSpace(col)).Insert(docs...)
}

func (this *MgoDB) Update(name, col string, selector interface{}, update interface{}) error {
	session := this.session.Clone()
	defer session.Close()
	return session.DB(strings.TrimSpace(name)).C(strings.TrimSpace(col)).Update(selector, update)
}

func (this *MgoDB) UpdateId(name, col string, id bson.M, update interface{}) error {
	session := this.session.Clone()
	defer session.Close()
	return session.DB(strings.TrimSpace(name)).C(strings.TrimSpace(col)).UpdateId(id, update)
}

func (this *MgoDB) UpdateAll(name, col string, selector interface{}, update interface{}) (*mgo.ChangeInfo, error) {
	session := this.session.Clone()
	defer session.Close()
	return session.DB(strings.TrimSpace(name)).C(strings.TrimSpace(col)).UpdateAll(selector, update)
}

// Not closed session
func (this *MgoDB) Find(name, col string, query interface{}) *MgoQuery {
	session := this.session.Clone()
	return &MgoQuery{session:session, query:session.DB(strings.TrimSpace(name)).C(strings.TrimSpace(col)).Find(query)}
}

func (this *MgoDB) Remove(name, col string, selector bson.M) error {
	session := this.session.Clone()
	defer session.Close()
	return session.DB(strings.TrimSpace(name)).C(strings.TrimSpace(col)).Remove(selector)
}

func (this *MgoDB) RemoveId(name, col string, id bson.M) error {
	session := this.session.Clone()
	defer session.Close()
	return session.DB(strings.TrimSpace(name)).C(strings.TrimSpace(col)).RemoveId(id)
}

func (this *MgoDB) RemoveAll(name, col string, selector bson.M) (*mgo.ChangeInfo, error) {
	session := this.session.Clone()
	defer session.Close()
	return session.DB(strings.TrimSpace(name)).C(strings.TrimSpace(col)).RemoveAll(selector)
}

/************************************/
/**********   Collection  ***********/
/************************************/

func (this *MgoDB) CollectionNames(name string) ([]string, error) {
	session := this.session.Clone()
	defer session.Close()
	return session.DB(strings.TrimSpace(name)).CollectionNames()
}


/************************************/
/**********  Find Query   ***********/
/************************************/
type MgoQuery struct {
	session *mgo.Session
	query *mgo.Query
}

func (q *MgoQuery) Close() {
	q.session.Close()
}

func (q *MgoQuery) One(result interface{}) error {
	defer q.Close()
	return q.query.One(result)
}

func (q *MgoQuery) All(result interface{}) error {
	defer q.Close()
	return q.query.All(result)
}

func (q *MgoQuery) Apply(change mgo.Change,result interface{}) (*mgo.ChangeInfo, error) {
	defer q.Close()
	return q.query.Apply(change, result)
}

func (q *MgoQuery) FindCount(name, col string, query bson.M) (int, error) {
	defer q.Close()
	return q.query.Count()
}

func (q *MgoQuery) Distinct(key string, result interface{}) error {
	defer q.Close()
	return q.query.Distinct(key, result)
}

func (q *MgoQuery) Explain(result interface{}) error {
	defer q.Close()
	return q.query.Explain(result)
}

func (q *MgoQuery) Query() *mgo.Query {
	return q.query
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
