package location

import (
	"math"
	"math/rand"

	"github.com/google/uuid"
)

func name() Name {
	return Name(uuid.New().String())
}

func randDist(radius float64) float64 {
	return math.Sqrt(rand.Float64()) * radius
}

func Generate(center Location, radius float64) (Location, float64) {
	return pointAtDistance(center, randDist(radius))
}
