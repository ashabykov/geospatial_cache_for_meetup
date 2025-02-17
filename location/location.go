package location

import "time"

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

func (lat Latitude) Float64() float64 {
	return float64(lat)
}

func (l Location) Key() string {
	return string(l.Name)
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
