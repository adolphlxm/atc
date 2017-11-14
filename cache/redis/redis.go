package redis

import (
	"github.com/garyburd/redigo/redis"
	"github.com/adolphlxm/atc/cache"
	"time"
	"strconv"
	"encoding/json"
	"errors"
)

const DefaultKey = "atcCacheRedis"

type Cache struct {
	p *redis.Pool
	conninfo map[string]string
}

func NewRedisCache() cache.Cache {
	return &Cache{}
}

// actually do the redis cmds
func (c *Cache) Do(commandName string, args ...interface{}) (reply interface{}, err error) {
	p := c.p.Get()
	defer p.Close()

	return p.Do(commandName, args...)
}

func (c *Cache) Get(key string) (reply interface{}, err error) {
	reply, err = c.Do("GET", key)
	return
}

func (c *Cache) Put(key string, val interface{}, timeout time.Duration) error {
	var err error
	if _, err = c.Do("SETEX", key, int64(timeout/time.Second), val); err != nil {
		return err
	}

	 _, err = c.Do("HSET", c.conninfo["key"], key, true)
	return err
}

func (c *Cache) Delete(key string) error {
	var err error
	if _, err = c.Do("DEL", key); err != nil {
		return err
	}
	_, err = c.Do("HDEL", c.conninfo["key"], key, true)
	return err
}

// Incr increase counter in redis.
func (c *Cache) Increment(key string) error {
	_, err := redis.Bool(c.Do("INCRBY", key, 1))
	return err
}

// Decr decrease counter in redis.
func (c *Cache) Decrement(key string) error {
	_, err := redis.Bool(c.Do("INCRBY", key, -1))
	return err
}

// FlushAll clear all cached in memcache.
func (c *Cache) ClearAll() error {
	cacheKeys, err := redis.Strings(c.Do("HKEYS", c.conninfo["key"]))
	if err != nil {
		return err
	}
	for _, key := range cacheKeys{
		if _, err = c.Do("DEL",key); err != nil {
			return err
		}
	}

	_, err = c.Do("DEL", c.conninfo["key"])
	return err
}

// StartAndGC start memcache adapter.
// if connecting error, return.
func (c *Cache) New(config string) error {
	err := json.Unmarshal([]byte(config), &c.conninfo)
	if err != nil {
		return err
	}

	if c.p == nil {
		if err := c.connectInit(); err != nil {
			return err
		}
	}
	return nil
}


func (c *Cache) connectInit() error {
	maxidle,_ := strconv.Atoi(c.conninfo["maxidle"])
	maxactive, _ := strconv.Atoi(c.conninfo["maxactive"])
	idletimeout, _ := strconv.Atoi(c.conninfo["idletimeout"])

	if c.conninfo["addr"] == "" {
		return errors.New("Redis addr is empty.")
	}
	if c.conninfo["key"] == "" {
		c.conninfo["key"] = DefaultKey
	}

	c.p = &redis.Pool{
		MaxIdle:maxidle,
		MaxActive:   maxactive,
		IdleTimeout: time.Duration(idletimeout) * time.Second,
		Dial: func() (redis.Conn, error) {
			r, err := redis.Dial("tcp", c.conninfo["addr"])
			if err != nil {
				return nil, err
			}
			if c.conninfo["password"] != "" {
				if _, err := r.Do("AUTH", c.conninfo["password"]); err != nil {
					r.Close()
					return nil, err
				}
			}
			return r, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < 3*time.Second {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}

	return nil
}

func init() {
	cache.Register("redis", NewRedisCache)
}