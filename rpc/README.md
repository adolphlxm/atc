# RPC
RPC引擎目前支持Thrift(client & serve)

# Thrift安装

    go get github.com/adolphlxm/atc/rpc/thrift
   
# Thrift服务端使用步骤

## 第一步：引入包
   
    import(
        "github.com/adolphlxm/atc/rpc/thrift"
    )
    
## 第二步：初始化服务
    
    // 创建服务实例
    ThriftRPC := thrift.NewThriftServe(`{"addr":"127.0.0,1:9090"}`)
    // DEBUG
	ThriftRPC.Debug(true)
	// 设置传输层、协议层
	ThriftRPC.Factory("binary", "framed")
	// 设置超时时间
	ThriftRPC.Timeout(10)

## 第三步：启动服务

    ThriftRPC.Run()
    
## 第四步：根据逻辑设置平滑退出

    ctx, _ := context.WithTimeout(context.Background(), time.Duration(Aconfig.ThriftQTimeout)*time.Second)
    ThriftRPC.Shutdown(ctx)
    
## 开始使用

Thrift RPC 路由 `router.go`
```go
func init() {
	processor := micro.NewMicroThriftProcessor(&MicroHandler{})
	atc.ThriftRPC.RegisterProcessor("user", processor)
}
```

* Go的Thrift包已经整合在ATC框架内了，无需重新安装和下载了！
* 使用.thrift IDL 生成 .go 请使用 `atc-tool` 工具(不然可能会报错.)

thrift命令行工具，用于 thrift IDL .go 文件生产

    $ atc-tool thrift [options] file
    
具体thrift命令可使用 `thrift --help` 查看
[atc-tool 工具](https://github.com/adolphlxm/atc-tool)

举例：

    $ atc-tool thrift -r --gen go xxx.thrift
    
# Thrift客户端使用步骤
## 第一步：引入包
    
    import(
        "github.com/adolphlxm/atc/rpc/thrift"
    )
    
## 第二步：初始化客户端连接池(由于thrift client 是非线程安全，So 提供了一个连接池管理)

    pool := NewThriftPool(net.JoinHostPort("127.0.0.1", "9090"),10,10,10)
    pool.SetFactory("binary", "framed")
    
## 第三步：开始使用
    
    // 获取一个可用连接
    conn,err := pool.Get()
    if err != nil {
        return err
    }
    mulProtocol := conn.NewTmultiplexedProtocol("user")
    
    // 这部分调用自动生成代码
    client := micro.NewMicroThriftClientProtocol(conn.GetTtransport(), conn.GetTprotocol(), mulProtocol)