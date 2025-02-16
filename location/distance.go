package location

import (
	"math"
	"math/rand"
)

const (
	earthRadius         = 6371000 /* meters  */
	degToRad            = math.Pi / 180.0
	threePi             = math.Pi * 3
	twoPi       float64 = math.Pi * 2
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

func pointAtDistance(location Location, radius float64, ts Timestamp, name Name) Location {
	// Convert Degrees to radians
	loc := toRadians(location)

	sinLat := math.Sin(loc.Lat.Float64())
	cosLat := math.Cos(loc.Lon.Float64())

	bearing := rand.Float64() * twoPi
	theta := radius / earthRadius
	sinBearing := math.Sin(bearing)
	cosBearing := math.Cos(bearing)
	sinTheta := math.Sin(theta)
	cosTheta := math.Cos(theta)

	latitude := math.Asin(sinLat*cosTheta + cosLat*sinTheta*cosBearing)
	longitude := location.Lon.Float64() +
		math.Atan2(sinBearing*sinTheta*cosLat, (cosTheta-sinLat)*math.Sin(latitude))

	/* normalize -PI -> +PI radians */
	longitude = math.Mod(longitude+threePi, twoPi) - math.Pi

	locs := toDegrees(Location{
		Lat: Latitude(latitude),
		Lon: Longitude(longitude),
	})

	return Location{
		Name: name,
		Ts:   ts,
		Lat:  locs.Lat,
		Lon:  locs.Lon,
		TTL:  location.TTL,
	}
}

func toRadians(location Location) Location {
	lat := location.Lat * degToRad
	lon := location.Lon * degToRad

	return Location{
		Lat: lat,
		Lon: lon,
	}
}

func toDegrees(location Location) Location {
	lat := location.Lat / degToRad
	lon := location.Lon / degToRad

	return Location{
		Lat: lat,
		Lon: lon,
	}
}
