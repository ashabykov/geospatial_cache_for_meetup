package geospatial_distributed_redis_cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"sync"

	"github.com/redis/go-redis/v9"

	"github.com/ashabykov/geospatial_cache_for_meetup/location"
	h3_radius_lookup "github.com/ashabykov/geospatial_cache_for_meetup/location/h3-radius-lookup"
)

type Cache struct {
	redis redis.UniversalClient
}

func New(redis redis.UniversalClient) *Cache {
	return &Cache{redis: redis}
}

func (c *Cache) Set(loc location.Location) error {
	var (
		shardKey = loc.ShardKey()
		ctx      = context.Background()
	)

	// Set contractor geolocation TTL
	if err := c.redis.ZAdd(
		ctx,
		redisListKey(shardKey),
		redis.Z{
			Score:  loc.Ts.Float64(),
			Member: loc.Key(),
		},
	).Err(); err != nil {
		return err
	}

	// Set contractor geolocation
	if err := c.redis.GeoAdd(
		ctx,
		redisLocationKey(shardKey),
		&redis.GeoLocation{
			Latitude:  loc.Lat.Float64(),
			Longitude: loc.Lon.Float64(),
			Name:      loc.Key(),
		},
	).Err(); err != nil {
		return err
	}

	// Set contractor's Location
	val, err := json.Marshal(loc)
	if err != nil {
		return err
	}
	if err = c.redis.Set(
		ctx,
		redisKey(shardKey, loc.Key()),
		val,
		loc.TTL,
	).Err(); err != nil {
		return err
	}
	return err
}

func (c *Cache) Near(target location.Location, radius float64, limit int) ([]location.Location, error) {
	var wg sync.WaitGroup

	var (
		ctx = context.Background()

		keys = h3_radius_lookup.KRingIndexesArea(
			target.Lat.Float64(),
			target.Lon.Float64(),
			radius,
			5,
		)
		results = make(chan location.Location)
	)

	for i := range keys {

		wg.Add(1)

		go func(key string) {

			defer wg.Done()

			locations, err := c.searchByKey(
				ctx,
				key,
				target,
				radius,
				limit,
			)
			if err != nil {
				return
			}
			for _, loc := range locations {
				results <- loc
			}
		}(keys[i])
	}

	go func() {
		wg.Wait()
		close(results)
	}()

	var resp = []location.Location{}

	for res := range results {
		resp = append(resp, res)
	}
	return resp, nil
}

func (c *Cache) searchByKey(
	ctx context.Context,
	indexKey string,
	target location.Location,
	radius float64,
	limit int,
) ([]location.Location, error) {

	values, err := c.read(ctx, indexKey, target, radius, limit)
	if err != nil {
		return nil, err
	}
	locations := make([]location.Location, 0, len(values))
	for i := range values {
		geoLoc, err := parse(values[i])
		if err != nil {
			continue
		}
		locations = append(locations, geoLoc)
	}
	return locations, nil
}

func (c *Cache) read(
	ctx context.Context,
	indexKey string,
	location location.Location,
	radius float64,
	limit int,
) ([]interface{}, error) {
	geoCmd := c.redis.GeoRadius(
		ctx,
		redisLocationKey(indexKey),
		location.Lon.Float64(),
		location.Lat.Float64(),
		&redis.GeoRadiusQuery{
			Radius:    radius,
			Count:     limit,
			Unit:      "m",
			Sort:      "ASC",
			WithDist:  true,
			WithCoord: true,
		},
	)
	if err := geoCmd.Err(); err != nil {
		return nil, err
	}
	locations := geoCmd.Val()
	if len(locations) < 1 {
		return nil, errors.New("locations are not found in geo radius")
	}
	keys := make([]string, 0, len(locations))
	for i := range locations {
		keys = append(keys, redisKey(indexKey, locations[i].Name))
	}
	mgetCmd := c.redis.MGet(ctx, keys...)
	if err := mgetCmd.Err(); err != nil {
		return nil, err
	}
	members := mgetCmd.Val()
	if len(members) < 1 {
		return nil, errors.New("members are not found")
	}
	return members, nil
}

func parse(ptr interface{}) (location.Location, error) {
	val, ok := ptr.([]byte)
	if !ok {
		return location.Location{}, errors.New("Unable to convert string")
	}

	loc := location.Location{}
	if err := json.Unmarshal(val, &loc); err != nil {
		return location.Location{}, errors.New("Unable to convert DTO")
	}
	return loc, nil
}

func redisLocationKey(shardKey string) string {
	return fmt.Sprintf("index:geo:{%s}", shardKey)
}

func redisListKey(shardKey string) string {
	return fmt.Sprintf("index:list:{%s}", shardKey)
}

func redisKey(shardKey string, key string) string {
	return fmt.Sprintf("index:geo:location:{%s}.%s", shardKey, key)
}
