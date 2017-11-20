package redis

import (
	"encoding/json"
	"errors"
	"github.com/adolphlxm/atc/cache"
	"github.com/garyburd/redigo/redis"
	"strconv"
	"time"
)

const DefaultKey = "atcCacheRedis"

type Cache struct {
	p        *redis.Pool
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
	for _, key := range cacheKeys {
		if _, err = c.Do("DEL", key); err != nil {
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
	maxidle, _ := strconv.Atoi(c.conninfo["maxidle"])
	maxactive, _ := strconv.Atoi(c.conninfo["maxactive"])
	idletimeout, _ := strconv.Atoi(c.conninfo["idletimeout"])

	if c.conninfo["addr"] == "" {
		return errors.New("Redis addr is empty.")
	}
	if c.conninfo["key"] == "" {
		c.conninfo["key"] = DefaultKey
	}

	c.p = &redis.Pool{
		MaxIdle:     maxidle,
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


/************************************/
/**********  Redis Reply  ***********/
/************************************/
func Strings(reply interface{}, err error) ([]string, error){
	return redis.Strings(reply, err)
}

func String(reply interface{}, err error) (string, error){
	return redis.String(reply, err)
}
// Int is a helper that converts a command reply to an integer. If err is not
// equal to nil, then Int returns 0, err. Otherwise, Int converts the
// reply to an int as follows:
//
//  Reply type    Result
//  integer       int(reply), nil
//  bulk string   parsed reply, nil
//  nil           0, ErrNil
//  other         0, error
func Int(reply interface{}, err error) (int, error){
	return redis.Int(reply, err)
}

func Int64(reply interface{}, err error) (int64, error){
	return redis.Int64(reply, err)
}

func Uint64(reply interface{}, err error) (uint64, error){
	return redis.Uint64(reply, err)
}

func Float64(reply interface{}, err error) (float64, error){
	return redis.Float64(reply, err)
}

func Bytes(reply interface{}, err error) ([]byte, error){
	return redis.Bytes(reply,err)
}

func Bool(reply interface{}, err error) (bool, error) {
	return redis.Bool(reply, err)
}

func Values(reply interface{}, err error) ([]interface{}, error) {
	return redis.Values(reply, err)
}

func Ints(reply interface{}, err error) ([]int, error){
	return redis.Ints(reply, err)
}

func StringMap(reply interface{}, err error) (map[string]string, error){
	return redis.StringMap(reply, err)
}

func IntMap(reply interface{}, err error) (map[string]int, error) {
	return redis.IntMap(reply, err)
}

func Int64Map(reply interface{}, err error) (map[string]int64, error){
	return redis.Int64Map(reply, err)
}

func init() {
	cache.Register("redis", NewRedisCache)
}