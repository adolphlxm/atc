package redis

import (
	"github.com/garyburd/redigo/redis"
	"github.com/adolphlxm/atc/cache"
	"time"
	"strconv"
	"encoding/json"
	"errors"
)

type Cache struct {
	p *redis.Pool
	conninfo map[string]string
}

func NewRedisCache() cache.Cache {
	return &Cache{}
}

// StartAndGC start memcache adapter.
// if connecting error, return.
func (c *Cache) New(config string) error {
	err := json.Unmarshal([]byte(config), &c.conninfo)
	if err != nil {
		return err
	}

	if c.p == nil {
		if err := c.connect(); err != nil {
			return err
		}
	}
	return nil
}

func (c *Cache) connect() error {
	maxidle,_ := strconv.Atoi(c.conninfo["maxidle"])
	maxactive, _ := strconv.Atoi(c.conninfo["maxactive"])
	idletimeout, _ := strconv.Atoi(c.conninfo["idletimeout"])

	if c.conninfo["addr"] == "" {
		return errors.New("Redis addr is empty.")
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