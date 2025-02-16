package location

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestGenerate(t *testing.T) {
	tests := []struct {
		name   string
		center Location
		radius float64
	}{
		{
			name: "",
			center: Location{
				Name: "target",
				Lat:  43.244555,
				Lon:  76.940012,
				Ts:   Timestamp(time.Now().Unix()),
				TTL:  10 * time.Minute,
			},
			radius: 5000.0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newLoc, dist := Generate(tt.center, tt.radius)
			assert.NotNil(t, newLoc)
			assert.Greater(t, tt.radius, dist)
		})
	}
}
