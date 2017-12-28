package queue

import (
	"fmt"
	"io"
	"sort"
	"sync"
	"time"

	"github.com/adolphlxm/atc/queue/message"
	"github.com/golang/protobuf/proto"
)

type RPCRequest struct {
	reply string
	Data  proto.Message
}

type RPCResponse struct {
	reply string
	Data  proto.Message
}

type RpcHandler func(*message.RpcMessage, *message.RpcMessage)

type basePublisher interface {
	Publish(subject string, msg *message.Message) error
    Request(subject string, req *message.RpcMessage, timeout time.Duration) (*message.RpcMessage, error)
}

type Publisher interface {
	basePublisher
	io.Closer
}

type Subscriber interface {
	NextMessage(timeout time.Duration) (*message.Message, error)
}

type baseConsumer interface {
	Subscribe(subject, group string) (Subscriber, error)
	RpcHandle(subject, group string, handler RpcHandler)
}

type Consumer interface {
	baseConsumer
	io.Closer
}

type Conn interface {
	basePublisher
	baseConsumer
	io.Closer
}

type driver interface {
	Open(addr string) (Conn, error)
}

var (
	drivers   = make(map[string]driver)
	driversMu sync.RWMutex
)

// Register makes a queue driver available by the provided name.
// If Register is called twice with the same name or if driver is nil,
// it panics.
func Register(name string, d driver) {
	driversMu.Lock()
	defer driversMu.Unlock()
	if d == nil {
		panic("queue: Register driver is nil")
	}

	if _, dup := drivers[name]; dup {
		panic("queue: Register called twice for driver " + name)
	}
	drivers[name] = d
}

// Drivers returns a sorted list of the names of the registered drivers.
func Drivers() []string {
	driversMu.RLock()
	defer driversMu.RUnlock()
	var list []string
	for name := range drivers {
		list = append(list, name)
	}
	sort.Strings(list)
	return list
}

func NewPublisher(driverName, queueAddrs string) (Publisher, error) {
	driversMu.RLock()
	d, ok := drivers[driverName]
	driversMu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("queue: unknown driver %q (forgotten import?)", driverName)
	}

	return d.Open(queueAddrs)
}

func NewConsumer(driverName, queueAddrs string) (Consumer, error) {
	driversMu.RLock()
	d, ok := drivers[driverName]
	driversMu.RUnlock()

	if !ok {
		return nil, fmt.Errorf("queue: unknown driver %q (forgotten import?)", driverName)
	}

	return d.Open(queueAddrs)
}
