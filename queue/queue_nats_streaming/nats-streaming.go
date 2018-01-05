package queue_nats_streaming

import (
	"github.com/adolphlxm/atc/queue"

	"github.com/nats-io/go-nats-streaming"
)

// Driver for nats-streaming
type NatsStreamingQueueDriver struct{}

// dsn format is `nats://[user:pass@host1:port],[user2:pass@host2:port2]/?cluster=xxx&client=xx&options`
//
// cluster     - The NATS Streaming cluster ID, required
// client      - unique client id, required
//
// all timeout option's unit is ms
// connectTimeout     - default is 2000ms, timeout for establishing a connection
// ackTimeout         - default is 30000ms, timeout for waiting for an ACK for a published message.
// maxPubAcksInflight - default is 16384, maximum number of published messages
//                      without outstanding ACKs from the server
//
func (d *NatsStreamingQueueDriver) Open(addr string) (queue.Conn, error) {
	info, err := parseURL(addr)
	if err != nil {
		return nil, err
	}

	nc, err := stan.Connect(info.ClusterID, info.ClientID, info.Options...)
	if err != nil {
		return nil, err
	}
	// here we need a block channel
	return &natsStreamingQueueConn{conn: nc, ch:make(chan *stan.Msg, 0)}, nil
}

func init() {
	queue.Register("nats-streaming", &NatsStreamingQueueDriver{})
}
