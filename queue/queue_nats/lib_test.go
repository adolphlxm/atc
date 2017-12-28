package queue_nats_test

import (
	"testing"
	"time"

	"github.com/adolphlxm/atc/queue"
	"github.com/adolphlxm/atc/queue/message"
	_ "github.com/adolphlxm/atc/queue/queue_nats"
	"github.com/adolphlxm/atc/queue/testdata"
	"github.com/adolphlxm/atc/queue/util"

	"github.com/golang/protobuf/ptypes"
	"github.com/nats-io/go-nats"
)

const (
	driverName  = "nats"
	natsDSN     = "nats://127.0.0.1:4222"
	testSubject = "testSubject"
	rpcSubject  = "rpcSubject"
)

func newSubscriber() (queue.Subscriber, error) {
	qc, err := queue.NewConsumer(driverName, natsDSN)
	if err != nil {
		return nil, err
	}
	return qc.Subscribe(testSubject, "test")
}

func TestNewPublisher(t *testing.T) {
	want := &testdata.Something{
		Name: "something",
		Age:  11,
	}

	msg := &message.Message{
		Priority: message.MsgPriority_PRIORITY0,
		Body:     util.MustMessageBody(want),
	}

	ch := make(chan struct{})
	go func() {
		sub, err := newSubscriber()
		ch <- struct{}{}
		if err != nil {
			t.Fatal(err)
		}

		m, err := sub.NextMessage(10 * time.Second)
		if err != nil {
			t.Error(err)
		}
		got := testdata.Something{}
		if err := ptypes.UnmarshalAny(m.Body, &got); err != nil {
			t.Error(err)
		}
		if msg.MessageId != m.MessageId {
			t.Errorf("message id: want %#x, got %#x", msg.MessageId, m.MessageId)
		}

		if want.Name != got.Name {
			t.Errorf("name: want %v, got %v", want.Name, got.Name)
		}

		if want.Age != got.Age {
			t.Errorf("Age: want %v, got %v", want.Age, got.Age)
		}

		ch <- struct{}{}
	}()

	<-ch
	qc, _ := queue.NewPublisher(driverName, natsDSN)
	if err := qc.Publish(testSubject, msg); err != nil {
		t.Error(err)
	}
	qc.Close()
	<-ch
}

func TestNewPublisherWithBadDSN(t *testing.T) {
	_, err := queue.NewPublisher(driverName, "nats://xxx")
	if err == nil {
		t.Error("want error, got nil")
	}
}

func TestDrivers(t *testing.T) {
	ds := queue.Drivers()
	if len(ds) != 1 {
		t.Errorf("driver want 1, got %d", len(ds))
	}
	if ds[0] != "nats" {
		t.Errorf("driver's name want nats, got %s", ds[0])
	}
}

func TestNewPublisherWithOptions(t *testing.T) {
	c, err := queue.NewPublisher(driverName, nats.DefaultURL+"/?maxReconnect=1")
	c.Close()
	if err != nil {
		t.Error(err)
	}

	_, err = queue.NewPublisher(driverName, nats.DefaultURL+"/?maxReconnect=true")
	if err == nil {
		t.Error("want error, got nil")
	}

	_, err = queue.NewPublisher(driverName, nats.DefaultURL+"/?connectTimeout=true")
	if err == nil {
		t.Error("want error, got nil")
	}

	c, err = queue.NewPublisher(driverName, nats.DefaultURL+"/?connectTimeout=1000")
	c.Close()
	if err != nil {
		t.Error(err)
	}

	c, err = queue.NewPublisher(driverName, nats.DefaultURL+"/?noRandomize=true")
	c.Close()
	if err != nil {
		t.Error(err)
	}

	c, err = queue.NewPublisher(driverName, nats.DefaultURL+"/?allowReconnect=false")
	c.Close()
	if err != nil {
		t.Error(err)
	}

	_, err = queue.NewPublisher(driverName, nats.DefaultURL+"/?notexists=xxx")
	if err == nil {
		t.Error("want error, got nil")
	}

	c, err = queue.NewPublisher(driverName, nats.DefaultURL+"/?name=test")
	c.Close()
	if err != nil {
		t.Error(err)
	}

	_, err = queue.NewPublisher(driverName, nats.DefaultURL+"/?name=")
	if err == nil {
		t.Error("want error, got nil")
	}
}

func TestNatsQueueConn_Publish(t *testing.T) {
	c, _ := queue.NewPublisher(driverName, nats.DefaultURL+"/?name=test")
	defer c.Close()
	err := c.Publish(testSubject, nil)
	if err == nil {
		t.Error("publish nil message, got non error")
	}
}

func TestNatsQueueConn_Request(t *testing.T) {
	c, _ := queue.NewPublisher(driverName, nats.DefaultURL+"/?name=test")
	defer c.Close()
	_, err := c.Request(testSubject, &message.RpcMessage{}, time.Millisecond)
	if err == nil {
		t.Error("request nil message, got non error")
	}

	_, err = c.Request(testSubject, &message.RpcMessage{}, time.Microsecond)
	if err == nil {
		t.Error("request message with short timeout, got non error")
	}

	_, err = c.Request(testSubject, nil, time.Microsecond)
	if err == nil {
		t.Error("request nil message, got non error")
	}
}

func TestNatsQueueConn_RPCHandler(t *testing.T) {
	c, _ := queue.NewConsumer(driverName, nats.DefaultURL+"/?name=rpcServer")

	req := &testdata.Something{Name: "request data"}
	wrappedReq := util.MustMessageBody(req)
	reqMsg := &message.RpcMessage{Body: wrappedReq}

	rsp := &testdata.Something{Name: "reply to aaa", Age: 10}
	wrappedRsq := util.MustMessageBody(rsp)

	c.RpcHandle(rpcSubject, "test", func(request *message.RpcMessage, response *message.RpcMessage) {
		response.Code = 0
		response.Body = &(*wrappedRsq)
	})

	// create a rpc request
	pc, _ := queue.NewPublisher(driverName, nats.DefaultURL+"/?name=rpcClient")
	got, err := pc.Request(rpcSubject, reqMsg, time.Second)
	if err != nil {
		t.Fatalf("Publisher Request with error: %v", err)
	}

	if got.Code != 0 {
		t.Fatalf("Publisher Request with error code: %d", got.Code)
	}

	if got.MessageId != reqMsg.MessageId {
		t.Errorf("Publisher Request id want %#x, got %#x", reqMsg.MessageId, got.MessageId)
	}
	gotData := testdata.Something{}
	if err = ptypes.UnmarshalAny(got.Body, &gotData); err != nil {
		t.Fatalf("Publisher Request Unmarshal Body err: %v", err)
	}

	if gotData.Name != rsp.Name {
		t.Errorf("Publisher Request data Name field want %s, got %s", rsp.Name, gotData.Name)
	}

	if gotData.Age != rsp.Age {
		t.Errorf("Publisher Request data Age field want %s, got %s", rsp.Age, gotData.Age)
	}
}
