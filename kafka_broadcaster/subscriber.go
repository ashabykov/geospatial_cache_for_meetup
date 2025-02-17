package kafka_broadcaster

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/segmentio/kafka-go"

	"github.com/ashabykov/geospatial_cache_for_meetup/location"
)

const (
	defaultResultsBufferSize = 1000
)

type Subscriber struct {
	partitions []*kafka.Reader
	timeOffset time.Duration
	timeout    time.Duration
	bufferSize int
}

func NewSubscriber(
	hosts []string,
	topic string, partitionsCount int,
	timeOffset time.Duration,
) *Subscriber {
	partitions := make([]*kafka.Reader, partitionsCount)
	for i := 0; i < partitionsCount; i++ {
		partitions[i] = kafka.NewReader(kafka.ReaderConfig{
			Brokers:          hosts,
			Topic:            topic,
			Partition:        i,
			MaxBytes:         int(10e6), // 10MB
			ReadBatchTimeout: 100 * time.Millisecond,
		})
	}
	return &Subscriber{
		partitions: partitions,
		timeOffset: timeOffset,
		timeout:    10 * time.Millisecond,
		bufferSize: defaultResultsBufferSize,
	}
}

func (s *Subscriber) Subscribe(ctx context.Context) (<-chan location.Location, error) {
	if err := s.resetOffset(ctx); err != nil {
		return nil, err
	}

	results := make(chan location.Location, defaultResultsBufferSize)

	wg := &sync.WaitGroup{}
	for _, reader := range s.partitions {

		wg.Add(1)

		go s.worker(ctx, results, reader, wg)
	}

	go func() {

		wg.Wait()

		close(results)
	}()

	return results, nil
}

func (s *Subscriber) worker(ctx context.Context, results chan<- location.Location, reader *kafka.Reader, wg *sync.WaitGroup) {

	defer wg.Done()

	ticker := time.NewTicker(s.timeout)

	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			message, err := reader.ReadMessage(ctx)
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
