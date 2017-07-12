package thrift

import (
	"github.com/adolphlxm/atc/pool"
	"github.com/adolphlxm/atc/rpc/thrift/lib/thrift"
	"time"
)

type Conn interface {
	// Close thrift connection
	Close() error

	GetTtransport() thrift.TTransport
	GetTprotocol() thrift.TProtocol
	NewTmultiplexedProtocol(serviceName string) *thrift.TMultiplexedProtocol
}

type ThriftPool struct {
	p           *pool.Pool
	c           *thrift.TSocket
	tprotocolF  thrift.TProtocolFactory
	ttransportF thrift.TTransportFactory

	ttransport thrift.TTransport
	tprotocol  thrift.TProtocol
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
		IdleTimeout: idleTimeout * time.Second,
	}
	return &ThriftPool{p: thriftPool}
}

func (this *ThriftPool) SetFactory(protocolFactory, transportFactory string) {
	switch transportFactory {
	case "framed":
		this.ttransportF = thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	case "memorybuffer":
		this.ttransportF = thrift.NewTMemoryBufferTransportFactory(1000)
	case "buffered":
		this.ttransportF = thrift.NewTBufferedTransportFactory(1024)
	}

	switch protocolFactory {
	case "binary":
		this.tprotocolF = thrift.NewTBinaryProtocolFactory(false, true)
	case "compact":
		this.tprotocolF = thrift.NewTCompactProtocolFactory()
	case "json":
		this.tprotocolF = thrift.NewTJSONProtocolFactory()
	case "simplejson":
		this.tprotocolF = thrift.NewTSimpleJSONProtocolFactory()
	}
}

func (this *ThriftPool) Get() (Conn, error) {
	c, err := this.p.Get()
	if err != nil {
		return nil, err
	}

	// Factory class used to create wrapped instance of Transports.
	// This is used primarily in servers, which get Transports from
	// a ServerTransport and then may want to mutate them (i.e. create
	// a BufferedTransport from the underlying base transport)
	ttransport := this.ttransportF.GetTransport(c.(*thrift.TSocket))
	tprotocol := this.tprotocolF.GetProtocol(ttransport)
	return &ThriftPool{p: this.p, c: c.(*thrift.TSocket), tprotocolF: this.tprotocolF, ttransportF: this.ttransportF, ttransport: ttransport, tprotocol: tprotocol}, err
}

func (this *ThriftPool) Close() error {
	return this.p.Put(this.c, false)
}

func (this *ThriftPool) GetTtransport() thrift.TTransport {
	return this.ttransport
}

func (this *ThriftPool) GetTprotocol() thrift.TProtocol {
	return this.tprotocol
}

func (this *ThriftPool) NewTmultiplexedProtocol(serviceName string) *thrift.TMultiplexedProtocol {
	return thrift.NewTMultiplexedProtocol(this.GetTprotocol(), serviceName)
}
