package main

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/17twenty/ipcsimple"
)

type DemoStruct struct {
	Hello string `json:"hello`
	Count int    `json:"count"`
}

func main() {

	temp := DemoStruct{
		Hello: "Heya",
		Count: 5,
	}

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	tempAsBytes, err := json.Marshal(temp)
	if err != nil {
		fmt.Println("error:", err)
	}

	var totalChecked int

	consumerA := ipcsimple.NewConsumer(9090)

	producer := ipcsimple.NewProducer([]string{"http://localhost:9090/ipcsimple"})

	go func() {
		for i := 0; i < 100; i++ {
			err = producer.SendMessage(&ipcsimple.ProducerMessage{
				Topic: "testA",
				Value: ipcsimple.ByteEncoder(tempAsBytes),
			})
			if err != nil {
				log.Println("About to panic()", err)
				time.Sleep(time.Second * 20)
				panic(err)
			}
			err = producer.SendMessage(&ipcsimple.ProducerMessage{
				Topic: "testB",
				Value: ipcsimple.ByteEncoder(tempAsBytes),
			})
			if err != nil {
				log.Println("About to panic()", err)
				time.Sleep(time.Second * 20)
				panic(err)
			}
			time.Sleep(time.Millisecond * 500)
		}
	}()
	log.Println("Closing...")
	if err := producer.Close(); err != nil {
		log.Println("About to panic()", err)
		time.Sleep(time.Second * 20)
		panic(err)
	}

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

	fmt.Printf("producer and consumer worked! %d messages ok\n", totalChecked)
}
