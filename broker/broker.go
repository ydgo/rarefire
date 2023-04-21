package broker

type Broker interface {
	Connect() error
	Close() error
	Publish(topic string, m *Message) error
	Subscribe(topic string, h Handler) error
}

type Message struct {
	body []byte
}

// message 是生产者放进去的？需不需要也封装一层？

// Event 是从消息队列拿出来后封装的？

type Handler func(Event) error

type Event interface {
	Topic() string
	Message() *Message
}
