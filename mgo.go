package atc

import (
	"github.com/adolphlxm/atc/mgo"
	"github.com/adolphlxm/atc/logs"
)

var mgoDBs map[string]*mgo.MgoDB

func RunMgoDBs() {
	mgoDBs = make(map[string]*mgo.MgoDB, 0)
	aliasnames := AppConfig.Strings("mgo.aliasnames")
	for _, aliasname := range aliasnames {
		keyPerfix := "mgo." + aliasname + "."
		addrs := AppConfig.String(keyPerfix + "addrs")
		logs.Tracef("mgo:[%s] starting...", aliasname)
		db, err := mgo.NewMgoDB(addrs)
		if err != nil {
			logs.Errorf("mgo:[%s] start fail err:%s", aliasname, err.Error())
			continue
		}

		mgoDBs[aliasname] = db
		logs.Tracef("mgo:[%s] Running on %s.", aliasname, addrs)
	}
}

func GetMgoDB(aliasname string) *mgo.MgoDB {
	return mgoDBs[aliasname]
}