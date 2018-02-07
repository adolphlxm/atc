package atc

import (
	"github.com/adolphlxm/atc/queue"
	"github.com/adolphlxm/atc/logs"
)


var queuePublisher map[string]QueuePblisherShutDown
var queueConsumer map[string]QueueConsumerShutDown

type QueuePblisherShutDown struct {
	publisher queue.Publisher
	Module string
}
func (this *QueuePblisherShutDown) ModuleID() string {
	return this.Module
}
func (this *QueuePblisherShutDown) Stop() error{
	logs.Tracef("grace: queue.publisher [%s] stop.", this.ModuleID())
	return this.publisher.Close()
}

type QueueConsumerShutDown struct {
	consumer queue.Consumer
	Module string
}
func (this *QueueConsumerShutDown) ModuleID() string {
	return this.Module
}
func (this *QueueConsumerShutDown) Stop() error{
	logs.Tracef("grace:queue.consumer [%s] stop.", this.ModuleID())
	return this.consumer.Close()
}

func RunQueuePublisher() {
	queuePublisher = make(map[string]QueuePblisherShutDown, 0)
	aliasnames := AppConfig.Strings("queue.publisher.aliasnames")
	for _, aliasname := range aliasnames {
		keyPerfix := "queue.publisher." + aliasname + "."
		addrs := AppConfig.String(keyPerfix + "addrs")
		drivername := AppConfig.String(keyPerfix + "driver")

		logs.Tracef("queue.publisher:[%s] starting....", aliasname)

		publisher, err := queue.NewPublisher(drivername, addrs)
		if err != nil {
			logs.Errorf("queue.publisher:[%s] start fail err:%s", aliasname, err.Error())
			panic(err)
		}

		shutDown := QueuePblisherShutDown{publisher:publisher,Module:aliasname}
		queuePublisher[aliasname] = shutDown
		GracePushFront(&shutDown)
		logs.Tracef("queue.publisher:[%s] Running on %s.", aliasname, addrs)
	}
}

func RunQueueConsumer() {
	queueConsumer = make(map[string]QueueConsumerShutDown, 0)
	aliasnames := AppConfig.Strings("queue.consumer.aliasnames")

	for _, aliasname := range aliasnames {
		keyPerfix := "queue.consumer." + aliasname + "."
		addrs := AppConfig.String(keyPerfix + "addrs")
		drivername := AppConfig.String(keyPerfix + "driver")

		logs.Tracef("queue.consumer:[%s] starting....", aliasname)

		consumer, err := queue.NewConsumer(drivername, addrs)
		if err != nil {
			logs.Errorf("queue.consumer:[%s] start fail err:%s", aliasname, err.Error())
			panic(err)
		}

		shutDown := QueueConsumerShutDown{consumer:consumer,Module:aliasname}
		queueConsumer[aliasname] = shutDown
		GracePushFront(&shutDown)
		logs.Tracef("queue.consumer:[%s] Running on %s.", aliasname, addrs)
	}
}

func GetPublisher(aliasname string) queue.Publisher {
	return queuePublisher[aliasname].publisher
}

func GetConsumer(aliasname string) queue.Consumer {
	return queueConsumer[aliasname].consumer
}