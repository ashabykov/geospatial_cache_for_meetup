package fanout_write_client

import (
	"context"
	"log"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/joho/godotenv"

	"github.com/ashabykov/geospatial_cache_for_meetup/cmd"
	"github.com/ashabykov/geospatial_cache_for_meetup/geospatial_client_side_cache"
	"github.com/ashabykov/geospatial_cache_for_meetup/geospatial_client_side_cache/kafka_broadcaster"
	"github.com/ashabykov/geospatial_cache_for_meetup/geospatial_client_side_cache/lru_cache"
	"github.com/ashabykov/geospatial_cache_for_meetup/geospatial_client_side_cache/rtree_index"
	"github.com/ashabykov/geospatial_cache_for_meetup/geospatial_client_side_cache/sorted_set"
	"github.com/ashabykov/geospatial_cache_for_meetup/location"
)

func BenchmarkClientNear_for_FunOutWrite(b *testing.B) {

	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	var (
		ctx           = cmd.WithContext(context.Background())
		kafka_addr    = os.Getenv("kafka_addr")
		kafka_topic   = os.Getenv("kafka_topic")
		partitions, _ = strconv.Atoi(os.Getenv("partitions"))
		timeOffset    = 5 * time.Minute
		ttl           = 20 * time.Minute
		capacity      = 10000
		sub           = kafka_broadcaster.NewSubscriber(
			[]string{kafka_addr},
			kafka_topic,
			partitions,
			timeOffset,
		)
		geo = geospatial_client_side_cache.New(
			ctx,
			rtree_index.NewIndex(),
			sorted_set.New(),
			lru_cache.New(ttl, capacity),
		)
		client = New(
			sub,
			geo,
		)
	)

	client.SubscribeOnUpdates(ctx)

	var (
		limit  = 1000
		radius = float64(4000)
		target = location.Location{
			Name: "target",
			Lat:  43.244555,
			Lon:  76.940012,
			Ts:   location.Timestamp(time.Now().UTC().Unix()),
			TTL:  10 * time.Minute,
		}
	)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			got, err := client.Near(target, radius, limit)
			if err != nil {
				log.Fatal(err)
			}
			println(len(got))
		}
	})
}
