package sorted_set

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/ashabykov/geospatial_cache_for_meetup/location"
)

func TestIndex_Remove(t *testing.T) {
	tests := []struct {
		name      string
		from, to  location.Timestamp
		target    location.Location
		locations []location.Location
	}{
		{
			name: "remove from index",
			from: location.Timestamp(time.Now().UTC().Add(-6 * time.Minute).Unix()),
			to:   location.Timestamp(time.Now().UTC().Add(1 * time.Minute).Unix()),
			target: location.NewLocation(
				"awesome-3",
				location.Timestamp(
					time.Now().UTC().Add(-5*time.Minute).Unix(),
				),
				45.786877,
				47.679879,
			),
			locations: []location.Location{
				location.NewLocation(
					"awesome-1",
					location.Timestamp(
						time.Now().UTC().Add(-15*time.Minute).Unix(),
					),
					45.68878,
					47.57867,
				),
				location.NewLocation(
					"awesome-2",
					location.Timestamp(
						time.Now().UTC().Add(-10*time.Minute).Unix(),
					),
					45.467467,
					47.456656,
				),
				location.NewLocation(
					"awesome-3",
					location.Timestamp(
						time.Now().UTC().Add(-5*time.Minute).Unix(),
					),
					45.786877,
					47.679879,
				),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			idx := New()

			for _, loc := range tt.locations {
				idx.Add(loc)
			}

			idx.Remove(tt.target)

			assert.Equal(t, 0, len(idx.Read(tt.from, tt.to)))
		})
	}
}

func TestIndex_Read(t *testing.T) {
	tests := []struct {
		name      string
		from, to  location.Timestamp
		locations []location.Location
		expected  []location.Name
	}{
		{
			name: "must return last two location's names",
			from: location.Timestamp(time.Now().UTC().Add(-6 * time.Minute).Unix()),
			to:   location.Timestamp(time.Now().UTC().Add(1 * time.Minute).Unix()),
			locations: []location.Location{
				location.NewLocation(
					"awesome-1",
					location.Timestamp(
						time.Now().UTC().Add(-15*time.Minute).Unix(),
					),
					45.68878,
					47.57867,
				),
				location.NewLocation(
					"awesome-2",
					location.Timestamp(
						time.Now().UTC().Add(-10*time.Minute).Unix(),
					),
					45.467467,
					47.456656,
				),
				location.NewLocation(
					"awesome-3",
					location.Timestamp(
						time.Now().UTC().Add(-5*time.Minute).Unix(),
					),
					45.786877,
					47.679879,
				),
				location.NewLocation(
					"awesome-4",
					location.Timestamp(
						time.Now().UTC().Unix(),
					),
					45.786877,
					47.679879,
				),
			},
			expected: []location.Name{
				"awesome-3",
				"awesome-4",
			},
		},
		{
			name: "must return fist two location's names",
			from: location.Timestamp(0),
			to:   location.Timestamp(time.Now().UTC().Add(-10 * time.Minute).Unix()),
			locations: []location.Location{
				location.NewLocation(
					"awesome-1",
					location.Timestamp(
						time.Now().UTC().Add(-15*time.Minute).Unix(),
					),
					45.68878,
					47.57867,
				),
				location.NewLocation(
					"awesome-2",
					location.Timestamp(
						time.Now().UTC().Add(-10*time.Minute).Unix(),
					),
					45.467467,
					47.456656,
				),
				location.NewLocation(
					"awesome-3",
					location.Timestamp(
						time.Now().UTC().Add(-5*time.Minute).Unix(),
					),
					45.786877,
					47.679879,
				),
				location.NewLocation(
					"awesome-4",
					location.Timestamp(
						time.Now().UTC().Unix(),
					),
					45.786877,
					47.679879,
				),
			},
			expected: []location.Name{
				"awesome-1",
				"awesome-2",
			},
		},
		{
			name: "all locations are expired",
			from: location.Timestamp(time.Now().UTC().Add(-4 * time.Minute).Unix()),
			to:   location.Timestamp(time.Now().UTC().Unix()),
			locations: []location.Location{
				location.NewLocation(
					"awesome-1",
					location.Timestamp(
						time.Now().UTC().Add(-15*time.Minute).Unix(),
					),
					45.68878,
					47.57867,
				),
				location.NewLocation(
					"awesome-2",
					location.Timestamp(
						time.Now().UTC().Add(-10*time.Minute).Unix(),
					),
					45.467467,
					47.456656,
				),
				location.NewLocation(
					"awesome-3",
					location.Timestamp(
						time.Now().UTC().Add(-5*time.Minute).Unix(),
					),
					45.786877,
					47.679879,
				),
			},
			expected: []location.Name{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			idx := New()

			for _, loc := range tt.locations {
				idx.Add(loc)
			}

			assert.Equal(t, tt.expected, idx.Read(tt.from, tt.to))
		})
	}
}
