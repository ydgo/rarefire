package broker

import "testing"

var (
	addr  = "127.0.0.1:6379"
	topic = "rarefire:braoke:redis:test"
	value = "hello yudong"
)

func TestNewRedisBroker(t *testing.T) {
	broker := NewRedisBroker(addr)
	if err := broker.Connect(); err != nil {
		t.Fatal("connect fail: ", err)
	}
}

func TestListElement_Message(t *testing.T) {
	broker := NewRedisBroker(addr)
	if err := broker.Connect(); err != nil {
		t.Fatal("connect fail: ", err)
	}

	if err := broker.Publish(topic, &Message{body: []byte(value)}); err != nil {
		t.Fatal("publish fail: ", err)
	}

	if err := broker.Subscribe(topic, func(event Event) error {
		data := NewListElement(event)
		if string(data.Message().body) != value {
			t.Fatal("data in and out mismatch")
		}
		return nil
	}); err != nil {
		t.Fatal("subscribe fail: ", err)
	}

}
