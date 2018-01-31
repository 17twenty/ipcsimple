package ipcsimple

import "github.com/gorilla/mux"

type RESTBroker struct {
	muxer *mux.Router
}

func NewRESTBroker(port int) (*RESTBroker, error) {
	b := &RESTBroker{}
	b.muxer = mux.NewRouter()
	return nil, nil
}

type RESTProducer struct {
	endpoints []string
}

func (ipc *RESTBroker) NewProducer(brokers []string) *RESTProducer {
	return &RESTProducer{
		endpoints: brokers,
	}
}

type RESTConsumer struct {
	endpoints []string
}

type ConsumerMessage struct {
	Key, Value []byte
	Topic      string
}

func (c *RESTConsumer) Errors() <-chan *error {
	return make(<-chan *error)
}

func (c *RESTConsumer) Messages() <-chan *ConsumerMessage {
	return make(<-chan *ConsumerMessage)
}

func (c *RESTConsumer) Close() error {
	return nil
}

func (ipc *RESTBroker) NewConsumer(brokers []string) *RESTConsumer {
	return &RESTConsumer{
		endpoints: brokers,
	}
}

func (ipc *RESTBroker) SubscribeToTopic(topic string) (chan<- []byte, error) {
	return nil, nil
}

func (ipc *RESTBroker) StartServer() {

}

func (ipc *RESTBroker) Unsubscribe(topic string) {

}
