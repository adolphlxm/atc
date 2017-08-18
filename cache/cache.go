package cache

import (
	"fmt"
	"time"
)

type Cache interface {
	Get(key string) interface{}
	Put(key string, val interface{}, timeout time.Duration) error
	Delete(key string) error
	FlushAll() error
	New(config string) error
}

type CacheFunc func() Cache

var adapters = make(map[string]CacheFunc)

func Register(name string, adapter CacheFunc) {
	if adapter == nil {

		panic("ATC cache: Register handler is nil")
	}
	if _, found := adapters[name]; found {
		panic("ATC cache: Register failed for handler " + name)
	}
	adapters[name] = adapter
}

func NewCache(adapterName, config string) (Cache, error) {
	if handler, ok := adapters[adapterName]; ok {
		return handler(), nil
	} else {
		return nil, fmt.Errorf("ATC cache: unknown adapter name %s failed.", adapterName)
	}
}
