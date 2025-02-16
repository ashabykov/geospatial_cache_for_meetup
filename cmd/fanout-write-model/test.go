package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"

	"github.com/ashabykov/geospatial_cache_for_meetup"
	"github.com/ashabykov/geospatial_cache_for_meetup/cmd"
	"github.com/ashabykov/geospatial_cache_for_meetup/geospatial_cache"
	"github.com/ashabykov/geospatial_cache_for_meetup/kafka_broadcaster"
	"github.com/ashabykov/geospatial_cache_for_meetup/location"
	"github.com/ashabykov/geospatial_cache_for_meetup/lru_cache"
	"github.com/ashabykov/geospatial_cache_for_meetup/rtree_index"
	"github.com/ashabykov/geospatial_cache_for_meetup/sorted_set"
)

func main() {

	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	var (
		ctx           = cmd.WithContext(context.Background())
		addr          = os.Getenv("addr")
		topic         = os.Getenv("topic")
		partitions, _ = strconv.Atoi(os.Getenv("partitions"))
		timeOffset    = 15 * time.Minute
		ttl           = 20 * time.Minute
		capacity      = 10000
		subscriber    = kafka_broadcaster.NewSubscriber(
			[]string{addr},
			topic,
			partitions,
			timeOffset,
		)
		geospatialCache = geospatial_cache.New(
			rtree_index.NewIndex(),
			sorted_set.New(),
			lru_cache.New(ttl, capacity),
		)
		client = geospatial_cache_for_meetup.New(
			subscriber,
			geospatialCache,
		)
	)

	go client.Subscribe(ctx)

	var (
		limit  = 1000
		radius = float64(4000)
		target = location.Location{
			Name: "target",
			Lat:  43.244555,
			Lon:  76.940012,
			Ts:   location.Timestamp(time.Now().Unix()),
			TTL:  10 * time.Minute,
		}
	)

	for i := 0; i < 10; i++ {

		fmt.Println("Found locations: ", len(client.Near(target, radius, limit)))

		time.Sleep(1 * time.Second)
	}

}
