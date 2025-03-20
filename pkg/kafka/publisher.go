package kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"

	"github.com/ashabykov/geospatial_cache_for_meetup/location"
)

type rateLimiter struct {
	ch *time.Ticker
}

func newRateLimiter(rps int) *rateLimiter {
	interval := time.Second / time.Duration(rps)
	return &rateLimiter{
		ch: time.NewTicker(interval),
	}
}

func (rl *rateLimiter) Block() {
	<-rl.ch.C
}

type Publisher struct {
	w  *kafka.Writer
	rl *rateLimiter
}

func NewPublisher(hosts []string, topic string, rps int) *Publisher {

	writer := &kafka.Writer{
		Addr:                   kafka.TCP(hosts...),
		Topic:                  topic,
		AllowAutoTopicCreation: true,
		Balancer:               &kafka.LeastBytes{},
	}

	return &Publisher{
		w:  writer,
		rl: newRateLimiter(rps),
	}
}

func (p *Publisher) Publish(ctx context.Context, locations ...location.Location) error {
	for i := range locations {

		p.rl.Block()

		key := []byte(locations[i].Key())

		data, err := location.Encode(locations[i])

		if err != nil {

			fmt.Println("encoding location error:", err)

			continue
		}

		msg := kafka.Message{
			Key:   key,
			Value: data,
		}

		if err = p.w.WriteMessages(ctx, msg); err != nil {

			fmt.Println("write messages error:", err)

			return err

		}

		fmt.Println("write messages ok")
	}

	return nil
}

func (p *Publisher) Close() error {
	return p.w.Close()
}
