package ipcsimple

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

var (
	ErrBadResponse = errors.New("ipcsimple: Response was not OK")
)

// RESTProducer and below defines the producer
type RESTProducer struct {
	endpoints []string
}

func NewProducer(brokers []string) *RESTProducer {
	return &RESTProducer{
		endpoints: brokers,
	}
}

type ProducerMessage struct {
	Topic string
	Value Encoder
}

func (p *RESTProducer) SendMessage(msg *ProducerMessage) error {

	b, err := json.Marshal(*msg)
	if err != nil {
		fmt.Println("error:", err)
	}

	// We iterate through all the 'subscribers' and post to their
	// <brokerEndpoint>/send/<topic>
	for _, url := range p.endpoints {
		req, err := http.NewRequest("POST", url, bytes.NewReader(b))
		req.Header.Set("User-Agent", "ipcsimple")
		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			continue
		}
		defer resp.Body.Close()

		body, _ := ioutil.ReadAll(resp.Body)

		output := struct {
			Status string `json:"status"`
		}{}
		err = json.Unmarshal(body, &output)
		if err != nil {
			log.Println("Unmarshal:", err)
			return ErrBadResponse
		}
		if output.Status != "ok" {
			return ErrBadResponse
		}
	}
	return nil
}

func (p *RESTProducer) Close() error {
	return nil
}

// RESTConsumer and below defines the consumer
type RESTConsumer struct {
	muxer       *mux.Router
	messageChan chan *ConsumerMessage
	errorChan   chan *error
}

type ConsumerMessage struct {
	Value []byte
	Topic string
}

func (c *RESTConsumer) Errors() <-chan *error {
	return c.errorChan
}

func (c *RESTConsumer) Messages() <-chan *ConsumerMessage {
	return c.messageChan
}

func (c *RESTConsumer) Close() error {
	return nil
}

func NewConsumer(port int) *RESTConsumer {
	b := &RESTConsumer{
		messageChan: make(chan *ConsumerMessage, 3),
		errorChan:   make(chan *error, 3),
	}
	b.muxer = mux.NewRouter()
	b.muxer.HandleFunc("/ipcsimple", b.consumerHandler).Methods("POST")
	go func(port int) {
		http.Handle("/", b.muxer)
		log.Println("Serving on", fmt.Sprintf(":%d", port))
		http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	}(port)

	// Give the consumer time to start listening...
	time.Sleep(1 * time.Second)

	return b
}

func (c *RESTConsumer) consumerHandler(wr http.ResponseWriter, req *http.Request) {
	var m ConsumerMessage
	b, _ := ioutil.ReadAll(req.Body)
	err := json.Unmarshal(b, &m)
	if err != nil {
		log.Println("Unmarshal:", err)
	}

	wr.WriteHeader(http.StatusOK)
	wr.Header().Set("Content-Type", "application/json")
	json.NewEncoder(wr).Encode(struct {
		Status string `json:"status"`
	}{
		Status: "ok",
	})

	c.messageChan <- &m
}
