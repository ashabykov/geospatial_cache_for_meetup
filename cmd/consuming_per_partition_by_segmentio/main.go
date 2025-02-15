package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/ashabykov/geospatial_cache_for_meetup"
	"github.com/ashabykov/geospatial_cache_for_meetup/kafka_broadcast"
)

func main() {
	var (
		ctx   = context.Background()
		addr  = []string{"localhost:9092"}
		topic = "awesom.orders"
		pub   = kafka_broadcast.NewPublisher(addr, topic)

		sub3 = kafka_broadcast.NewSubscriber(addr, topic, 1*time.Second, "Consumer 3")
		sub1 = kafka_broadcast.NewSubscriber(addr, topic, 1*time.Second, "Consumer 1")
		sub2 = kafka_broadcast.NewSubscriber(addr, topic, 1*time.Second, "Consumer 2")
	)

	defer pub.Close()
	defer sub1.Close()
	defer sub2.Close()
	defer sub3.Close()

	time.Sleep(1 * time.Second)
	if err := pub.Publish(ctx, messages("pack 1", 100)...); err != nil {
		fmt.Println("Publish err:", err, "\n")
	}

	time.Sleep(1 * time.Second)
	if err := pub.Publish(ctx, messages("pack 2", 100)...); err != nil {
		fmt.Println("Publish err:", err, "\n")
	}

	time.Sleep(1 * time.Second)
	if err := pub.Publish(ctx, messages("pack 3", 30)...); err != nil {
		fmt.Println("Publish err:", err, "\n")
	}
	var wg sync.WaitGroup
	wg.Add(3)

	printLn := func(ctx context.Context, sub geospatial_cache_for_meetup.Sub) {

		defer wg.Done()

		results, err := sub.Subscribe(ctx)
		if err != nil {
			fmt.Println("Subscribe err:", err)
			return
		}

		for msg := range results {
			fmt.Println(sub.Name(), msg)
		}
	}

	go printLn(ctx, sub1)
	go printLn(ctx, sub2)
	go printLn(ctx, sub3)

	wg.Wait()
}

func messages(pref string, count int) []geospatial_cache_for_meetup.Message {
	msgs := make([]geospatial_cache_for_meetup.Message, 0, count)
	for i := 0; i < count; i++ {
		msgs = append(msgs, geospatial_cache_for_meetup.Message{
			Key: []byte(fmt.Sprintf("%d", i)),
			Val: []byte(fmt.Sprintf("order %s value %d", pref, i)),
		})
	}
	return msgs
}
