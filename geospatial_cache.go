package geospatial_cache_for_meetup

import (
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
}

type Cache struct {
	geospatial geospatial
	timestamp  timestamp
	cache      cache
}

func New(g geospatial, t timestamp, c cache) *Cache {
	return &Cache{cache: c, geospatial: g, timestamp: t}
}

func (c *Cache) Near(target location.Location, radius float64, limit int) []location.Location {
	var (
		now  = time.Now()
		from = location.Timestamp(now.UTC().Add(-target.TTL).Unix())
		to   = location.Timestamp(now.UTC().Unix())
	)

	loc1 := c.timestamp.Read(from, to)

	loc2 := c.geospatial.Nearby(target, radius, limit)

	if len(loc1) == 0 && len(loc2) == 0 {
		return []location.Location{}
	}

	if len(loc1) > len(loc2) {
		return c.get(intersect(loc2, loc1)...)
	}

	return c.get(intersect(loc1, loc2)...)
}

func (c *Cache) Get(name location.Name) (location.Location, bool) {
	return c.cache.Get(name.String())
}

func (c *Cache) Set(target location.Location) {
	c.geospatial.Add(target)
	c.timestamp.Add(target)
	c.cache.Set(target.Key(), target, target.TTL)
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

func intersect(self, other []location.Name) []location.Name {

	hs := make(map[location.Name]struct{}, len(self))
	for _, name := range self {
		hs[name] = struct{}{}
	}

	ret := make([]location.Name, 0, len(self))
	for _, name := range other {
		if _, ok := hs[name]; ok {
			ret = append(ret, name)
		}
	}
	return ret
}
