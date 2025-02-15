package kafka_broadcast

import (
	"context"
	"time"

	"github.com/segmentio/kafka-go"

	"github.com/ashabykov/geospatial_cache_for_meetup"
)

type Subscriber struct {
	partitions map[int]*kafka.Reader

	name      string
	timeRange time.Duration
	timout    time.Duration
}

func NewSubscriber(
	hosts []string,
	topic string, partitionsCount int,
	timeRange time.Duration,
	name string,
) *Subscriber {
	partitions := make(map[int]*kafka.Reader, partitionsCount)
	for i := 0; i < partitionsCount; i++ {
		partitions[i] = kafka.NewReader(kafka.ReaderConfig{
			Brokers:   hosts,
			Topic:     topic,
			Partition: i,
			MaxBytes:  int(10e6), // 10MB
		})
	}
	return &Subscriber{
		partitions: partitions,
		timeRange:  timeRange,
		timout:     1 * time.Millisecond,
		name:       name,
	}
}

func (s *Subscriber) Subscribe(ctx context.Context) (<-chan geospatial_cache_for_meetup.Message, error) {
	for partition := range s.partitions {
		if err := s.partitions[partition].SetOffsetAt(
			context.Background(),
			time.Now().UTC().Add(-s.timeRange),
		); err != nil {
			return nil, err
		}
	}

	var (
		results = make(chan geospatial_cache_for_meetup.Message, 1000)
		ticker  = time.NewTicker(s.timout)
	)

	for partition := range s.partitions {
		go func(ctx context.Context, partition int) {
			for {
				select {
				case <-ticker.C:

					m, err := s.partitions[partition].ReadMessage(ctx)
					if err != nil {
						results <- geospatial_cache_for_meetup.NewMessage(m.Key, m.Value, m.Partition, m.Offset, err)
						continue
					}
					results <- geospatial_cache_for_meetup.NewMessage(m.Key, m.Value, m.Partition, m.Offset, nil)

				case <-ctx.Done():

					close(results)

					return
				}

			}
		}(ctx, partition)
	}
	return results, nil
}

func (s *Subscriber) Name() string {
	return s.name
}

func (s *Subscriber) Close() error {
	for _, partition := range s.partitions {
		if err := partition.Close(); err != nil {
			return err
		}
	}
	return nil
}
