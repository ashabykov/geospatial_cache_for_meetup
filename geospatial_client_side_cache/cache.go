package geospatial_client_side_cache

import (
	"context"
	"fmt"
	"github.com/pkg/errors"
	"time"

	"github.com/ashabykov/geospatial_cache_for_meetup/location"
)

// use as geospatial index
type geospatial interface {
	Add(location location.Location)
	Nearby(target location.Location, radius float64, limit int) []location.Name
	Remove(location location.Location)
}

// use as sorted timestamp index
type timestamp interface {
	Add(location location.Location)
	Read(from, to location.Timestamp) []location.Name
	Remove(location location.Location)
}

// use as in-memory key/val storage
type cache interface {
	Set(string, location.Location, time.Duration) bool
	Get(string) (location.Location, bool)
	Del(string)
	TTL() time.Duration
}

type Cache struct {
	geospatial   geospatial
	timestamp    timestamp
	cache        cache
	cleanTimeout time.Duration
	cleanRange   time.Duration
}

func New(ctx context.Context, g geospatial, t timestamp, c cache) *Cache {
	ccc := &Cache{
		cache:        c,
		geospatial:   g,
		timestamp:    t,
		cleanTimeout: 5 * time.Second,
		cleanRange:   c.TTL(),
	}

	go ccc.clean(ctx)

	return ccc
}

func (c *Cache) clean(ctx context.Context) {

	ticker := time.NewTicker(c.cleanTimeout)

	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			var (
				from = location.Timestamp(0)
				to   = location.Timestamp(time.Now().UTC().Add(-c.cleanRange).Unix())
			)
			for _, name := range c.timestamp.Read(from, to) {
				if loc, ok := c.cache.Get(name.String()); ok {

					c.Del(loc)

					fmt.Println("Removed from cache: ", loc)
				}
			}
		case <-ctx.Done():
			return
		}
	}
}

func (c *Cache) Near(target location.Location, radius float64, limit int) ([]location.Location, error) {
	var (
		now  = time.Now().UTC()
		from = location.Timestamp(now.Add(-c.cache.TTL()).Unix())
		to   = location.Timestamp(now.Unix())
	)

	loc1 := c.timestamp.Read(from, to)

	loc2 := c.geospatial.Nearby(target, radius, limit)

	if len(loc1) == 0 || len(loc2) == 0 {
		return []location.Location{}, errors.New("no locations found")
	}

	if len(loc1) > len(loc2) {
		return c.get(intersect(loc2, loc1)...), nil
	}

	return c.get(intersect(loc1, loc2)...), nil
}

func (c *Cache) Get(name location.Name) (location.Location, bool) {
	return c.cache.Get(name.String())
}

func (c *Cache) Set(target location.Location) error {
	c.cache.Set(target.Key(), target, target.TTL)
	c.timestamp.Add(target)
	c.geospatial.Add(target)
	return nil
}

func (c *Cache) Del(target location.Location) {
	c.geospatial.Remove(target)
	c.timestamp.Add(target)
	c.cache.Del(target.Key())
}

func (c *Cache) get(names ...location.Name) []location.Location {

	ret := make([]location.Location, 0, len(names))
	for _, name := range names {
		if loc, ok := c.cache.Get(name.String()); ok {
			ret = append(ret, loc)
		}
	}
	return ret
}

func intersect(shortest, longest []location.Name) []location.Name {

	hs := make(map[location.Name]struct{}, len(shortest))
	for _, name := range shortest {
		hs[name] = struct{}{}
	}

	ret := make([]location.Name, 0, len(shortest))
	for _, name := range longest {
		if _, ok := hs[name]; ok {
			ret = append(ret, name)
		}
	}
	return ret
}
