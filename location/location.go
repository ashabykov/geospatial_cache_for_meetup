package location

import (
	"math"
	"time"
)

const (
	angle = math.Pi / 180
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

func (l Location) EuclideanDistance(other Location) float64 {
	var (
		X1 = l.Lon.Float64()
		X2 = other.Lon.Float64()
		Y1 = l.Lat.Float64()
		Y2 = other.Lat.Float64()
	)
	return math.Sqrt((X1-X2)*(X1-X2) + (Y1-Y2)*(Y1-Y2))
}

func (l Location) CosineDistance(other Location) float64 {

	lon1, lat1 := l.Lon.Float64(), l.Lat.Float64()
	lon2, lat2 := other.Lon.Float64(), other.Lat.Float64()
	theta := lon1 - lon2
	dist := math.Sin(deg2Rad(lat1))*math.Sin(deg2Rad(lat2)) + math.Cos(deg2Rad(lat1))*math.Cos(deg2Rad(lat2))*math.Cos(deg2Rad(theta))
	dist = math.Acos(dist)
	dist = rad2Deg(dist)
	meters := dist * 60 * 1.1515 * 1.609344 * 1000

	if math.IsNaN(meters) {
		return 0
	}
	return meters
}

// Deg2Rad converts from degree measure to radiance.
func deg2Rad(degree float64) float64 { return degree * angle }

// rad2Deg converts from radians to degree measure.
func rad2Deg(rad float64) float64 { return rad / angle }

func NewNeighbour(loc Location, distance float64) Neighbour {
	return Neighbour{
		Location: loc,
		Distance: distance,
	}
}

func (n Neighbours) Len() int { return len(n) }
