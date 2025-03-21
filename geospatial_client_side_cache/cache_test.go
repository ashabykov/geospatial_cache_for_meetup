package geospatial_client_side_cache

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/ashabykov/geospatial_cache_for_meetup/geospatial_client_side_cache/lru_cache"
	"github.com/ashabykov/geospatial_cache_for_meetup/geospatial_client_side_cache/rtree_index"
	"github.com/ashabykov/geospatial_cache_for_meetup/geospatial_client_side_cache/sorted_set"
	"github.com/ashabykov/geospatial_cache_for_meetup/location"
)

func TestCache_Del(t *testing.T) {
	now := time.Now().UTC()
	tests := []struct {
		name      string
		radius    float64
		limit     int
		locations []location.Location
		target    location.Location
	}{
		{
			name:   "success",
			radius: 10,
			limit:  10,
			target: location.Location{
				Name: "awesome-1-near",
				Ts:   newTimestamp(now, -3*time.Minute),
				TTL:  15 * time.Minute,
				Lat:  43.241705,
				Lon:  76.909756,
			},
			locations: []location.Location{
				{
					Name: "awesome-1-near",
					Ts:   newTimestamp(now, -3*time.Minute),
					TTL:  15 * time.Minute,
					Lat:  43.241705,
					Lon:  76.909756,
				},
				{
					Name: "awesome-2-far",
					Ts:   newTimestamp(now, -3*time.Minute),
					TTL:  15 * time.Minute,
					Lat:  43.248489,
					Lon:  76.923511,
				},
				{
					Name: "awesome-3-near",
					Ts:   newTimestamp(now, -3*time.Minute),
					TTL:  15 * time.Minute,
					Lat:  43.246410,
					Lon:  76.916558,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			idx := New(
				context.Background(),
				rtree_index.NewIndex(),
				sorted_set.New(),
				lru_cache.New(tt.target.TTL, 10000),
			)

			for i := range tt.locations {
				idx.Set(tt.locations[i])
			}

			// testing method
			idx.Del(tt.target)

			// check deleted
			got, exist := idx.Get(tt.target.Name)
			assert.Equal(t, location.Location{}, got)
			assert.False(t, exist)

			locations, err := idx.Near(tt.target, tt.radius, tt.limit)
			assert.NotNil(t, err)
			assert.Equal(t, 0, len(locations))
		})
	}
}

func TestCache_Near(t *testing.T) {
	now := time.Now().UTC()
	tests := []struct {
		name   string
		target location.Location
		radius float64
		limit  int

		locations []location.Location
		expected  []location.Location
	}{
		{
			name: "near locations",
			target: location.Location{
				Name: "target",
				Ts:   newTimestamp(now, 0*time.Minute),
				TTL:  15 * time.Minute,
				Lat:  43.246645,
				Lon:  76.909713,
			},
			radius: 1000,
			limit:  5,
			locations: []location.Location{
				{
					Name: "awesome-1-near",
					Ts:   newTimestamp(now, -3*time.Minute),
					TTL:  15 * time.Minute,
					Lat:  43.241705,
					Lon:  76.909756,
				},
				{
					Name: "awesome-2-far",
					Ts:   newTimestamp(now, -3*time.Minute),
					TTL:  15 * time.Minute,
					Lat:  43.248489,
					Lon:  76.923511,
				},
				{
					Name: "awesome-3-near",
					Ts:   newTimestamp(now, -3*time.Minute),
					TTL:  15 * time.Minute,
					Lat:  43.246410,
					Lon:  76.916558,
				},
				{
					Name: "awesome-4-far",
					Ts:   newTimestamp(now, -3*time.Minute),
					TTL:  15 * time.Minute,
					Lat:  43.222179,
					Lon:  76.798691,
				},
				{
					Name: "awesome-5-expired",
					Ts:   newTimestamp(now, -20*time.Minute),
					TTL:  15 * time.Minute,
					Lat:  43.246410,
					Lon:  76.916558,
				},
			},
			expected: []location.Location{
				{
					Name: "awesome-1-near",
					Ts:   newTimestamp(now, -3*time.Minute),
					TTL:  15 * time.Minute,
					Lat:  43.241705,
					Lon:  76.909756,
				},
				{
					Name: "awesome-3-near",
					Ts:   newTimestamp(now, -3*time.Minute),
					TTL:  15 * time.Minute,
					Lat:  43.246410,
					Lon:  76.916558,
				},
			},
		},
		{
			name: "the most near locations",
			target: location.Location{
				Name: "target",
				Ts:   newTimestamp(now, -5*time.Minute),
				TTL:  15 * time.Minute,
				Lat:  43.248723,
				Lon:  76.923489,
			},
			radius: 1000,
			limit:  1,
			locations: []location.Location{
				{
					Name: "awesome-1-near",
					Ts:   newTimestamp(now, -3*time.Minute),
					TTL:  15 * time.Minute,
					Lat:  43.247192,
					Lon:  76.923875,
				},
				{
					Name: "awesome-2-the-most-near",
					Ts:   newTimestamp(now, -3*time.Minute),
					TTL:  15 * time.Minute,
					Lat:  43.248723,
					Lon:  76.923489,
				},
				{
					Name: "awesome-3-near",
					Ts:   newTimestamp(now, -3*time.Minute),
					TTL:  15 * time.Minute,
					Lat:  43.246410,
					Lon:  76.916558,
				},
				{
					Name: "awesome-4-far",
					Ts:   newTimestamp(now, -3*time.Minute),
					TTL:  15 * time.Minute,
					Lat:  43.222179,
					Lon:  76.798691,
				},
			},
			expected: []location.Location{
				{
					Name: "awesome-2-the-most-near",
					Ts:   newTimestamp(now, -3*time.Minute),
					TTL:  15 * time.Minute,
					Lat:  43.248723,
					Lon:  76.923489,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			idx := New(
				context.Background(),
				rtree_index.NewIndex(),
				sorted_set.New(),
				lru_cache.New(tt.target.TTL, 10000),
			)

			for i := range tt.locations {
				idx.Set(tt.locations[i])
			}

			got, err := idx.Near(tt.target, tt.radius, tt.limit)
			assert.Nil(t, err)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func newTimestamp(now time.Time, daly time.Duration) location.Timestamp {
	return location.Timestamp(now.Add(daly).Unix())
}
