package fanout_write_client

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
		Near(target location.Location, radius float64, limit int) ([]location.Location, error)
		Set(target location.Location) error
	}

	Client struct {
		subscriber subscriber
		geospatial geospatial
	}
)

func (cl *Client) Near(target location.Location, radius float64, limit int) ([]location.Location, error) {
	return cl.geospatial.Near(target, radius, limit)
}

func (cl *Client) SubscribeOnUpdates(ctx context.Context) {
	results, err := cl.subscriber.Subscribe(ctx)
	if err != nil {

		fmt.Println("Client subscriber error:", err)

		return
	}

	for result := range results {

		if err = cl.geospatial.Set(result); err != nil {
			fmt.Println("Client subscriber set error:", err)
		}

		fmt.Println("Client geospatial set:", result)
	}
}

func New(sub subscriber, geospatial geospatial) *Client {
	return &Client{
		subscriber: sub,
		geospatial: geospatial,
	}
}
