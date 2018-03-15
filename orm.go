package atc

import (
	"time"
	"strconv"
	"github.com/adolphlxm/atc/orm"
	"github.com/adolphlxm/atc/orm/util"
	"github.com/adolphlxm/atc/logs"
	"github.com/go-xorm/xorm"
	_ "github.com/adolphlxm/atc/orm/xorm"
)

var dbs orm.Orm

func RunOrms() {
	dbs, _ = orm.NewOrm("xorm")
	for _, aliasname := range Aconfig.OrmAliasNames {
		addrs := AppConfig.Strings("orm." + aliasname)
		logs.Tracef("orm:[%s] starting...", aliasname)
		err := dbs.Open(aliasname, addrs)
		if err != nil {
			panic(err)
		}

		// Check orm connection
		dns, _ := util.ExtractURL(addrs[0])
		if pingtime, err := strconv.Atoi(dns.Options["pingtime"]); err != nil {
			go timerTask(aliasname, int64(pingtime), dbs)
		}

		dbs.SetLevel(aliasname, Aconfig.OrmLogLevel)

		logs.Tracef("orm:[%s] Running.", aliasname)

	}

	return
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
					logs.Tracef("orm:[%s] reconnection Running.", aliasname)
				}
				t.Reset(timeDuration * time.Second)
			}
		}
	}
}

func DB(aliasname string) *xorm.EngineGroup {
	return dbs.Use(aliasname)
}
