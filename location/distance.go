package location

import (
	"math"

	gener "gopkg.in/gujarats/GenerateLocation.v1"
)

const (
	degToRad = math.Pi / 180.0
)

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
func deg2Rad(degree float64) float64 { return degree * degToRad }

// rad2Deg converts from radians to degree measure.
func rad2Deg(rad float64) float64 { return rad / degToRad }

func pointAtDistance(self Location, radius float64) (Location, float64) {
	// Convert Degrees to radians
	tmp := gener.New(self.Lat.Float64(), self.Lon.Float64())
	newLoc := tmp.GenerateLocation(radius/1000, radius/1000)
	other := Location{
		Name: name(),
		Ts:   Now(),
		Lat:  Latitude(newLoc[0].Lat),
		Lon:  Longitude(newLoc[0].Lon),
		TTL:  self.TTL,
	}
	return other, other.CosineDistance(self)
}
