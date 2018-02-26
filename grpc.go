package atc

import (
	"net/url"

	"google.golang.org/grpc"
	"github.com/adolphlxm/atc/rpc/pgrpc"
	"github.com/adolphlxm/atc/logs"
)

var grpcserve pgrpc.Grpc

func lazyInitGrpcServer() error {
	addrs := AppConfig.String("grpc.addrs")
	logs.Trace("grpc.serve:starting...")
	addrUrl, err := url.Parse(addrs)

	grpcserve = pgrpc.NewGrpc()
	err = grpcserve.NewServer("tcp", addrUrl.Host)
	if err != nil {
		logs.Fatalf("grpc.serve:start addrs fail err:%s", err.Error())
		panic(err)
	}

	return nil
}

func runGrpcServe() {
	end := make(chan struct{})
	go func() {
		close(end)
		err := grpcserve.Serve()
		if err != nil {
			logs.Fatalf("grpc.serve:start fail err:%s", err.Error())
			panic(err)
		}
	}()
	<-end

	GracePushFront(&grpcServeShutDown{})
	// TODO serviceInfo
	addrs := AppConfig.String("grpc.addrs")
	addrUrl, _ := url.Parse(addrs)
	logs.Tracef("grpc.serve:Running on %s.", addrUrl.Host)
}


type grpcServeShutDown struct {}
func (this *grpcServeShutDown) ModuleID() string {
	return "pgrpc"
}
func (this *grpcServeShutDown) Stop() error {
	grpcserve.GracefulStop()
	return nil
}

func GetGrpcServer() *grpc.Server{
	return grpcserve.GetServer()
}