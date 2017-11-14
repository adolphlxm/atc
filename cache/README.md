# cache

缓存模块，目前支持的引擎有Redis、Memcache

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


## Redis使用
config json 配置字段说明

`{"addr":"127.0.0.1:6379","maxidle":"2","maxactive":"2","idletimeout":"5"}`
* addr：连接地址及端口
* maxidle：最大空闲连接数
* maxactive：最大连接数
* idletimeout：空闲连接超时时间

```go
package main

import (
    "github.com/adolphlxm/atc/cache"
)

func main() {

    red, err := cache.NewCache("redis",`{"addr":"127.0.0.1:6379","maxidle":"2","maxactive":"2","idletimeout":"5"}`)
    if err != nil {

    }
    red.Do("RPUSH", key, "data")
    red.Put("atc","act framework", 10 * time.Second)
    red.Get("atc")
    red.Delete("atc")
    red.ClearAll()
    ...
}

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
    mem.ClearAll()
    ...
}

```