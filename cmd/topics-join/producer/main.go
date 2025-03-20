package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/ashabykov/geospatial_cache_for_meetup/cmd"
	"github.com/ashabykov/geospatial_cache_for_meetup/location"
	"github.com/ashabykov/geospatial_cache_for_meetup/pkg/kafka"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	var (
		rps                 = 1000
		ctx                 = cmd.WithContext(context.Background())
		kafkaAddr           = os.Getenv("kafka_addr")
		kafkaTopicJoinLeft  = os.Getenv("kafka_topic_join_left")
		kafkaTopicJoinRight = os.Getenv("kafka_topic_join_right")
		leftPub             = kafka.NewPublisher([]string{kafkaAddr}, kafkaTopicJoinLeft, rps)
		rightPub            = kafka.NewPublisher([]string{kafkaAddr}, kafkaTopicJoinRight, rps)
	)

	defer func() {
		if err := leftPub.Close(); err != nil {
			log.Printf("Error closing publisher: %v", err)
		}
		if err := rightPub.Close(); err != nil {
			log.Printf("Error closing publisher: %v", err)
		}
	}()

	var (

		//
		count  = 100000
		radius = float64(5000)
		target = location.Location{
			Name: "target",
			Lat:  43.244555,
			Lon:  76.940012,
			Ts:   location.Timestamp(time.Now().UTC().Unix()),
			TTL:  10 * time.Minute,
		}
	)
	for i, loc := range locations(count, radius, target) {
		loc.Ts = location.Now()
		loc.Name = location.Name(fmt.Sprintf("location-%d", i%100))
		if err := leftPub.Publish(ctx, loc); err != nil {
			fmt.Println("Publish err:", err)
		}
		if err := rightPub.Publish(ctx, loc); err != nil {
			fmt.Println("Publish err:", err)
		}
	}

	fmt.Println("Published messages:", count)
}

func locations(count int, radius float64, center location.Location) []location.Location {
	msgs := make([]location.Location, 0, count)
	for i := 0; i < count; i++ {
		loc, _ := location.Generate(
			center,
			radius,
		)
		msgs = append(msgs, loc)
	}
	return msgs
}
