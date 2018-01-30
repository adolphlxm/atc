package mgo

import (
	"time"

	"gopkg.in/mgo.v2"
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