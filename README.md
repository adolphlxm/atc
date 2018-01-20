# ATC

ATC 是一个快速开发GO应用程序的开源框架，支持RESTful API, Thrift RPC, Redis, Nats队列的框架.可根据自身业务逻辑选择性的卸载中间件的功能，均支持平滑退出。

要求GO版本 >= 1.8
当前版本: 0.9.6 (Beta 2018-01-20)

## 安装ATC

    go get github.com/adolphlxm/atc
    
  **用到的第三方 go package TAG**
  
  ```config 
    // 配置文件加载包,之前写的一个Github账户，后续会迁移过来，方便管理
    github.com/lxmgo/config
    
    // 官方websocket包
    code.google.com/p/go.net/websocket
    
    // xorm 包
    github.com/go-sql-driver/mysql
    github.com/go-xorm/xorm
    github.com/go-xorm/core
  ```
   
## RESTful API 经典案例
一个经典的ATC例子 `atc.go`
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

// 登录过滤器
// 可以通过自定义过滤器来实现登录状态、权限检查等功能
func AfterLogin(ctx *context.Context){
	// 错误输出
	error := atc.NewError(ctx)
	error.Code(401,10000).JSON()
}

func init(){
        // 分组版本控制
        
        v1 := atc.Route.Group("V1")
        {
                // V1版本路由
                v1.AddRouter("users.login",&LoginHandler{})
                ...
        }
        V2 := atc.Route.Group("V2")
        {
                // V2版本过滤器, 根据路由规则加载。
                // 支持三种过滤器：
                //      1. EFORE_ROUTE                    //匹配路由之前
                //      2. BEFORE_HANDLER                 //匹配到路由后,执行Handler之前
                //      3. AFTER                          //执行完所有逻辑后
                v2.AddFilter(atc.BEFORE_ROUTE,"users.*",AfterLogin)
                // V2版本路由
                v2.AddRouter("users.login",&LoginHandler{})
                ...
        }
}
```
    
然后在浏览器访问 

1. `http://localhost/V1/users/login`
    * 将会得到一个json返回

2. `http://localhost/V2/users/login`
    * 先执行`BFORE_ROUTE`过滤器, 未通过则得到一个json返回。
    * 通过过滤器后 加载 `Get()`, 将会得到一个json返回。

## 路由规则

**固定路由**

固定路由也就是全匹配的路由，如下所示：

`注：ATC路由模块之间请使用"."来区分，而不是"/"`
```router
// 匹配根目录
atc.AddRouter(".", &index.MainHandler{})

// 匹配 /api
atc.AddRouter("api",&api.IndexHandler{})

// 匹配 /api/list
atc.AddRouter("api.list", &api.ListHandler{})

// 匹配 /api/group
atc.AddRouter("api.group", &api.GroupHandler{})
```

如上所示的固定路由就是我们最常用的路由方式，根据用户请求方法不同请求控制器中对应的方法，典型的 RESTful 方式。

**正则路由**

 为了更加方便的路由设置，ATC 支持多种方式的路由：

* {id:[\\d]`*`}, 其中 '*' 表示匹配{}表达式零次或多次
* {id:[\\d]`?`}, 其中 '?' 表示匹配{}表达式零次或一次
* {id:[\\d]`+`}, 其中 '+' 表示匹配{}表达式一次或多次

```router
// 匹配 /api/123 , 同时可匹配 /api | /api/ 这两个URL
        // id = 123 , 可通过this.Ctx.ParamIndex64("id") 获取参数值
        // 匹配{id:[0-9]*}表达式零次或多次
atc.AddRouter("api.{id:[0-9]*}",&api.IndexHandler{})

// 匹配 /api/123, 不匹配 /api
atc.AddRouter("api.{id:[0-9]+}",&api.IndexHandler{})
atc.AddRouter("api.{id}",&api.IndexHandler{})

// 自定义正则匹配, 匹配 /api/adolph , name = adolph
atc.AddRouter("api.{name:[\w]+}",&api.IndexHandler{})
```

## 顺序平滑退出

使用双向链表list实现

* 客户端实现TT接口

        type TT interface {
        	ModuleID() string
        	Stop() error
        }
     - ModuleID() 方法返回string, 表示退出模块名称
     - Stop() 方法根据客户端业务自行实现退出逻辑

