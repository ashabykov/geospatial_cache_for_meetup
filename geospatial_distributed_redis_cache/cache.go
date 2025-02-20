package geospatial_distributed_redis_cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"

	"github.com/ashabykov/geospatial_cache_for_meetup/location"
)

type Cache struct {
	redis redis.UniversalClient

	ttl time.Duration
}

func New(redis redis.UniversalClient, ttl time.Duration) *Cache {
	return &Cache{redis: redis, ttl: ttl}
}

func (c *Cache) Set(loc location.Location) error {
	var (
		shardKey = loc.GeoHash()
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
		ctx     = context.Background()
		shards  = target.NearHex(5)[:3]
		results = make(chan location.Location)
	)

	for i := range shards {

		wg.Add(1)

		go func(key string) {

			defer wg.Done()

			locations, err := c.readShard(
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
		}(shards[i])
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

func (c *Cache) readShard(
	ctx context.Context,
	shardKey string,
	target location.Location,
	radius float64,
	limit int,
) ([]location.Location, error) {

	values, err := c.read(ctx, shardKey, target, radius, limit)
	if err != nil {
		return nil, err
	}
	locations := make([]location.Location, 0, len(values))
	for i := range values {
		if geoLoc, err := parse(values[i]); err == nil {
			locations = append(locations, geoLoc)
		}
	}
	return locations, nil
}

func (c *Cache) read(
	ctx context.Context,
	shardKey string,
	location location.Location,
	radius float64,
	limit int,
) ([]interface{}, error) {
	geoCmd := c.redis.GeoRadius(
		ctx,
		redisLocationKey(shardKey),
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

	geoMembers := geoCmd.Val()
	if len(geoMembers) == 0 {
		return nil, errors.New("locations are not found in geo radius")
	}
	locNames1 := make([]string, 0, len(geoMembers))
	for i := range geoMembers {
		locNames1 = append(locNames1, geoMembers[i].Name)
	}

	var (
		now  = time.Now().UTC()
		from = strconv.FormatInt(now.Add(-c.ttl).Unix(), 10)
		to   = strconv.FormatInt(now.Unix(), 10)
	)

	zCmd := c.redis.ZRangeByScore(
		ctx,
		redisListKey(shardKey),
		&redis.ZRangeBy{
			Min:    from,
			Max:    to,
			Offset: 0,
			Count:  0,
		},
	)
	if err := zCmd.Err(); err != nil {
		return nil, err
	}

	locNames2 := zCmd.Val()
	if len(locNames2) == 0 {
		return nil, errors.New("locations are not found in geo radius")
	}

	names := intersect(locNames1, locNames2)
	keys := make([]string, 0, len(names))
	for i := range names {
		keys = append(keys, redisKey(shardKey, names[i]))
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
	val, ok := ptr.(string)
	if !ok {
		return location.Location{}, errors.New("Unable to convert string")
	}

	loc := location.Location{}
	if err := json.Unmarshal([]byte(val), &loc); err != nil {
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

func intersect(shortest, longest []string) []string {

	hs := make(map[string]struct{}, len(shortest))
	for _, name := range shortest {
		hs[name] = struct{}{}
	}

	ret := make([]string, 0, len(shortest))
	for _, name := range longest {
		if _, ok := hs[name]; ok {
			ret = append(ret, name)
		}
	}
	return ret
}
