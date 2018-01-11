package atc

import (
	"encoding/json"
	"errors"
	"time"

	"github.com/adolphlxm/atc/logs"
	"github.com/adolphlxm/atc/orm"
	_ "github.com/adolphlxm/atc/orm/xorm"
	"github.com/go-xorm/xorm"
)

var dbs orm.Orm

type OrmConfig struct {
	Driver   string `json:"driver"`
	Host     string `json:"host"`
	User     string `json:"user"`
	Password string `json:"password"`
	DbName   string `json:"dbname"`
	LogLevel string `json:"loglevel"`
}

func RunOrms() {
	var (
		maxidleconns int
		maxopenconns int
		pingtime     int
	)
	maxidleconns = AppConfig.DefaultInt("orm.maxidleconns", 0)
	maxopenconns = AppConfig.DefaultInt("orm.maxopenconns", 0)
	pingtime = AppConfig.DefaultInt("orm.pingtime", 0)

	dbs, _ = orm.NewOrm("xorm")
	for _, aliasname := range Aconfig.OrmAliasNames {
		keyPerfix := "orm." + aliasname
		cfg, err := newEngineConfig(keyPerfix)
		if err != nil {
			panic(err)
		}

		if conns1 := AppConfig.DefaultInt(keyPerfix+".c.maxidleconns", 0); conns1 > 0 {
			maxidleconns = conns1
		}
		if conns2 := AppConfig.DefaultInt(keyPerfix+".c.maxopenconns", 0); conns2 > 0 {
			maxopenconns = conns2
		}
		if apingtime := AppConfig.DefaultInt(keyPerfix+".c.pingtime", 0); apingtime > 0 {
			pingtime = apingtime
		}

		if err := dbs.Open(aliasname, cfg); err != nil {
			panic(err)
		}
		dbs.Debug(aliasname, Aconfig.Debug)
		dbs.SetMaxIdleConns(aliasname, maxidleconns)
		dbs.SetMaxOpenConns(aliasname, maxopenconns)

		// Check orm connection
		go timerTask(aliasname, int64(pingtime), dbs)
	}

	return
}

func newEngineConfig(keyPerfix string) (string, error) {

	cf := &OrmConfig{
		Driver:   AppConfig.DefaultString(keyPerfix+".driver", ""),
		Host:     AppConfig.DefaultString(keyPerfix+".host", ""),
		User:     AppConfig.DefaultString(keyPerfix+".user", ""),
		Password: AppConfig.DefaultString(keyPerfix+".password", ""),
		DbName:   AppConfig.DefaultString(keyPerfix+".dbname", ""),
		LogLevel: Aconfig.OrmLogLevel,
	}
	if cf.Host == "" || cf.User == "" {
		return "", errors.New("Host is empty.")
	}
	cfJson, err := json.Marshal(cf)
	return string(cfJson), err
}

func timerTask(aliasname string, timeout int64, db orm.Orm) {
	if timeout > 0 {
		timeDuration := time.Duration(timeout)
		t := time.NewTimer(timeDuration * time.Second)
		for {
			select {
			case <-t.C:
				if err := db.Ping(aliasname); err != nil {
					db.Clone(aliasname)
					logs.Tracef("ATC orm: reconnection successful to %s", aliasname)
				}
				t.Reset(timeDuration * time.Second)
			}
		}
	}
}


func DB(aliasname string) *xorm.Engine {
	return dbs.Use(aliasname)
}