* 通过atc.Grace...() 方法插入实现的接口struct
* ATC默认在队头插入了`http`,`thrift`,`queuePublisher`,`queueConsumer`四个TT退出接口，也可根据客户端需要自行控制退出顺序。
* 当ATC收到退出信号时(参阅信号说明)，从队尾开始,逆向顺序逐个退出

**使用方法**

        // 双向链表队头插入退出接口
        atc.GracePushFront(TT)
        // 双向链表队尾插入退出接口
        atc.GracePushBack(TT)
        // 在"atc"模块的链表之后插入退出接口
        atc.GraceInsertAfter("atc",TT)
        // 在"atc"模块的链表之前插入退出接口
        atc.GraceInsertBefore("atc",TT)
        // 在链表中移除"atc"退出接口
        atc.GraceRemove("atc")
        // 在链表中将"atc"接口移动到"http"之后。
        atc.GraceMoveAfter("atc", "http")
        // 在链表中将"atc"接口移动到"http"之前。
        atc.GraceMoveAfter("atc", "http")

## Queue 队列
支持redis、nats两个队列

* atc.QueuePublisher [生产实例]
* QueueConsumer [消费实例]
* 实例初始化配置可在`app.ini`内写

**配置参数**

* queue.publisher.support
    - 参数：true | false
    - 注释：是否支持开启生产实例
* queue.publisher.`aliasnames`
    - 参数: 别名,逗号隔开,如:p1,p2...
    - 注释：用来初始化多实例
* queue.publisher.`[aliasname]`.driver
    - 参数: redis | nats
    - 注释：支持的队列驱动
* queue.publisher.`[aliasname]`.addrs
    - 参数：redis://127.0.0.1:6379
    - 注释：redis连接地址
* queue.consumer.support
    - 参数：true | false
    - 注释：是否支持开启消费者实例

**使用方法**

```go
// Publisher
atc.GetPublisher("aliasname").Publish("subject", &message.Message{
    Body: util.MustMessageBody(nil, /* point to your protobuffer struct */ ),
})


// Consumer
sub,_ := atc.GetPublisher("aliasname").Subscribe("subject", "cluster-group")
msg, _ := sub.NextMessage(time.Second)
// logic for msg
```


## RPC 经典案例

**Thrift RPC**
关于Thrift RPC 具体可以 度娘、谷爹查看

