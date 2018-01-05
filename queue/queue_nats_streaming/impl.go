package queue_nats_streaming

import (
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/adolphlxm/atc/queue"
	"github.com/adolphlxm/atc/queue/message"
	"github.com/adolphlxm/atc/queue/util"

	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes"
	"github.com/nats-io/go-nats-streaming"
)

var (
	ErrNotSupport = errors.New("queue nats-streaming: Not Support")
)

type natsStreamingQueueConn struct {
	conn stan.Conn
	ch   chan *stan.Msg

	// subscribe init once
	sbErr     error
	sub       stan.Subscription
	isSubInit uint32
	m         sync.Mutex
}

func (c *natsStreamingQueueConn) Subscribe(subject, group string) (queue.Subscriber, error) {
	return nil, ErrNotSupport
}

func (c *natsStreamingQueueConn) RpcHandle(subject, group string, handler queue.RpcHandler) {
	panic(ErrNotSupport)
}

func (c *natsStreamingQueueConn) Publish(subject string, msg *message.Message) error {
	return ErrNotSupport
}

func (c *natsStreamingQueueConn) Request(subject string, req *message.RpcMessage, timeout time.Duration) (*message.RpcMessage, error) {
	return nil, ErrNotSupport
}

func (c *natsStreamingQueueConn) Enqueue(subject string, msg *message.Message) error {
	if msg == nil {
		return errors.New("nil message to enqueue")
	}
	if msg.MessageId == "" {
		msg.MessageId = util.GenMsgID()
	}
	buf, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	return c.conn.Publish(subject, buf)
}

func (c *natsStreamingQueueConn) initSubscribe(subject, group string, timeout time.Duration) {
	if atomic.LoadUint32(&c.isSubInit) == 1 {
		return
	}

	c.m.Lock()
	if c.isSubInit == 0 {
		c.sub, c.sbErr = c.conn.QueueSubscribe(
			subject,
			group,
			func(msg *stan.Msg) {
				c.ch <- msg
			},
			stan.DurableName(fmt.Sprintf("%s-%s", subject, group)),
		)
		atomic.StoreUint32(&c.isSubInit, 1)
	}
	c.m.Unlock()
}

func (c *natsStreamingQueueConn) Dequeue(subject, group string, timeout time.Duration, dst proto.Message) (*message.Meta, error) {
	c.initSubscribe(subject, group, timeout)
	if c.sbErr != nil {
		return nil, c.sbErr
	}

	meta := &message.Meta{}
	natsMsg := <-c.ch

	ret := &message.Message{}
	err := proto.Unmarshal(natsMsg.Data, ret)
	if err != nil {
		return nil, err
	}
	meta.FormMessage(ret)
	meta.Src = natsMsg.Subject

	err = ptypes.UnmarshalAny(ret.Body, dst)
	return meta, err
}

func (c *natsStreamingQueueConn) Close() error {
	if c.sub != nil {
		c.sub.Close()
		close(c.ch)
	}
	return c.conn.Close()
}
