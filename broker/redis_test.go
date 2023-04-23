package broker

import (
	"fmt"
	"reflect"
	"sort"
	"testing"
	"time"
)

var (
	addr  = "127.0.0.1:6379"
	topic = "rarefire:broker:test1"
)

func TestNewRedisBroker(t *testing.T) {
	broker := NewRedisBroker(addr)
	if err := broker.Connect(); err != nil {
		t.Fatal("connect fail: ", err)
	}
}

func subscribe(t *testing.T, b Broker, topic string, handle Handler) Subscriber {
	s, err := b.Subscribe(topic, handle)
	if err != nil {
		t.Error(err)
	}
	return s
}

func publish(t *testing.T, b Broker, topic string, msg *Message) {
	if err := b.Publish(topic, msg); err != nil {
		t.Error(err)
	}
}

func unsubscribe(t *testing.T, s Subscriber) {
	if err := s.Unsubscribe(); err != nil {
		t.Error(err)
	}
}

func TestBroker(t *testing.T) {
	broker := NewRedisBroker(addr)
	if err := broker.Connect(); err != nil {
		t.Fatal("connect fail: ", err)
	}
	msgs := make(chan string, 10)
	go func() {
		s1 := subscribe(t, broker, topic, func(p Event) error {
			m := p.Message()
			msgs <- fmt.Sprintf("s1:%s", string(m.Body))
			return nil
		})
		s2 := subscribe(t, broker, topic, func(p Event) error {
			m := p.Message()
			msgs <- fmt.Sprintf("s2:%s", string(m.Body))
			return nil
		})

		publish(t, broker, topic, &Message{Body: []byte("hello")})
		publish(t, broker, topic, &Message{Body: []byte("world")})

		unsubscribe(t, s1)
		time.Sleep(time.Second)

		publish(t, broker, topic, &Message{Body: []byte("other")})

		unsubscribe(t, s2)
		time.Sleep(time.Second)

		publish(t, broker, topic, &Message{Body: []byte("none")})

		time.Sleep(time.Second)
		close(msgs)
	}()

	var actual []string
	for msg := range msgs {
		actual = append(actual, msg)
	}
	exp := []string{
		"s1:hello",
		"s2:hello",
		"s1:world",
		"s2:world",
		"s2:other",
	}

	// Order is not guaranteed.
	sort.Strings(actual)
	sort.Strings(exp)
	if !reflect.DeepEqual(actual, exp) {
		t.Fatalf("expected %v, got %v", exp, actual)
	}

}
