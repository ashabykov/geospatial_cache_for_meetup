package group

import (
	"context"
	"log"
	"sync"

	"github.com/segmentio/kafka-go"

	"github.com/ashabykov/geospatial_cache_for_meetup/location"
)

const (
	defaultResultsBufferSize = 1000
)

type Subscriber struct {
	name   string
	reader *kafka.Reader
}

func NewSubscriber(
	name string,
	hosts []string,
	topic string, groupID string,
) *Subscriber {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:  hosts,
		GroupID:  groupID,
		Topic:    topic,
		MaxBytes: int(10e6), // 10MB
	})

	return &Subscriber{
		name:   name,
		reader: reader,
	}
}

func (s *Subscriber) Subscribe(ctx context.Context) (<-chan location.Location, error) {

	results := make(chan location.Location, defaultResultsBufferSize)

	wg := &sync.WaitGroup{}
	wg.Add(1)

	go s.worker(ctx, results, wg)

	go func() {

		wg.Wait()

		close(results)
	}()

	return results, nil
}

func (s *Subscriber) worker(ctx context.Context, results chan<- location.Location, wg *sync.WaitGroup) {

	defer wg.Done()

	for {
		select {
		default:
			message, err := s.reader.ReadMessage(ctx)
			if err != nil {
				log.Println("Subscriber: read message error:", err)
				continue
			}

			loc, err := location.Decode(message.Value)
			if err != nil {
				log.Println("Subscriber: decode message error:", err)
				continue
			}
			results <- loc

		case <-ctx.Done():
			return
		}
	}
}
