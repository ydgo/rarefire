package consumer

import (
	"context"
	"github.com/redis/go-redis/v9"
	"log"
	"strings"
	"time"
)

type XMessage struct {
	ID     string
	Values map[string]interface{}
}

type Consumer interface {
	Consume(func(m XMessage)) error
	Close() error
}

type Conf struct {
	Redis     string
	Stream    string
	Group     string
	Consumers []string
}

type groupConsumer struct {
	conf   *Conf
	client *redis.Client
	done   chan struct{}
}

func (c *Conf) New() Consumer {
	client := redis.NewClient(&redis.Options{Addr: c.Redis})
	return &groupConsumer{
		conf:   c,
		client: client,
		done:   make(chan struct{}),
	}
}

func (g *groupConsumer) Consume(handle func(m XMessage)) error {
	log.Printf("stream: %s, group: %s, consumer: %s start consuming...\n", g.conf.Stream, g.conf.Group, g.conf.Consumers)
	ctx := context.Background()
	args := redis.XReadGroupArgs{
		Group:    g.conf.Group,
		Consumer: g.conf.Consumers[0],
		Streams:  []string{g.conf.Stream, ">"},
		Block:    time.Second * 3,
	}
	err := g.client.XGroupCreateMkStream(ctx, g.conf.Stream, g.conf.Group, "0").Err()
	if err != nil {
		if !strings.Contains(err.Error(), "BUSYGROUP") {
			log.Println("xgroup create err: ", err.Error())
			return err
		}
	}
	for {
		select {
		case <-g.done:
			return nil
		default:
			xStreams, err := g.client.XReadGroup(ctx, &args).Result()
			if err != nil {
				if err == redis.Nil {
					continue
				}
				log.Println("xread group err: ", err)
				continue
			}
			for _, xStream := range xStreams {
				for _, msg := range xStream.Messages {
					xMessage := XMessage{ID: msg.ID, Values: msg.Values}
					log.Printf("handle message: %+v\n", xMessage)
					handle(xMessage)
					err = g.client.XAck(ctx, g.conf.Stream, g.conf.Group, msg.ID).Err()
					if err != nil {
						log.Printf("xack: %s,  err: %s\n", msg.ID, err.Error())
					}
				}
			}

		}
	}
}

func (g *groupConsumer) Close() error {
	log.Println("consumer close.")
	g.done <- struct{}{}
	return g.client.Close()
}
