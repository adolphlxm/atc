package queue_nats_streaming_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/adolphlxm/atc/queue"
	"github.com/adolphlxm/atc/queue/message"
	_ "github.com/adolphlxm/atc/queue/queue_nats_streaming"
	"github.com/adolphlxm/atc/queue/testdata"
	"github.com/adolphlxm/atc/queue/util"
)

const (
	driverName  = "nats-streaming"
	natsPubDSN  = "nats://127.0.0.1:4222?cluster=test-cluster&client=xxx_pub"
	natsSubDSN  = "nats://127.0.0.1:4222?cluster=test-cluster&client=xxx_sub"
	testSubject = "testSubject"
)

func TestNatsStreamingQueueDriver_Open(t *testing.T) {
	qc, err := queue.NewPublisher(driverName, natsPubDSN)
	if err != nil {
		t.Fatal(err)
	}
	qc.Close()
}

func TestNatsStreamingQueueConn_Dequeue(t *testing.T) {
	want := testdata.Something{Age: int32(rand.Uint32())}
	msg := &message.Message{Body: util.MustMessageBody(&want)}
	done := make(chan struct{})

	// sub
	go func() {
		defer func() {
			done <- struct{}{}
		}()

		de, err := queue.NewConsumer(driverName, natsSubDSN)
		if err != nil {
			t.Fatal(err)
		}
		defer de.Close()

		got := &testdata.Something{}
		meta, err := de.Dequeue(testSubject, "test", 10*time.Second, got)
		if err != nil {
			t.Fatal(err)
		}
		if msg.MessageId != meta.MessageId {
			t.Errorf("dequeue meta's message id want %#x, got %#x", msg.MessageId, meta.MessageId)
		}
		if got.Age != want.Age {
			t.Errorf("dequeue data's age want %d, got %d", want.Age, got.Age)
		}
	}()

	// pub
	qc, err := queue.NewPublisher(driverName, natsPubDSN)
	if err != nil {
		t.Fatal(err)
	}
	defer qc.Close()

	if err = qc.Enqueue(testSubject, msg); err != nil {
		t.Error(err)
	}
	<-done
}
