package orm

import (
	"encoding/json"
	"time"
)

/************************************/
/********** Request input ***********/
/************************************/

var DB Orm

type Orms struct {
	debug bool
	maxidleconns int
	maxopenconns int
	pingtime int
	ormAliasNames map[string]*OrmConfig
}

type OrmConfig struct {
	Driver   string `json:"driver"`
	Host     string `json:"host"`
	User     string `json:"user"`
	Password string `json:"password"`
	DbName   string `json:"dbname"`
	LogLevel string `json:"loglevel"`
	Maxidleconns int
	Maxopenconns int
	Pingtime int
}

func NewOrms(maxidleconns, maxopenconns, pingtime int, ormAliasNames map[string]*OrmConfig, debug bool) *Orms{
	return &Orms{
		debug:debug,
		maxidleconns:maxidleconns,
		maxopenconns:maxopenconns,
		pingtime:pingtime,
		ormAliasNames:ormAliasNames,
	}
}

func (this *Orms) Run() (db Orm){
	db, _ = NewOrm("xorm")

	for aliasname, cfg := range this.ormAliasNames {
		if cfg.Maxidleconns > 0 {
			this.maxidleconns = cfg.Maxidleconns
		}
		if cfg.Maxopenconns > 0 {
			this.maxopenconns = cfg.Maxopenconns
		}
		if cfg.Pingtime > 0 {
			this.pingtime = cfg.Pingtime
		}
		cfJson, err := json.Marshal(cfg)
		if err != nil {
			panic(err)
		}
		if err := db.Open(aliasname, string(cfJson)); err != nil {
			panic(err)
		}
		db.Debug(aliasname, this.debug)
		db.SetMaxIdleConns(aliasname, this.maxidleconns)
		db.SetMaxOpenConns(aliasname, this.maxopenconns)

		// Check orm connection
		go timerTask(aliasname, int64(this.pingtime))
	}
	return
}

func timerTask(aliasname string, timeout int64) {
	if timeout > 0 {
		timeDuration := time.Duration(timeout)
		t := time.NewTimer(timeDuration * time.Second)
		for {
			select {
			case <-t.C:
				//if err := Orm.Ping(aliasname); err != nil {
				//	Orm.Clone(aliasname)
				//	Logger.Trace("ATC orm: reconnection successful to %s", aliasname)
				//}
				t.Reset(timeDuration * time.Second)
			}
		}
	}
}
