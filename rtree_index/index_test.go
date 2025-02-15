package rtree_index

import (
	"github.com/ashabykov/geospatial_cache_for_meetup/location"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndex_Nearby(t *testing.T) {
	tests := []struct {
		name      string
		radius    float64
		limit     int
		locations []location.Location
		target    location.Location
		dist      float64
	}{
		{
			name:   "nearby",
			radius: 1000.0,
			limit:  2,
			locations: []location.Location{
				{
					Name: "location 1",
					Lat:  43.244555,
					Lon:  76.940012,
				},
				{
					Name: "location 2",
					Lat:  43.244331,
					Lon:  76.929712,
				},
				{
					Name: "location 3",
					Lat:  43.226188,
					Lon:  76.869333,
				},
				{
					Name: "location 4",
					Lat:  43.256870,
					Lon:  76.893835,
				},
				{
					Name: "location 5",
					Lat:  43.256880,
					Lon:  76.893837,
				},
				{
					Name: "location 6",
					Lat:  43.256890,
					Lon:  76.893838,
				},
			},
			target: location.Location{
				Name: "target",
				Lat:  43.248752,
				Lon:  76.932523,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := NewIndex()
			for _, l := range tt.locations {
				i.Add(l)
			}
			got := i.Nearby(tt.target, tt.radius, tt.limit)
			assert.Equal(t, 2, len(got))
		})
	}
}

func TestIndex_Remove(t *testing.T) {
	tests := []struct {
		name      string
		locations []location.Location
		target    location.Location
		size      int
	}{
		{

			name: "must remove",
			locations: []location.Location{
				{
					Name: "location 1",
					Lat:  43.244555,
					Lon:  76.940012,
				},
				{
					Name: "location 2",
					Lat:  43.244331,
					Lon:  76.929712,
				},
				{
					Name: "location 3",
					Lat:  43.226188,
					Lon:  76.869333,
				},
			},
			target: location.Location{
				Name: "location 2",
				Lat:  43.244331,
				Lon:  76.929712,
			},
			size: 2,
		},
		{
			name: "not remove, no match by name",
			locations: []location.Location{
				{
					Name: "location 1",
					Lat:  43.244555,
					Lon:  76.940012,
				},
				{
					Name: "location 2",
					Lat:  43.244331,
					Lon:  76.929712,
				},
				{
					Name: "location 3",
					Lat:  43.226188,
					Lon:  76.869333,
				},
			},
			target: location.Location{
				Name: "location N",
				Lat:  43.244331,
				Lon:  76.929712,
			},
			size: 3,
		},
		{
			name: "not remove, no match by coordinates",
			locations: []location.Location{
				{
					Name: "location 1",
					Lat:  43.244555,
					Lon:  76.940012,
				},
				{
					Name: "location 2",
					Lat:  43.244331,
					Lon:  76.929712,
				},
				{
					Name: "location 3",
					Lat:  43.226188,
					Lon:  76.869333,
				},
			},
			target: location.Location{
				Name: "location 2",
				Lat:  43.226188,
				Lon:  76.869333,
			},
			size: 3,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := NewIndex()
			for _, l := range tt.locations {
				i.Add(l)
			}
			i.Remove(tt.target)
			assert.Equal(t, tt.size, i.Size())
		})
	}
}
