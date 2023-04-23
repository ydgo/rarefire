package broker

import (
	"errors"
	"github.com/go-redis/redis"
	"sync"
)

type publication struct {
	topic   string
	message *Message
	err     error
}

func (p *publication) Topic() string {
	return p.topic
}

func (p *publication) Message() *Message {
	return p.message
}

type subscriber struct {
	pubSub *redis.PubSub
	topic  string
	handle Handler
}

func (s *subscriber) recv() {
	defer s.pubSub.Close()
	for msg := range s.pubSub.Channel() {
		var m Message
		m.Body = []byte(msg.Payload)
		p := publication{
			topic:   msg.Channel,
			message: &m,
		}
		if p.err = s.handle(&p); p.err != nil {
			break
		}
		// handle error?
	}
}

// Unsubscribe unsubscribes the subscriber and frees the connection.
func (s *subscriber) Unsubscribe() error {
	return s.pubSub.Unsubscribe(s.topic)
}

// Topic returns the topic of the subscriber.
func (s *subscriber) Topic() string {
	return s.topic
}

type redisBroker struct {
	sync.Mutex
	addr      string
	connected bool
	c         *redis.Client
}

func (r *redisBroker) Connect() error {
	r.Lock()
	defer r.Unlock()
	if r.connected {
		return errors.New("already connected")
	}
	r.c = redis.NewClient(&redis.Options{
		Addr:         r.addr,
		MinIdleConns: 5,
		DB:           0,
	})
	return r.c.Ping().Err()
}

func (r *redisBroker) Close() error {
	r.Lock()
	defer r.Unlock()
	r.connected = false
	return r.c.Close()
}

func (r *redisBroker) Publish(topic string, m *Message) error {
	return r.c.Publish(topic, string(m.Body)).Err()
}

func (r *redisBroker) Subscribe(topic string, h Handler) (Subscriber, error) {
	s := subscriber{
		pubSub: r.c.Subscribe(topic),
		topic:  topic,
		handle: h,
	}
	go s.recv()
	return &s, nil
}

func NewRedisBroker(addr string) Broker {
	return &redisBroker{
		addr:      addr,
		connected: false,
	}
}
