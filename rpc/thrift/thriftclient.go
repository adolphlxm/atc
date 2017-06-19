package thrift

import (
	"github.com/adolphlxm/atc/rpc/thrift/lib/thrift"
	"github.com/adolphlxm/atc/pool"
	"time"
)

type ThriftPool struct {
	pool      *pool.Pool
	protocolFactory  string
	transportFactory string
}


func NewThriftPool(addr string, maxActive, maxIdle int, idleTimeout time.Duration) *ThriftPool {
	thriftPool := &pool.Pool{
		Dial: func() (interface{}, error) {
			transport, err := thrift.NewTSocket(addr)
			if err != nil {
				return nil, err
			}
			if err = transport.Open(); err != nil {
				return nil, err
			}
			return transport, err
		},
		Close: func(c interface{}) error {
			c.(*thrift.TSocket).Close()
			return nil
		},
		MaxActive:   maxActive,
		MaxIdle:     maxIdle,
		IdleTimeout: idleTimeout,
	}
	return &ThriftPool{pool: thriftPool}
}

func (this *ThriftPool) SetFactory(protocolFactory, transportFactory string) {
	this.protocolFactory = protocolFactory
	this.transportFactory = transportFactory
}


func (this *ThriftPool) GetTtransport() (thrift.TTransport,error) {
	var transportFactory thrift.TTransportFactory
	switch this.transportFactory {
	case "framed":
		transportFactory = thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	case "memorybuffer":
		transportFactory = thrift.NewTMemoryBufferTransportFactory(1000)
	case "buffered":
		transportFactory = thrift.NewTBufferedTransportFactory(1024)
	}

	v, err := this.pool.Get()
	if err != nil {
		return nil, err
	}
	// Factory class used to create wrapped instance of Transports.
	// This is used primarily in servers, which get Transports from
	// a ServerTransport and then may want to mutate them (i.e. create
	// a BufferedTransport from the underlying base transport)

	return transportFactory.GetTransport(v.(thrift.TTransport)),nil
}

func (this *ThriftPool) GetTprotocol(ttransport thrift.TTransport) (tprotocol thrift.TProtocol) {
	switch this.protocolFactory {
	case "binary":
		tprotocol = thrift.NewTBinaryProtocol(ttransport, false, true)
	case "compact":
		tprotocol = thrift.NewTCompactProtocol(ttransport)
	case "json":
		tprotocol = thrift.NewTJSONProtocol(ttransport)
	case "simplejson":
		tprotocol = thrift.NewTSimpleJSONProtocol(ttransport)
	}
	return
}

func (this *ThriftPool) NewTmultiplexedProtocol(serviceName string,ttransport thrift.TTransport) *thrift.TMultiplexedProtocol {
	return thrift.NewTMultiplexedProtocol(this.GetTprotocol(ttransport),serviceName)
}