package queue_nats

import (
	"errors"
	"time"

	"github.com/adolphlxm/atc/queue"
	"github.com/adolphlxm/atc/queue/message"
	"github.com/adolphlxm/atc/queue/util"

	"github.com/golang/protobuf/proto"
	"github.com/nats-io/go-nats"
)

var (
	errRpcUnmarshal []byte
	errRpcMarshal   []byte
)

func init() {
	buf, err := proto.Marshal(&message.RpcMessage{Code: 1})
	if err != nil {
		panic(err)
	}
	errRpcUnmarshal = buf

	buf, err = proto.Marshal(&message.RpcMessage{Code: 2})
	if err != nil {
		panic(err)
	}
	errRpcMarshal = buf
}

type natsSubscriber struct {
	subject string
	pbConn  *nats.Subscription
}

func (s *natsSubscriber) NextMessage(timeout time.Duration) (*message.Message, error) {
	msg, err := s.pbConn.NextMsg(timeout)
	if err != nil {
		return nil, err
	}
	ret := &message.Message{}
	if err = proto.Unmarshal(msg.Data, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

type natsQueueConn struct {
	conn *nats.Conn
}

func (c *natsQueueConn) Subscribe(subject, group string) (queue.Subscriber, error) {
	sub, err := c.conn.QueueSubscribeSync(subject, group)
	if err != nil {
		return nil, err
	}
	return &natsSubscriber{pbConn: sub, subject: subject}, nil
}

func (c *natsQueueConn) RpcHandle(subject, group string, handler queue.RpcHandler) {
	c.conn.QueueSubscribe(subject, group, func(msg *nats.Msg) {
		req := &message.RpcMessage{}
		if err := proto.Unmarshal(msg.Data, req); err != nil {
			c.conn.Publish(msg.Reply, errRpcUnmarshal)
			c.conn.Flush()
			return
		}
		ret := &message.RpcMessage{MessageId: req.MessageId}
		handler(req, ret)

		buf, err := proto.Marshal(ret)
		if err != nil {
			c.conn.Publish(msg.Reply, errRpcMarshal)
			c.conn.Flush()
			return
		}

		c.conn.Publish(msg.Reply, buf)
		c.conn.Flush()
	})
}

func (c *natsQueueConn) Publish(subject string, msg *message.Message) error {
	if msg == nil {
		return errors.New("nil message to publish")
	}

	msg.MessageId = util.GenMsgID()

	buf, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	c.conn.Publish(subject, buf)
	c.conn.Flush()
	return c.conn.LastError()
}

func (c *natsQueueConn) Request(subject string, req *message.RpcMessage, timeout time.Duration) (*message.RpcMessage, error) {
	if req == nil {
		return nil, errors.New("nil message to request")
	}
	req.MessageId = util.GenMsgID()

	buf, err := proto.Marshal(req)
	if err != nil {
		return nil, err
	}
	msg, err := c.conn.Request(subject, buf, timeout)
	if err != nil {
		return nil, err
	}
	ret := &message.RpcMessage{}
	if err = proto.Unmarshal(msg.Data, ret); err != nil {
		return nil, err
	}
	return ret, nil
}

func (c *natsQueueConn) Enqueue(subject string, msg *message.Message) error {
	return nil
}

func (c *natsQueueConn) Dequeue(subject, group string, timeout time.Duration, dst proto.Message) (*message.Meta, error) {
	return nil,nil
}

func (c *natsQueueConn) Close() error {
	c.conn.Close()
	return nil
}
