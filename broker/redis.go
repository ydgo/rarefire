package broker

import (
	"errors"
	"github.com/go-redis/redis"
	"sync"
)

// 目前仅支持list类型的topic

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
	return r.c.LPush(topic, string(m.body)).Err()
}

func (r *redisBroker) Subscribe(topic string, h Handler) error {
	v, err := r.c.RPop(topic).Result()
	if err != nil {
		return err
	}
	return h(&listElement{Key: topic, Body: []byte(v)})
}

func NewRedisBroker(addr string) Broker {
	return &redisBroker{
		addr:      addr,
		connected: false,
	}
}

// 实现一个列表类型的消息队列进行订阅和发布
type listElement struct {
	Key  string
	Body []byte
}

func (l *listElement) Topic() string {
	return l.Key
}

func (l *listElement) Message() *Message {
	return &Message{body: l.Body}
}

func NewListElement(event Event) Event {
	return &listElement{
		Key:  event.Topic(),
		Body: event.Message().body,
	}
}
