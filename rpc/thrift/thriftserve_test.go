package thrift

import (
	"context"
	"encoding/json"
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/adolphlxm/atc/rpc/thrift/gen/atcrpc"
	"github.com/adolphlxm/atc/rpc/thrift/gen/micro"
	"net"
	"testing"
	"time"
)

var (
	addr = net.JoinHostPort("127.0.0.1", "9191")
)

type AtcrpcHandler struct {
}

func (this *AtcrpcHandler) CallBack(a *atcrpc.ReqHandler, body map[string]string) (r string, err error) {
	jsonApi, err := json.Marshal(a)
	r = string(jsonApi)
	time.Sleep(3 * time.Second)
	fmt.Println("CallBack is running.")
	return
}

func (this *AtcrpcHandler) Ping() (err error) {
	return
}

// Local thrift testing.
func TestThriftRPC(t *testing.T) {
	_thriftServer(t, "binary", "framed")
}

// Remote thrift testing.
func TestRemoteThriftRPC(t *testing.T) {
	_remoteClient(t, "binary", "framed")
}


// Remote thrift pool client testing.
func TestRemoteThriftPoolThriftRPC(t *testing.T) {
	_remoteThriftPoolClient(t, "binary", "framed")
}

// Thrift serve.
func _thriftServer(t *testing.T, protocolT, transportT string) {
	// Initialize the thrift
	thriftServe := NewThriftServe(`{"addr":"` + addr + `","secure":"false"}`)
	thriftServe.Debug(false)
	thriftServe.Factory(protocolT, transportT)
	thriftServe.Timeout(3)

	// registerProcessor
	processor := atcrpc.NewAtcrpcThriftProcessor(&AtcrpcHandler{})
	thriftServe.RegisterProcessor("thriftTest", processor)

	// Run thrift serve
	err := thriftServe.Run()
	if err != nil {
		panic(err)
	}

	// Run client testing
	fmt.Println("thrift client start.")
	go _localClient(t, protocolT, transportT)

	time.Sleep(1 * time.Millisecond)
	ctx, _ := context.WithTimeout(context.Background(), 4*time.Second)
	if err := thriftServe.Shutdown(ctx); err != nil {
		fmt.Printf("thrift client %s finish.\n", err.Error())
	} else {
		fmt.Println("thrift client finish.")
	}

}

func _localClient(t *testing.T, protocolT, transportT string) {
	transport, err := thrift.NewTSocket(addr)
	if err != nil {
		t.Fatalf("[thrift client] Error resolving address:%v", err)
	}
	var (
		transportFactory thrift.TTransportFactory
		protocol         thrift.TProtocol
	)
	switch transportT {
	case "framed":
		transportFactory = thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	case "memorybuffer":
		transportFactory = thrift.NewTMemoryBufferTransportFactory(1000)
	case "buffered":
		transportFactory = thrift.NewTBufferedTransportFactory(1024)
	}

	useTransport := transportFactory.GetTransport(transport)
	switch protocolT {
	case "binary":
		protocol = thrift.NewTBinaryProtocol(useTransport, false, true)
	case "compact":
		protocol = thrift.NewTCompactProtocol(useTransport)
	case "json":
		protocol = thrift.NewTJSONProtocol(useTransport)
	case "simplejson":
		protocol = thrift.NewTSimpleJSONProtocol(useTransport)
	}


	if err := transport.Open(); err != nil {
		t.Fatalf("[thrift client] Error opening socket to %s,err:", addr, err)
	}
	defer transport.Close()

	rpc := &atcrpc.ReqHandler{
		Version: "V2",
		Method:  "GET",
		Handler: "users.test",
	}
	body := map[string]string{
		"a": "1",
		"b": "2",
	}

	mProtocol := thrift.NewTMultiplexedProtocol(protocol, "thriftTest")
	client := atcrpc.NewAtcrpcThriftClientProtocol(useTransport, mProtocol, mProtocol)
	r, err := client.CallBack(rpc, body)
	if err != nil {
		t.Fatalf("[thrift client] Error clinet CallBack err:", err)
	}
	rpcResult := &atcrpc.ReqHandler{}
	err = json.Unmarshal([]byte(r), rpcResult)
	if err != nil {
		t.Fatalf("[thrift client] Error clinet Unmarshal err:", err)
	}

	if rpcResult.Handler != rpc.Handler {
		t.Fatalf("[thrift client] Error result")
	}
}

/************************************/
/***** Remote service testing *******/
/************************************/
func _remoteClient(t *testing.T, protocolT, transportT string) {
	transport, err := thrift.NewTSocket(net.JoinHostPort("127.0.0.1", "9090"))

	if err != nil {
		t.Fatalf("[thrift client] Error resolving address:%v", err)
	}
	var (
		transportFactory thrift.TTransportFactory
		protocol         thrift.TProtocol
	)
	switch transportT {
	case "framed":
		transportFactory = thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	case "memorybuffer":
		transportFactory = thrift.NewTMemoryBufferTransportFactory(1000)
	case "buffered":
		transportFactory = thrift.NewTBufferedTransportFactory(1024)
	}

	useTransport := transportFactory.GetTransport(transport)
	switch protocolT {
	case "binary":
		protocol = thrift.NewTBinaryProtocol(useTransport, false, true)
	case "compact":
		protocol = thrift.NewTCompactProtocol(useTransport)
	case "json":
		protocol = thrift.NewTJSONProtocol(useTransport)
	case "simplejson":
		protocol = thrift.NewTSimpleJSONProtocol(useTransport)
	}

	mProtocol := thrift.NewTMultiplexedProtocol(protocol, "user")
	client := micro.NewMicroThriftClientProtocol(useTransport, mProtocol, mProtocol)
	if err := transport.Open(); err != nil {
		t.Fatalf("[thrift client] Error opening socket to %s,err:", addr, err)
	}
	defer transport.Close()

	body := map[string]string{
		"a": "1",
		"b": "2",
	}
	r, err := client.CallBack(1, "adolph", body)

	if err != nil {
		t.Fatalf("[thrift client] Error clinet CallBack err:", err)
	}
	fmt.Println(r)
}

/************************************/
/**** Remote pool client testing ****/
/************************************/
func _remoteThriftPoolClient(t *testing.T, protocolT, transportT string) {
	pool := NewThriftPool(net.JoinHostPort("127.0.0.1", "9090"),10,10,10)
	pool.Factory(protocolT, transportT)
	protocol := thrift.NewTMultiplexedProtocol(pool.GetTProtocol(),"user")

	client := micro.NewMicroThriftClientProtocol(pool.GetTransport(), protocol, protocol)
	body := map[string]string{
		"a": "1",
		"b": "2",
	}
	r, err := client.CallBack(1, "adolph", body)

	if err != nil {
		t.Fatalf("[thrift client] Error clinet CallBack err:", err)
	}
	fmt.Println(r)
}