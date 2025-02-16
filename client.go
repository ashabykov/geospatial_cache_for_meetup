package geospatial_cache_for_meetup

import (
	"context"
	"fmt"

	"github.com/ashabykov/geospatial_cache_for_meetup/location"
)

type (
	subscriber interface {
		Subscribe(ctx context.Context) (<-chan location.Location, error)
	}

	geospatial interface {
		Near(target location.Location, radius float64, limit int) []location.Location
		Set(target location.Location)
	}

	Client struct {
		subscriber subscriber
		geospatial geospatial
	}
)

func (cl *Client) Near(target location.Location, radius float64, limit int) []location.Location {
	return cl.geospatial.Near(target, radius, limit)
}

func (cl *Client) Subscribe(ctx context.Context) {
	results, err := cl.subscriber.Subscribe(ctx)
	if err != nil {

		fmt.Println("Client Subscriber error:", err)

		return
	}

	for result := range results {
		cl.geospatial.Set(result)
	}
}

func New(sub subscriber, geospatial geospatial) *Client {
	return &Client{
		subscriber: sub,
		geospatial: geospatial,
	}
}
