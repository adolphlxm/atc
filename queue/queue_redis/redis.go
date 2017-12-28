package queue_redis

import (
	"time"

	"github.com/adolphlxm/atc/queue"
	"github.com/garyburd/redigo/redis"
)

type RedisQueueDriver struct {
}

// dsn format is `redis://:pass@host1:port/db?options`
//
// options can be:
// maxIdle - default is 1
// maxActive - default is 1
//
// all timeout option's unit is ms
// idleTimeout    - default is 0ms
// connectTimeout - default is 3000ms
// readTimeout    - default is 0ms
// writeTimeout   - default is 0ms

func (d *RedisQueueDriver) Open(addr string) (queue.Conn, error) {
	info, err := parseURL(addr)
	if err != nil {
		return nil, err
	}

	c := &redisQueueConn{
		cs: &redis.Pool{
			Dial: func() (redis.Conn, error) {
				conn, err := redis.DialURL(
					info.Url,
					redis.DialConnectTimeout(info.ConnectTimeout),
					redis.DialReadTimeout(info.ReadTimeout),
					redis.DialWriteTimeout(info.WriteTimeout),
				)
				if err != nil {
					return nil, err
				}

				return conn, nil
			},

			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				if time.Since(t) < 3*time.Second {
					return nil
				}
				_, err := c.Do("PING")
				return err
			},
			MaxIdle:     info.MaxIdle,
			MaxActive:   info.MaxActive,
			IdleTimeout: info.IdleTimeout,
			Wait:        true,
		},
	}
	return c, nil
}

func init() {
	queue.Register("redis", &RedisQueueDriver{})
}
