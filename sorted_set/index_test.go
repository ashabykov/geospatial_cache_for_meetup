package sorted_set

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/ashabykov/geospatial_cache_for_meetup/location"
)

func TestIndex_Read(t *testing.T) {
	tests := []struct {
		name      string
		from, to  location.Timestamp
		locations []location.Location
		expected  []location.Name
	}{
		{
			name: "must return last two location's names",
			from: location.Timestamp(time.Now().Add(-6 * time.Minute).Unix()),
			to:   location.Timestamp(time.Now().Add(1 * time.Minute).Unix()),
			locations: []location.Location{
				location.NewLocation(
					"awesome-1",
					location.Timestamp(
						time.Now().Add(-15*time.Minute).Unix(),
					),
					45.68878,
					47.57867,
				),
				location.NewLocation(
					"awesome-2",
					location.Timestamp(
						time.Now().Add(-10*time.Minute).Unix(),
					),
					45.467467,
					47.456656,
				),
				location.NewLocation(
					"awesome-3",
					location.Timestamp(
						time.Now().Add(-5*time.Minute).Unix(),
					),
					45.786877,
					47.679879,
				),
				location.NewLocation(
					"awesome-4",
					location.Timestamp(
						time.Now().Unix(),
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
			name: "all locations are expired",
			from: location.Timestamp(time.Now().Add(-4 * time.Minute).Unix()),
			to:   location.Timestamp(time.Now().Unix()),
			locations: []location.Location{
				location.NewLocation(
					"awesome-1",
					location.Timestamp(
						time.Now().Add(-15*time.Minute).Unix(),
					),
					45.68878,
					47.57867,
				),
				location.NewLocation(
					"awesome-2",
					location.Timestamp(
						time.Now().Add(-10*time.Minute).Unix(),
					),
					45.467467,
					47.456656,
				),
				location.NewLocation(
					"awesome-3",
					location.Timestamp(
						time.Now().Add(-5*time.Minute).Unix(),
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

			index := New()

			for _, loc := range tt.locations {
				index.Add(loc)
			}

			assert.Equal(t, tt.expected, index.Read(tt.from, tt.to))
		})
	}
}
