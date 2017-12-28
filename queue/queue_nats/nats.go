package queue_nats

import (
	"github.com/adolphlxm/atc/queue"
	"github.com/nats-io/go-nats"
)

type NatsQueueDriver struct{}

// dsn format is `nats://[user:pass@host1:port],[user2:pass@host2:port2]/?options`
//
// options can be:
// allowReconnect - default is true
// name           - unique client name, default is random
// noRandomize    - if true, then turn off randomizing the server pool, default is false
// maxReconnect   - maxReconnects is an Option to set the maximum number of reconnect attempts
//
// all timeout option's unit is ms
// connectTimeout - default is 2000ms, connectTimeout is an Option to set the timeout for Dial on a connection.
//
func (d *NatsQueueDriver) Open(addr string) (queue.Conn, error) {
	info, err := parseURL(addr)
	if err != nil {
		return nil, err
	}

	nc, err := nats.Connect(info.Url, info.Options...)
	if err != nil {
		return nil, err
	}
	return &natsQueueConn{conn: nc}, nil
}

func init() {
	queue.Register("nats", &NatsQueueDriver{})
}
