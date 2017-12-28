## Description
Queue Lib for ATC, internal support Redis, Nats, protobuffer etc. 

## Install
  `go get github.com/adolphlxm/atc/queue`
   
## Usage

#### Publisher
```go
import (
    "github.com/adolphlxm/atc/queue/message"
    
    "github.com/adolphlxm/atc/queue"
    _ "github.com/adolphlxm/atc/queue/queue_redis"
)

pub, _ := queue.NewPublisher("redis", "redis://127.0.0.1:6379")
pub.Publish("subject", &message.Message{
    Body: util.MustMessageBody(nil, /* point to your protobuffer struct */ ),
})
pub.Close()
```

#### Consumer
```go
import (
    "github.com/adolphlxm/atc/queue/message"
    
    "github.com/adolphlxm/atc/queue"
    _ "github.com/adolphlxm/atc/queue/queue_redis"
)

con, _ := queue.NewConsumer("redis", "redis://127.0.0.1:6379")
sub,_ := con.Subscribe("subject", "cluster-group")
msg, _ := sub.NextMessage(time.Second)
// logic for msg
con.Close()
```
