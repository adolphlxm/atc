package atc

import (
	"net/url"

	"github.com/adolphlxm/atc/cache"
	_ "github.com/adolphlxm/atc/cache/memcache"
	_ "github.com/adolphlxm/atc/cache/redis"
	"github.com/adolphlxm/atc/logs"
	"fmt"
)

var aCache map[string]cache.Cache

func RunCaches() {
	aCache = make(map[string]cache.Cache, 0)
	aliasnames := AppConfig.Strings("cache.aliasnames")
	for _, aliasname := range aliasnames {
		keyPerfix := "cache." + aliasname + "."
		addr := AppConfig.String(keyPerfix + "addrs")

		var config string

		addrUrl, err := url.Parse(addr)
		if err != nil {
			logs.Errorf("cache: aliasname:%s,Parse addrs err:%s", aliasname, err.Error())
			continue
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
			logs.Errorf("cache: aliasname:%s,NewCache err:%s", aliasname, err.Error())
			continue
		}
		logs.Tracef("cache: aliasname:%s,initialization is successful.", aliasname)
	}
}

func GetCache(aliasname string) cache.Cache {
	return aCache[aliasname]
}
