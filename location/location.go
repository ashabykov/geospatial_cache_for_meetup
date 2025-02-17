package location

import (
	"time"

	"github.com/uber/h3-go/v4"
)

type (
	Name      string
	Longitude float64
	Latitude  float64
	Timestamp int64

	Location struct {
		Name Name          `json:"name"`
		Ts   Timestamp     `json:"ts"`
		TTL  time.Duration `json:"ttl"`
		Lon  Longitude     `json:"lon"`
		Lat  Latitude      `json:"lat"`
	}

	Neighbour struct {
		Location Location
		Distance float64
	}

	Neighbours []Neighbour
)

func Now() Timestamp {
	return Timestamp(time.Now().UTC().Unix())
}

func NewLocation(name Name, ts Timestamp, lon Longitude, lat Latitude) Location {
	return Location{
		Name: name,
		Ts:   ts,
		Lon:  lon,
		Lat:  lat,
	}
}

func (name Name) String() string {
	return string(name)
}

func (lng Longitude) Float64() float64 {
	return float64(lng)
}

func (ts Timestamp) Int64() int64 {
	return int64(ts)
}

func (ts Timestamp) Float64() float64 {
	return float64(ts)
}

func (lat Latitude) Float64() float64 {
	return float64(lat)
}

func (l Location) Key() string {
	return string(l.Name)
}

func (l Location) ShardKey() string {
	return h3.LatLngToCell(
		h3.NewLatLng(l.Lat.Float64(), l.Lon.Float64()), 5,
	).String()
}

func (l Location) ShardKeys() string {
	LatLng := h3.NewLatLng(l.Lat.Float64(), l.Lon.Float64())
	h3.
}

func (l Location) List() [2]float64 {
	return [2]float64{l.Lon.Float64(), l.Lat.Float64()}
}

func NewNeighbour(loc Location, distance float64) Neighbour {
	return Neighbour{
		Location: loc,
		Distance: distance,
	}
}

func (n Neighbours) Len() int { return len(n) }
