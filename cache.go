package atc

import (
	"net/url"

	"github.com/adolphlxm/atc/cache"
	_ "github.com/adolphlxm/atc/cache/memcache"
	_ "github.com/adolphlxm/atc/cache/redis"
	"github.com/adolphlxm/atc/logs"
)

var aCache map[string]cache.Cache

func RunCaches() {
	aCache = make(map[string]cache.Cache, 0)
	aliasnames := AppConfig.Strings("cache.aliasnames")
	for _, aliasname := range aliasnames {
		keyPerfix := "cache." + aliasname + "."
		addr := AppConfig.String(keyPerfix + "addrs")
		logs.Tracef("cache:[%s] starting...", aliasname)
		var config string

		addrUrl, err := url.Parse(addr)
		if err != nil {
			logs.Errorf("cache:[%s] parse addrs err:%s", aliasname, err.Error())
			panic(err)
		}
		drivename := addrUrl.Scheme

		switch drivename {
		case "memcache":
			config = addrUrl.Host
		case "redis":
			redisAddr := addrUrl.Host
			queryValue := addrUrl.Query()
			password := ""
			if userInfo := addrUrl.User; userInfo != nil {
				password, _ = addrUrl.User.Password()
			}
			config = `{"addr":"` + redisAddr + `","maxidle":"` + queryValue.Get("maxIdle") + `","maxactive":"` + queryValue.Get("maxActive") + `","idletimeout":"` + queryValue.Get("idleTimeout") + `","password":"` + password + `"}`
		default:
			continue
		}

		aCache[aliasname], err = cache.NewCache(drivename, config)
		if err != nil {
			logs.Errorf("cache:[%s] start fail err:%s", aliasname, err.Error())
			panic(err)
		}

		logs.Tracef("cache:[%s] Running on %s.", aliasname, addrUrl.Host)
	}
}

func GetCache(aliasname string) cache.Cache {
	return aCache[aliasname]
}
