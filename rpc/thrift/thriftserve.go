// Copyright 2015 The ATC Authors. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// ATC is an open-source, automated test framework for the Go programming language.
// more infomation: http://atc.wiki
package thrift

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"strconv"
	"time"

	"git.apache.org/thrift.git/lib/go/thrift"
)

type Thrift interface {
	Run() error
	RegisterProcessor(name string, p thrift.TProcessor)
	Timeout(timeout int)
	Factory(protocol, transport string)
	Debug(ok bool)
	Addr() string

	Shutdown(ctx context.Context) error
}

type ThriftServe struct {
	debug     bool
	protocol  string
	transport string
	// TMultiplexedProcessor is a TProcessor allowing
	// a single TServer to provide multiple services.
	processor *thrift.TMultiplexedProcessor
	// Factory interface for constructing protocol instances.
	protocolFactory thrift.TProtocolFactory
	// Factory class used to create wrapped instance of Transports.
	transportFactory thrift.TTransportFactory

	addr    string
	secure  bool
	timeout time.Duration

	server *thrift.TSimpleServer

	quit chan struct{}
}

func NewThriftServe(config string) Thrift {
	var cf map[string]string
	json.Unmarshal([]byte(config), &cf)
	secure, _ := strconv.ParseBool(cf["secure"])

	return &ThriftServe{
		processor: thrift.NewTMultiplexedProcessor(),
		addr:      cf["addr"],
		secure:    secure,
	}
}

// Register processor of multiple interfaces.
func (t *ThriftServe) RegisterProcessor(name string, p thrift.TProcessor) {
	t.processor.RegisterProcessor(name, p)
}

// Run thrift serve
func (t *ThriftServe) Run() error {
	// Support the transport protocol factory.
	//
	// Format:
	//	TBinaryProtocol is binary format.
	//	TcompactProtocol is compressed format.
	//	TJSONProtocol is JSON format.
	//	TsimpleJSONProtocol is JSON write agreement only.
	//	TDebugProtocol is text format
	switch t.protocol {
	case "binary":
		t.protocolFactory = thrift.NewTBinaryProtocolFactory(false, true)
	case "compact":
		t.protocolFactory = thrift.NewTCompactProtocolFactory()
	case "json":
		t.protocolFactory = thrift.NewTJSONProtocolFactory()
	case "simplejson":
		t.protocolFactory = thrift.NewTSimpleJSONProtocolFactory()
	default:
		return errors.New("Invalid protocol.")
	}

	// If not thrift debug
	if t.debug {
		t.protocolFactory = thrift.NewTDebugProtocolFactory(t.protocolFactory, "ATC_logs_")
	}

	// Support the transport data factory.
	// Format:
	//	TFramedTransport
	//	TMemoryBuffer
	//	TBuffered
	transport := thrift.NewTTransportFactory()
	switch t.transport {
	case "framed":
		t.transportFactory = thrift.NewTFramedTransportFactory(transport)
	case "memorybuffer":
		t.transportFactory = thrift.NewTMemoryBufferTransportFactory(1024)
	case "buffered":
		t.transportFactory = thrift.NewTBufferedTransportFactory(10240)
	default:
		return errors.New("Invalid transport.")
	}

	return t.runServe()
}

func (t *ThriftServe) Debug(ok bool) {
	t.debug = ok
}

func (t *ThriftServe) Factory(protocol, transport string) {
	t.protocol = protocol
	t.transport = transport
}

func (t *ThriftServe) Timeout(timeout int) {
	t.timeout = time.Duration(timeout)
}

func (t *ThriftServe) Addr() string {
	return t.addr
}

func (t *ThriftServe) runServe() error {
	var tServer thrift.TServerTransport
	var err error
	if t.secure {
		cfg := new(tls.Config)
		if cert, err := tls.LoadX509KeyPair("server.crt", "server.key"); err == nil {
			cfg.Certificates = append(cfg.Certificates, cert)
		} else {
			return err
		}
		tServer, err = thrift.NewTSSLServerSocketTimeout(t.addr, cfg, t.timeout*time.Second)
	} else {
		tServer, err = thrift.NewTServerSocketTimeout(t.addr, t.timeout*time.Second)
	}

	if err != nil {
		return err
	}

	t.server = thrift.NewTSimpleServer4(t.processor, tServer, t.transportFactory, t.protocolFactory)

	go t.server.Serve()

	return nil
}

// shutdownPollInterval is how often we poll for quiescence
// during Server.Shutdown. This is lower during tests, to
// speed up tests.
// Ideally we could find a solution that doesn't involve polling,
// but which also doesn't have a high runtime cost (and doesn't
// involve any contentious mutexes), but that is left as an
// exercise for the reader.
var shutdownPollInterval = 500 * time.Millisecond

func (t *ThriftServe) Shutdown(ctx context.Context) error {
	err := t.server.Stop()
	ticker := time.NewTicker(shutdownPollInterval)
	for {
		if !t.server.ShuttingDown() {
			return err
		}

		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
		}
	}
}
