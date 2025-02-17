package main

import (
	"context"

	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"

	"github.com/ashabykov/geospatial_cache_for_meetup/fan-out-read-client"
	"github.com/ashabykov/geospatial_cache_for_meetup/fan-out-write-client"
	"github.com/ashabykov/geospatial_cache_for_meetup/geospatial_client_side_cache"
	"github.com/ashabykov/geospatial_cache_for_meetup/geospatial_client_side_cache/kafka_broadcaster"
	"github.com/ashabykov/geospatial_cache_for_meetup/geospatial_client_side_cache/lru_cache"
	"github.com/ashabykov/geospatial_cache_for_meetup/geospatial_client_side_cache/rtree_index"
	"github.com/ashabykov/geospatial_cache_for_meetup/geospatial_client_side_cache/sorted_set"
	"github.com/ashabykov/geospatial_cache_for_meetup/geospatial_distributed_redis_cache"
)

func Init(ctx context.Context) (*fatout_read_client.Client, *fanout_write_client.Client) {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	var (
		kafkaAddr     = os.Getenv("kafka_addr")
		kafkaTopic    = os.Getenv("kafka_topic")
		partitions, _ = strconv.Atoi(os.Getenv("partitions"))
		timeOffset    = 10 * time.Minute
		ttl           = 10 * time.Minute
		capacity      = 10000
		redisAddr     = os.Getenv("redis_addr")
		geoV1         = geospatial_distributed_redis_cache.New(
			redis.NewUniversalClient(&redis.UniversalOptions{
				Addrs:                 []string{redisAddr},
				ReadOnly:              false,
				RouteByLatency:        false,
				RouteRandomly:         true,
				ContextTimeoutEnabled: true,
				ConnMaxIdleTime:       170 * time.Second,
			}),
			ttl,
		)
		sub = kafka_broadcaster.NewSubscriber(
			[]string{kafkaAddr},
			kafkaTopic,
			partitions,
			timeOffset,
		)
		geoV2 = geospatial_client_side_cache.New(
			ctx,
			rtree_index.NewIndex(),
			sorted_set.New(),
			lru_cache.New(ttl, capacity),
		)
	)
	return fatout_read_client.New(geoV1), fanout_write_client.New(sub, geoV2)
}