[之前写过一篇Thrift简单使用教程：GO/PHP使用指南](http://blog.csdn.net/liuxinmingcode/article/details/45696237)

[获取Thrfit-官方](http://thrift.apache.org)

[官方各种DEMO](https://git1-us-west.apache.org/repos/asf?p=thrift.git;a=tree;f=tutorial;h=d69498f9f249afaefd9e6257b338515c0ea06390;hb=HEAD)

[协议库IDL文件参考资料](https://my.oschina.net/helight/blog/195015)

这里就不累赘描述Thrift具体用法了。

* Go的Thrift包已经整合在ATC框架内了，无需重新安装和下载了！
* 由于Go的Thrift源码有修改，支持RPC平滑退出
* 使用.thrift IDL 生成 .go 请使用 `atc-tool` 工具(不然可能会报错.)

[atc-tool 工具](https://github.com/adolphlxm/atc-tool)

Thrift RPC 路由 `router.go`
```go
func init() {
	processor := micro.NewMicroThriftProcessor(&MicroHandler{})
	atc.ThriftRPC.RegisterProcessor("user", processor)
}
```

**gRPC...**

## ORM
* 通过`app.ini`配置文件加载多库初始化方法
* 返回 `orm interface ` 接口

```go
    atc.DB("库名").Where("id=?",1).Get(...)
```

## Cache
* 通过`app.ini`配置redis,memcache

``` go
atc.GetCache("别名").Put()
atc.GetCache("别名").Put()
atc.GetCache("别名").Get()
atc.GetCache("别名").Delete()
...
```

## Mongodb
* 通过`app.ini`配置文件加载多库初始化方法

```go
    atc.GetMgoDB("别名").Insert()
```

## 日志处理
**通过`app.ini`配置日志输出引擎**

**参数说明**

* logdir 存放日志路径
* maxsize 日志分割文件尺寸, 单位:byte
* buffersize 日志缓冲区大小，单位:byte
* flushinterval 定时刷新日志到磁盘的间隔时间，单位:s

**级别说明(枚举)**

0. LevelFatal = iota 打印文件名及行号
1. LevelError 打印文件名及行号
2. LevelWarn 打印文件名及行号
3. LevelNotice 不打印文件名及行号
4. LevelInfo 不打印文件名及行号
5. LevelTrace 不打印文件名及行号
6. LevelDebug 打印文件名及行号


**通用使用方式**

```go
    logs.Debug("")
    logs.Debugf("%v", "debugf")
    logs.Info("")
    logs.Infof("%v", "infof")
    logs.Warn("")
    logs.Warnf("%v", "warnf")
    logs.Error("")
    logs.Errorf("%v", "errorf")
    ...
    logs.Flush()
```

## 编译并运行

    go build atc.go
    ./atc
    
   Flag参数说明：
   ```config 
     -c string
           use -c <config file> (default "conf/app.ini")
     -m string
           Use -m <config mode> (default "dev")
     -v    Use -v <current version>
   ```
   
   ATC信号控制：
   
| 信号量  | 退出 |
|:------------- | -------------:|
| TERM,INT      | 立即终止 |
| QUIT      | 优雅的关闭进程,即等所有请求结束后再关闭 |


## ATC项目结构
<pre>
├── conf
│   ├── app.ini
│   └── error.ini
├── front
│   └── HTML...
├── bin
├── src
│   ├── httprouter
│         ├── V1
│         └── router.go
│   └── thriftrpc
│         ├── idl
│         ├── gen-go
│         ├── ...(.go)
│         └── router.go
└── atc.go
</pre>

## 特性

* 支持RESTful HTTP通信 及 平滑退出
* 支持Websoeckt通信
* 支持RPC通信 及 平滑退出
    - Thrift 
    
## 更新日志
* 2017.5 
    - 日志支持file文件写入,通过`app.ini` 配置日志类型
    - 支持Flag参数(-c 配置文件, -m 配置环境, -v 当前版本号)
    - 优雅的关闭进程
    - 优化错误码配置文件加载,通过`app.ini` 配置错误码文件
    - 优化DEBUG模式
* 2017.6
    - utils/encrypt包增加RSA/DES/AES 加解密
    - 修复ORM BUG
    - rpc/thrift 增加 thrift client 实现, 结合pool对连接进行管理(线程安全，thrift的本身client端是线程不安全的)
    - pool包 通用连接池管理
    - 优化thrift RPC，支持atc-tool生成工具，使用更方便
    - 支持 API 跨域配置, 通过`app.ini` 配置跨域
* 2017.7
    - 修复不同环境配置`app.ini`不生效 BUG
    - 优化日志模块 `logs` 包，使用更方便
    - 重构thrift client pool 封装，修复若干BUG
    - 修复日志级别BUG
* 2017.8
    - 增加cache缓存模块
    - 支持memcache
* 2017.10
    - 优化路由，Context struct增加`RunHandler`字段，便于在REST中获取和过滤运行Handler
    - 优化JSON输出，变更前仅支持 this.Ctx.SetDatas( map[string]interface{}) -> this.JSON() ，变更后兼容之前使用方式同时支持 this.JSON(interface{})
* 2017.11
    - cache缓存模块支持redis
* 2017.12
    - 优化RESTFul router 正则匹配
    - Context提供更多的路由正则参数解析方法 及 更多的表单解析方法
    - 增加客户端顺序平滑退出接口
    - 优化logs包，日志写入缓冲区，定时刷新缓存区到磁盘(也可以退出程序时，调用logs.Flush()方法主动刷取)及配置文件
    - 增加了queue队列封装包
* 2018.1
    - ATC框架管理orm初始化实例
    - 增加cache,queue多实例管理
    - 增加mongodb封装，多实例管理
    - 优化初始化

## 即将支持特性(待定稿)

* HTTPS / POST附件上传压缩
* gRPC
* 服务注册与发现模块

# LICENSE

ATC is Licensed under the Apache License, Version 2.0 (the "License")
(http://www.apache.org/licenses/LICENSE-2.0.html).