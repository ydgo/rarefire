package broker

type Broker interface {
	Connect() error
	Close() error
	Publish(topic string, m *Message) error
	Subscribe(topic string, h Handler) (Subscriber, error)
}

// Message 消息队列中保存的消息数据
type Message struct {
	Body []byte
}

// Handler 处理消息的方法
type Handler func(Event) error

// Event 将消息发布的整个过程视为一个事件
type Event interface {
	Topic() string
	Message() *Message
}

// Subscriber 订阅者
type Subscriber interface {
	Topic() string
	Unsubscribe() error
}
