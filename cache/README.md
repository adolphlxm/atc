# cache

缓存模块，目前支持的引擎有memcache

# 安装

    go get github.com/adolphlxm/atc/cache
   
# 使用步骤
##引入包
    
    import(
        "github.com/adolphlxm/atc/cache"
    )
    
## 接口
```go
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
```
## Memcache使用

```go
package main

import (
    "github.com/adolphlxm/atc/cache"
)

func main() {

    mem, err := cache.NewCache("memcache","127.0.0.1:11211")
    if err != nil {
    
    }
    
    mem.Put("atc","act framework", 10 * time.Second)
    mem.Get("atc")
    mem.Delete("atc")
    mem.FlushAll()
    ...
}

```