package cache

import (
	"fmt"
	"time"
)

type Cache interface {
	// Only user redis.
	// Do sends a command to the server and returns the received reply.
	Do(commandName string, args ...interface{}) (reply interface{}, err error)
	// Get cached value by key.
	Get(key string) (interface{}, error)
	// Put cached value with key and expire time.
	Put(key string, val interface{}, timeout time.Duration) error
	// Increment cached int value by key, as a counter.
	Increment(key string) error
	// Decrement cached int value by key, as a counter.
	Decrement(key string) error
	// Delete cached value by key.
	Delete(key string) error
	// Clear all cache.
	ClearAll() error

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
		adapter := handler()
		err := adapter.New(config)
		return adapter, err
	} else {
		return nil, fmt.Errorf("ATC cache: unknown adapter name %s failed.", adapterName)
	}
}
