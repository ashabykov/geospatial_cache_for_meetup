package kafka_broadcast

import (
	"context"

	"github.com/segmentio/kafka-go"

	"github.com/ashabykov/geospatial_cache_for_meetup"
)

type Publisher struct {
	w *kafka.Writer
}

func NewPublisher(hosts []string, topic string) *Publisher {
	return &Publisher{
		w: &kafka.Writer{
			Addr:                   kafka.TCP(hosts...),
			Topic:                  topic,
			AllowAutoTopicCreation: true,
			Balancer:               &kafka.LeastBytes{},
		},
	}
}

func (p *Publisher) Publish(ctx context.Context, message ...geospatial_cache_for_meetup.Message) error {
	msgs := make([]kafka.Message, 0, len(message))
	for i := range message {
		msgs = append(msgs, kafka.Message{
			Key:   message[i].Key,
			Value: message[i].Val,
		})
	}
	return p.w.WriteMessages(ctx, msgs...)
}

func (p *Publisher) Close() error {
	return p.w.Close()
}
