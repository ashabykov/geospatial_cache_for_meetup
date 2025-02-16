package geospatial_cache_for_meetup

import (
	"context"
	"log"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/joho/godotenv"

	"github.com/ashabykov/geospatial_cache_for_meetup/cmd"
	"github.com/ashabykov/geospatial_cache_for_meetup/geospatial_cache"
	"github.com/ashabykov/geospatial_cache_for_meetup/kafka_broadcaster"
	"github.com/ashabykov/geospatial_cache_for_meetup/location"
	"github.com/ashabykov/geospatial_cache_for_meetup/lru_cache"
	"github.com/ashabykov/geospatial_cache_for_meetup/rtree_index"
	"github.com/ashabykov/geospatial_cache_for_meetup/sorted_set"
)

func BenchmarkClientNear_for_FunOutWrite(b *testing.B) {

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
		sub           = kafka_broadcaster.NewSubscriber(
			[]string{addr},
			topic,
			partitions,
			timeOffset,
		)
		geo = geospatial_cache.New(
			rtree_index.NewIndex(),
			sorted_set.New(),
			lru_cache.New(ttl, capacity),
		)
		client = New(
			sub,
			geo,
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

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			client.Near(target, radius, limit)
		}
	})
}
