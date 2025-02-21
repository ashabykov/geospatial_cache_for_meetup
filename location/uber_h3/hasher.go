package uber_h3

import (
	"github.com/ashabykov/geospatial_cache_for_meetup/location"
	"github.com/uber/h3-go/v4"
)

type Hasher struct {
	resolution int
}

func (h *Hasher) Hash(loc location.Location) string {
	return h3.LatLngToCell(
		h3.NewLatLng(loc.Lat.Float64(), loc.Lon.Float64()),
		h.resolution,
	).String()
}
