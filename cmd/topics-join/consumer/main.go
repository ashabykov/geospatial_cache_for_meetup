package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/joho/godotenv"

	"github.com/ashabykov/geospatial_cache_for_meetup/cmd"
	"github.com/ashabykov/geospatial_cache_for_meetup/pkg/kafka/consumer/group"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}
	ctx := cmd.WithContext(context.Background())
	kafkaAddr := os.Getenv("kafka_addr")
	kafkaGroupId := os.Getenv("kafka_group_id")
	kafkaTopicJoinLeft := os.Getenv("kafka_topic_join_left")
	kafkaTopicJoinRight := os.Getenv("kafka_topic_join_right")

	leftSub := group.NewSubscriber(
		"left",
		[]string{kafkaAddr},
		kafkaTopicJoinLeft,
		kafkaGroupId,
	)

	rightSub := group.NewSubscriber(
		"right",
		[]string{kafkaAddr},
		kafkaTopicJoinRight,
		kafkaGroupId,
	)

	leftResults, err := leftSub.Subscribe(ctx)
	if err != nil {
		panic(err)
	}

	wg := sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for result := range leftResults {
			fmt.Printf("left: %v \n", result)
		}
	}()

	rightResults, err := rightSub.Subscribe(ctx)
	if err != nil {
		panic(err)
	}

	wg.Add(1)
	go func() {
		defer wg.Done()
		for result := range rightResults {
			fmt.Printf("right: %v\n", result)
		}
	}()

	wg.Wait()

	return
}
