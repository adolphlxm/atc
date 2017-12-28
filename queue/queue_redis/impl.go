package queue_redis

import (
	"errors"
	"time"

	"github.com/adolphlxm/atc/queue"
	"github.com/adolphlxm/atc/queue/message"

	"github.com/garyburd/redigo/redis"
	"github.com/golang/protobuf/proto"
    "github.com/adolphlxm/atc/queue/util"
)

var (
	ErrNotIMPL = errors.New("queue redis: Not implement")
	ErrSubFail = errors.New("queue redis: subscribe fail")
)

type redisSubscriber struct {
	subject string
	pbConn  redis.PubSubConn
}

func (s *redisSubscriber) NextMessage(timeout time.Duration) (*message.Message, error) {

	for {
		switch n := s.pbConn.Receive().(type) {
		case redis.Message:
			ret := &message.Message{}
			err := proto.Unmarshal(n.Data, ret)
			return ret, err
		case redis.Subscription:
			if n.Count == 0 {
				return nil, ErrSubFail
			}
		case error:
			return nil, n
		}
	}

	return nil, nil
}

type redisQueueConn struct {
	cs *redis.Pool
}

func (d *redisQueueConn) peekAvailableConn() (c redis.Conn, err error) {

	limit := 3
	for {
		c = d.cs.Get()
		if err = c.Err(); err != nil {
			c.Close()
			limit--
			if limit == 0 {
				return nil, err
			}
			continue
		}
		return
	}
}

func (d *redisQueueConn) Publish(subject string, msg *message.Message) error {
	msg.MessageId = util.GenMsgID()
	c, err := d.peekAvailableConn()
	if err != nil {
		return err
	}

	buf, err := proto.Marshal(msg)
	if err != nil {
		return err
	}
	c.Send("PUBLISH", subject, buf)
	c.Flush()
	return c.Err()
}

func (d *redisQueueConn) Subscribe(subject, group string) (queue.Subscriber, error) {
	c, err := d.peekAvailableConn()
	if err != nil {
		return nil, err
	}

	psC := redis.PubSubConn{Conn: c}
	if err = psC.Subscribe(subject); err != nil {
		return nil, err
	}
	return &redisSubscriber{pbConn: psC, subject: subject}, err
}

func (d *redisQueueConn) Request(subject string, req *message.RpcMessage, timeout time.Duration) (*message.RpcMessage, error) {
	return nil, ErrNotIMPL
}

func (d *redisQueueConn) RpcHandle(subject, group string, handler queue.RpcHandler) {
}

func (d *redisQueueConn) Close() error {
	return d.cs.Close()
}
