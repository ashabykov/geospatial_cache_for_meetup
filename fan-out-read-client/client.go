package fatout_read_client

import "github.com/ashabykov/geospatial_cache_for_meetup/location"

type (
	geospatial interface {
		Near(target location.Location, radius float64, limit int) ([]location.Location, error)
	}

	Client struct {
		geospatial geospatial
	}
)

func (cl *Client) Near(target location.Location, radius float64, limit int) ([]location.Location, error) {
	return cl.geospatial.Near(target, radius, limit)
}

func New(geospatial geospatial) *Client {
	return &Client{
		geospatial: geospatial,
	}
}
