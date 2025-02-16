package kafka_broadcaster

import (
	"context"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"

	"github.com/ashabykov/geospatial_cache_for_meetup/location"
)

const (
	defaultResultsBufferSize = 1000
)

type Subscriber struct {
	partitions map[int]*kafka.Reader
	timeOffset time.Duration
	timeout    time.Duration
	bufferSize int
}

func NewSubscriber(
	hosts []string,
	topic string, partitionsCount int,
	timeOffset time.Duration,
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
		timeOffset: timeOffset,
		timeout:    time.Nanosecond,
		bufferSize: defaultResultsBufferSize,
	}
}

func (s *Subscriber) Subscribe(ctx context.Context) (<-chan location.Location, error) {
	if err := s.resetOffset(ctx); err != nil {
		return nil, err
	}

	results := make(chan location.Location, defaultResultsBufferSize)
	ticker := time.NewTicker(s.timeout)
	for partition := range s.partitions {
		go func(ctx context.Context, partition int) {
			for {
				select {

				case <-ticker.C:

					message, err := s.partitions[partition].ReadMessage(ctx)
					if err != nil {

						fmt.Println("Subscriber: read message error:", err)

						return
					}

					loc, err := location.Decode(message.Value)
					if err != nil {

						fmt.Println("Subscriber: decode message error:", err)

						continue
					}

					results <- loc

				case <-ctx.Done():

					close(results)

					return
				}
			}
		}(ctx, partition)
	}
	return results, nil
}

func (s *Subscriber) resetOffset(ctx context.Context) error {

	startOffset := time.Now().UTC().Add(-s.timeOffset)

	for partition := range s.partitions {
		if err := s.partitions[partition].SetOffsetAt(ctx, startOffset); err != nil {
			return err
		}
	}
	return nil
}

func (s *Subscriber) Close() error {
	for _, partition := range s.partitions {
		if err := partition.Close(); err != nil {
			return err
		}
	}
	return nil
}
