# cache

缓存模块，目前支持的引擎有memcache

# 安装

    go get github.com/adolphlxm/atc/cache
   
# 使用步骤
##引入包
    
    import(
        "github.com/adolphlxm/atc/cache"
    )
    
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