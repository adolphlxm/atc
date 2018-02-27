# RPC
RPC引擎目前支持
* Grpc
* Thrift(client & serve)

# Grpc安装

## 第一步：引入包
    github.com/adolphlxm/atc/rpc/pgrpc

## 第二步：初始化服务

        grpcserve = pgrpc.NewGrpc()
    	err = grpcserve.NewServer("tcp", "localhost:50005")
    	if err != nil {
    		logs.Fatalf("grpc.serve:start addrs fail err:%s", err.Error())
    		panic(err)
    	}

## 第三步：简单GRPC服务端实现及注册服务
        type server struct {}

        func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
        	return &pb.HelloReply{Message: "Hello " + in.Name}, nil
        }

        pb.RegisterGreeterServer(grpcserve.GetServer(), &server{})

## 第四步：运行GRPC serve

        grpcserve.Serve()

## 第五步：编写GRPC client 代码测试

```go
package main

//client.go

import (
	"log"
	"os"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	pb "example/rpc/helloworld"
)

const (
	address     = "127.0.0.1:50005"
	defaultName = "world"
)


func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatal("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	name := defaultName
	if len(os.Args) >1 {
		name = os.Args[1]
	}
	r, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: name})
	if err != nil {
		log.Fatal("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.Message)

}
```

# <del>Thrift安装<del>

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