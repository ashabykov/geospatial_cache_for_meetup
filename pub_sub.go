package geospatial_cache_for_meetup

import (
	"context"
	"fmt"
)

type Message struct {
	Key, Val  []byte
	Offset    int64
	Partition int
	Err       error
}

func (m Message) String() string {
	return fmt.Sprintf("{Key:%s, Val:%s, Partition:%d Offset:%d, Err:%v}", string(m.Key), string(m.Val), m.Partition, m.Offset, m.Err)
}

func NewMessage(key, val []byte, partition int, offset int64, err error) Message {
	return Message{
		Key:       key,
		Val:       val,
		Partition: partition,
		Offset:    offset,
		Err:       err,
	}
}

type Pub interface {
	Publish(ctx context.Context, m ...Message) error
	Close() error
}

type Sub interface {
	Name() string
	Subscribe(ctx context.Context) (<-chan Message, error)
	Close() error
}
