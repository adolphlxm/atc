# ATC

ATC 是一个快速开发GO应用程序的开源框架，支持RESTful API 及 Thrift RPC的框架.可根据自身业务逻辑选择性的卸载中间件的功能，均支持平滑退出。

当前版本: 0.1.0 (Beta 2017-04-28)

More info [atc.wiki](http://atc.wiki)

[老版本GITHUB](https://github.com/lxmgo)

## 安装ATC

    go get github.com/adolphlxm/atc
    go get github.com/lxmgo/config
    
   用到了配置加载包 `github.com/lxmgo/config`, 该GITHUB内容将逐步迁移过来(为了方便管理新开了Adolphlxm)
   
## RESTful API 经典案例
一个经典的ATC例子 `main.go`
```go
package main

import (
	"github.com/adolphlxm/atc"
)

type LoginHandler struct {
        atc.Handler
}

func (this *LoginHandler) Get(){
        // 已登录
        if true {
               loginData := map[string]interface{}{
                       "username" : "ATC",
                       "regtime":"2017-04-28",
               }
               this.Ctx.SetData("data",loginData)
               this.Ctx.SetData("ID",this.Ctx.Query("id"))
               this.JSON()
               return 
        }
        
        // 未登录
        // 自定义错误提示内容
        this.Error406(-1).Message("没有权限查看").JSON()
        // error.ini错误码匹配的提示内容
        // this.Error406(-1).JSON()
        return 
}

func main(){
    // 根据配置文件注入依赖中间件
    // 目前支持：HTTP/Websoeckt、Thrift、ORM(xorm、其它待开发)
	atc.Run()
}
```

路由加载 `router.go`

```go
func init(){
        // 分组版本
        v1 := atc.Route.Group("V1")
        {
                // 该分组的版本
                v1.AddRouter("users",&LoginHandler{})
        }
}
```
    
然后在浏览器访问 `http://localhost/users/login`, 将会得到一个json返回

##RPC 经典案例

### Thrift RPC
关于Thrift RPC 具体可以 度娘、谷爹查看

[之前写过一篇Thrift简单使用教程：GO/PHP使用指南](http://blog.csdn.net/liuxinmingcode/article/details/45696237)

[获取Thrfit-官方](http://thrift.apache.org)

[官方各种DEMO](https://git1-us-west.apache.org/repos/asf?p=thrift.git;a=tree;f=tutorial;h=d69498f9f249afaefd9e6257b338515c0ea06390;hb=HEAD)

[协议库IDL文件参考资料](https://my.oschina.net/helight/blog/195015)

这里就不累赘描述Thrift具体用法了。

Go的Thrift包安装：
```go
go get git.apache.org/thrift.git/lib/go/thrift
```

注：为了实现Thrift RPC平滑退出，改了Thrift Go源码，安装后需要重新覆盖下源码。
复制 `github.com/adolphlxm/atc/rpc/thrift/lib/thrift` 目录下的源码到 `git.apache.org/thrift.git/lib/go/thrift` 即可。

Thrift RPC 路由 `router.go`
```go
func init() {
	processor := micro.NewMicroThriftProcessor(&MicroHandler{})
	atc.ThriftRPC.RegisterProcessor("user", processor)
}
```

### gRPC...

## 特性

* 支持RESTful HTTP通信 及 平滑退出
* 支持Websoeckt通信
* 支持RPC通信 及 平滑退出
    - Thrift 
    
## 即将支持特性(待定稿)

* gRPC通信
* 消息通信
* 服务注册与发现模块

文档陆续更新中...
DEMO更新中...

# LICENSE

ATC is Licensed under the Apache License, Version 2.0 (the "License")
(http://www.apache.org/licenses/LICENSE-2.0.html).