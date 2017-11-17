package redis

import (
	"github.com/adolphlxm/atc/cache"
	"github.com/garyburd/redigo/redis"
	"testing"
	"time"
)

func TestNewRedisCache(t *testing.T) {
	r, err := cache.NewCache("redis", `{"addr":"127.0.0.1:6379","maxidle":"2","maxactive":"2","idletimeout":"5"}`)
	if err != nil {
		t.Errorf("Redis init err:%v", err.Error())
	}
	key := "atc"

	// test put
	timeoutDuration := 10 * time.Second
	if err = r.Put(key, "atc framework", timeoutDuration); err != nil {
		t.Errorf("Redis put err:%v", err.Error())
	}

	// test get
	v, err := redis.String(r.Get(key))
	if err != nil {
		t.Errorf("Redis Get err:%v", err.Error())
	}

	if v != "atc framework" {
		t.Errorf("Redis Get value error.")
	}

	// test delete
	err = r.Delete(key)
	if err != nil {
		t.Errorf("Redis Delete err:%v", err.Error())
	}
	b, err := redis.Bool(r.Do("EXISTS", key))
	if err != nil {
		t.Errorf("Redis Do err:%v", err.Error())
	}
	if b {
		t.Errorf("Redis %s key is EXISTS.", key)
	}

	// test increment
	err = r.Increment(key + "_rement")
	if err != nil {
		t.Errorf("Redis Increment err:%v", err.Error())
	}

	// test decrement
	err = r.Decrement(key + "_rement")
	if err != nil {
		t.Errorf("Redis Decrement err:%v", err.Error())
	}

}
