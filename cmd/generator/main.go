package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"

	"github.com/ashabykov/geospatial_cache_for_meetup/kafka_broadcaster"
	"github.com/ashabykov/geospatial_cache_for_meetup/location"
)

func main() {

	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	var (
		ctx   = context.Background()
		addr  = os.Getenv("addr")
		topic = os.Getenv("topic")
		rps   = 30
		pub   = kafka_broadcaster.NewPublisher([]string{addr}, topic, rps)
	)

	defer func() {
		if err := pub.Close(); err != nil {
			log.Printf("Error closing publisher: %v", err)
		}
	}()

	var (
		//
		count  = 1000
		radius = float64(8000)
		target = location.Location{
			Lat: 43.244555,
			Lon: 76.940012,
			TTL: 10 * time.Minute,
		}
	)
	if err := pub.Publish(ctx, locations(count, radius, target)...); err != nil {
		fmt.Println("Publish err:", err, "\n")
	}

	fmt.Println("Published messages:", count, "\n")
}

func locations(count int, radius float64, center location.Location) []location.Location {
	msgs := make([]location.Location, 0, count)
	for i := 0; i < count; i++ {
		loc := location.Generate(
			center,
			radius,
		)
		msgs = append(msgs, loc)
	}
	return msgs
}
