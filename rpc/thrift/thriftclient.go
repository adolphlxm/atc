package thrift

import (
	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/adolphlxm/atc/pool"
	"time"
)

type ThriftPool struct {
	pool      *pool.Pool
	protocol  thrift.TProtocol
	transport thrift.TTransport
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

func (this *ThriftPool) Factory(protocol, transport string) error {
	var transportFactory thrift.TTransportFactory
	switch transport {
	case "framed":
		transportFactory = thrift.NewTFramedTransportFactory(thrift.NewTTransportFactory())
	case "memorybuffer":
		transportFactory = thrift.NewTMemoryBufferTransportFactory(1000)
	case "buffered":
		transportFactory = thrift.NewTBufferedTransportFactory(1024)
	}

	v, err := this.pool.Get()
	if err != nil {
		return err
	}

	this.transport = transportFactory.GetTransport(v.(thrift.TTransport))

	switch protocol {
	case "binary":
		this.protocol = thrift.NewTBinaryProtocol(this.transport, false, true)
	case "compact":
		this.protocol = thrift.NewTCompactProtocol(this.transport)
	case "json":
		this.protocol = thrift.NewTJSONProtocol(this.transport)
	case "simplejson":
		this.protocol = thrift.NewTSimpleJSONProtocol(this.transport)
	}
	return nil
}

func (this *ThriftPool) GetTProtocol() thrift.TProtocol {
	return this.protocol
}

func (this *ThriftPool) GetTransport() thrift.TTransport {
	return this.transport
}
