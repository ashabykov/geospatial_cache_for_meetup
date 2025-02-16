//  based on impl: https://github.com/hashicorp/golang-lru/v2/expirable

package lru_cache

import (
	"time"

	"github.com/hashicorp/golang-lru/v2/expirable"

	"github.com/ashabykov/geospatial_cache_for_meetup/location"
)

const (
	defaultCapacity = 100000
)

type Cache struct {
	base *expirable.LRU[string, location.Location]
}

func New(ttl time.Duration, capacity int) *Cache {
	if capacity == 0 {
		capacity = defaultCapacity
	}

	return &Cache{
		base: expirable.NewLRU[string, location.Location](capacity, nil, ttl),
	}
}

func (c *Cache) Get(key string) (location.Location, bool) {
	return c.base.Get(key)
}

func (c *Cache) Set(key string, value location.Location, _ time.Duration) bool {
	return c.base.Add(key, value)
}

func (c *Cache) Del(key string) {
	c.base.Remove(key)
}
