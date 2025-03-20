package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"

	"github.com/ashabykov/geospatial_cache_for_meetup/cmd"
	"github.com/ashabykov/geospatial_cache_for_meetup/geospatial_client_side_cache/kafka_broadcaster"
	"github.com/ashabykov/geospatial_cache_for_meetup/geospatial_distributed_redis_cache"
	"github.com/ashabykov/geospatial_cache_for_meetup/location"
)

func main() {

	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	var (
		ctx        = cmd.WithContext(context.Background())
		kafkaAddr  = os.Getenv("kafka_addr")
		kafkaTopic = os.Getenv("kafka_topic")
		redisAddr  = os.Getenv("redis_addr")
		rps        = 1000
		pub        = kafka_broadcaster.NewPublisher([]string{kafkaAddr}, kafkaTopic, rps)
		geo        = geospatial_distributed_redis_cache.New(
			redis.NewUniversalClient(&redis.UniversalOptions{
				Addrs:                 []string{redisAddr},
				ReadOnly:              false,
				RouteByLatency:        false,
				RouteRandomly:         true,
				ContextTimeoutEnabled: true,
				ConnMaxIdleTime:       170 * time.Second,
			}),
			10*time.Minute,
		)
	)

	defer func() {
		if err := pub.Close(); err != nil {
			log.Printf("Error closing publisher: %v", err)
		}
	}()

	var (

		//
		count  = 10000
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
		if err := pub.Publish(ctx, loc); err != nil {
			fmt.Println("Publish err:", err)
		}
		if err := geo.Set(loc); err != nil {
			fmt.Println("Set err:", err)
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
