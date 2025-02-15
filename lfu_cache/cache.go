//  based on impl: https://github.com/hypermodeinc/ristretto

package lfu_cache

import (
	"time"

	"github.com/dgraph-io/ristretto/v2"

	"github.com/ashabykov/geospatial_cache_for_meetup/location"
)

type Cache struct {
	base *ristretto.Cache[string, location.Location]
}

func New() *Cache {
	base, _ := ristretto.NewCache(&ristretto.Config[string, location.Location]{
		NumCounters: 1e7,     // number of keys to track frequency of (10M).
		MaxCost:     1 << 30, // maximum cost of cache (1GB).
		BufferItems: 64,      // number of keys per Get buffer.
	})
	return &Cache{base: base}
}

func (c *Cache) Get(key string) (location.Location, bool) {
	return c.base.Get(key)
}

func (c *Cache) Set(key string, value location.Location, ttl time.Duration) bool {
	return c.base.SetWithTTL(key, value, 1, ttl)
}

func (c *Cache) Del(key string) {
	c.base.Del(key)
}
