package redis

import (
	"github.com/garyburd/redigo/redis"
	"github.com/adolphlxm/atc/cache"
)

type Cache struct {
	p *redis.Pool
	conninfo string
}

func NewRedisCache() cache.Cache {
	return &Cache{}
}

func (c *Cache) connect() {

}