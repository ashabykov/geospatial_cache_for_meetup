package main

import (
	"context"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"

	"github.com/ashabykov/geospatial_cache_for_meetup"
	"github.com/ashabykov/geospatial_cache_for_meetup/geospatial_cache"
	"github.com/ashabykov/geospatial_cache_for_meetup/kafka_broadcaster"
	"github.com/ashabykov/geospatial_cache_for_meetup/lru_cache"
	"github.com/ashabykov/geospatial_cache_for_meetup/rtree_index"
	"github.com/ashabykov/geospatial_cache_for_meetup/sorted_set"
)

func Init(ctx context.Context) *geospatial_cache_for_meetup.Client {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	var (
		addr          = os.Getenv("addr")
		topic         = os.Getenv("topic")
		partitions, _ = strconv.Atoi(os.Getenv("partitions"))
		timeOffset    = 10 * time.Minute
		ttl           = 10 * time.Minute
		capacity      = 10000
		sub           = kafka_broadcaster.NewSubscriber(
			[]string{addr},
			topic,
			partitions,
			timeOffset,
		)
		geo = geospatial_cache.New(
			ctx,
			rtree_index.NewIndex(),
			sorted_set.New(),
			lru_cache.New(ttl, capacity),
		)
	)

	client := geospatial_cache_for_meetup.New(sub, geo)

	return client
}
