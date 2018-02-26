package pgrpc

import (
	"net"
	"google.golang.org/grpc"
)

type Grpc interface {
	NewServer(network, addr string) error
	GetServer() *grpc.Server
	Serve() error
	GetServiceInfo() map[string]grpc.ServiceInfo
	GracefulStop()
}

type GrpcServe struct {
	lis net.Listener
	server *grpc.Server
}

func NewGrpc() Grpc {
	return &GrpcServe{}
}

// NewServer creates a gRPC server which has no service registered and has not
// started to accept requests yet.
func (this *GrpcServe) NewServer(network, addr string) error {
	var err error
	this.lis, err = net.Listen(network, addr)
	if err != nil {
		return err
	}

	this.server = grpc.NewServer()

	return nil
}

func (this *GrpcServe) GetServer() *grpc.Server {
	return this.server
}

func (this *GrpcServe) Serve() error {
	return this.server.Serve(this.lis)
}

func (this *GrpcServe) GetServiceInfo() map[string]grpc.ServiceInfo {
	return this.server.GetServiceInfo()
}

func (this *GrpcServe) GracefulStop() {
	this.server.GracefulStop()
}