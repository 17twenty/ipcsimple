# Simple IPC

This is a simple IPC used to move development forward when you can't be arsed with SQS or Kafka.

## Getting Started

### Producer

```golang
// The brokers are all the endpoints we want to post to
producer := ipcsimple.NewProducer([]string{"http://localhost:9090/ipcsimple"})
...
err = producer.SendMessage(&ipcsimple.ProducerMessage{
    Topic: "testA",
    Value: ipcsimple.ByteEncoder(tempAsBytes),
})
```

### Consumer

```golang
consumerA := ipcsimple.NewConsumer(9090) // local port to bind too
...
log.Println("Reading Messages...")
for msg := range consumerA.Messages() {
    switch msg.Topic {
    case "testA":
        var data DemoStruct
        log.Println(fmt.Sprintf("%v", string(msg.Value)))
        err := json.Unmarshal([]byte(msg.Value), &data)
        log.Println(err, fmt.Sprintf("%#+v", data))
        log.Println("msg is ok")

    case "testB":
        log.Println("Got RAW testB message", msg.Value)
    default:
        log.Println("ERR: Unknown topic", string(msg.Topic))
    }
}
```

## Issues

* Server timeouts could be a problem
* Doesn't fail that gracefully
* Close doesn't do much
* Could take a logger as an input rather than stdoutting it
* Doesn't ever exit inner go routine
