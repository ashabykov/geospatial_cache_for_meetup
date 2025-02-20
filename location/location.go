package location

import (
	"sort"
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

func (l Location) GeoHash() string {
	return h3.LatLngToCell(
		h3.NewLatLng(l.Lat.Float64(), l.Lon.Float64()), 5,
	).String()
}

type HexDistance struct {
	Hex      string
	Distance float64
}

func (l Location) NearHex(resolution int) []string {
	latLng := h3.NewLatLng(l.Lat.Float64(), l.Lon.Float64())
	originCell := h3.LatLngToCell(latLng, resolution)
	originCellLatLng := originCell.LatLng()

	hexDistances := make([]HexDistance, 0, 7)
	hexDistances = append(hexDistances, HexDistance{
		Hex: originCell.String(),
		Distance: l.CosineDistance(Location{
			Lat: Latitude(originCellLatLng.Lat),
			Lon: Longitude(originCellLatLng.Lng),
		}),
	})

	directedEdges := originCell.DirectedEdges()
	for i := range directedEdges {
		cellLatLng := directedEdges[i].Destination().LatLng()
		distance := l.CosineDistance(Location{
			Lat: Latitude(cellLatLng.Lat),
			Lon: Longitude(cellLatLng.Lng),
		})
		hexDistances = append(hexDistances, HexDistance{
			Hex:      directedEdges[i].Destination().String(),
			Distance: distance,
		})
	}

	sort.Slice(hexDistances, func(i, j int) bool {
		return hexDistances[i].Distance < hexDistances[j].Distance
	})

	result := make([]string, len(hexDistances))
	for i := 0; i < len(hexDistances); i++ {
		result[i] = hexDistances[i].Hex
	}

	return result
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
