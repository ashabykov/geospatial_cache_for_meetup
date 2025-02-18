package main

import (
	"context"
	"encoding/json"
	"net/http"

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

type handler struct {
	fatOutReadClient  *fatout_read_client.Client
	fanOutWriteClient *fanout_write_client.Client
}

func (h *handler) Star(ctx context.Context) {
	go h.fanOutWriteClient.SubscribeOnUpdates(ctx)
}

func (h *handler) FanOutReadClientHandler(w http.ResponseWriter, r *http.Request) {
	query := NearBy{}
	err := json.NewDecoder(r.Body).Decode(&query)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	results, err := h.fatOutReadClient.Near(query.Location, query.Radius, query.Limit)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	jsonItem, err := json.Marshal(results)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonItem)
}

func (h *handler) FanOutWriteClientHandler(w http.ResponseWriter, r *http.Request) {
	query := NearBy{}
	err := json.NewDecoder(r.Body).Decode(&query)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	results, err := h.fanOutWriteClient.Near(query.Location, query.Radius, query.Limit)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	jsonItem, err := json.Marshal(results)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(jsonItem)
}

func New(ctx context.Context) *handler {
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
	return &handler{
		fatOutReadClient:  fatout_read_client.New(geoV1),
		fanOutWriteClient: fanout_write_client.New(sub, geoV2),
	}
}
