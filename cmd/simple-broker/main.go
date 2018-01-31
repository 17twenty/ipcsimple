package main

import (
	"fmt"
	"log"

	ipcsimple "github.com/17twenty/ipc-simple"
	"github.com/Shopify/sarama"
)

type check struct {
	partition int32
	offset    int64
	message   string
}

const (
	topic         = "test_topic"
	messageCount  = 15
	clientID      = "test_client"
	numPartitions = int32(8)
	brokerAddr    = "127.0.0.1:9092"
	raftAddr      = "127.0.0.1:9093"
	serfAddr      = "127.0.0.1:9094"
	httpAddr      = "127.0.0.1:9095"
	logDir        = "logdir"
	brokerID      = 0
)

func main() {
	broker, err := ipcsimple.NewRESTBroker(8080)

	brokers := []string{brokerAddr}
	producer := broker.NewProducer(brokers)

	pmap := make(map[int32][]check)

	for i := 0; i < messageCount; i++ {
		message := fmt.Sprintf("Hello from Jocko #%d!", i)
		partition, offset, err := producer.SendMessage(&ipcsimple.ProducerMessage{
			Topic: topic,
			Value: sarama.StringEncoder(message),
		})
		if err != nil {
			panic(err)
		}
		pmap[partition] = append(pmap[partition], check{
			partition: partition,
			offset:    offset,
			message:   message,
		})
	}
	if err = producer.Close(); err != nil {
		panic(err)
	}

	var totalChecked int
	for partitionID := range pmap {
		checked := 0
		consumer := broker.NewConsumer(brokers)

		for msg := range consumer.Messages() {
			if string(msg.Value) != check.message {
				log.Fatalln("msg values not equal")
			}
			// if msg.Offset != check.offset {
			// 	logger.Fatal("msg offsets not equal", log.Int32("partition", msg.Partition), log.Int64("offset", msg.Offset))
			// }
			log.Println(fmt.Sprintf("%v", string(msg.Value)))
			log.Println("msg is ok")
		}
	}
	fmt.Printf("producer and consumer worked! %d messages ok\n", totalChecked)
}
