package city_based_hasher

import "github.com/ashabykov/geospatial_cache_for_meetup/location"

type City struct {
	UUID   string
	Radius float64
	Center location.Location
}

type Hasher struct {
	cities []City
}

func (h *Hasher) Hash(loc location.Location) string {
	for _, city := range h.cities {
		if loc.CosineDistance(city.Center) <= city.Radius {
			return city.UUID
		}
	}
	return "unknown"
}
