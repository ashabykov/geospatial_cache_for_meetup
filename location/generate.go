package location

import (
	"math"
	"math/rand"
	"time"

	"github.com/google/uuid"
)

func ts() Timestamp {
	return Timestamp(time.Now().Unix())
}

func name() Name {
	return Name(uuid.New().String())
}

func randDist(radius float64) float64 {
	return math.Sqrt(rand.Float64()) * radius
}

func Generate(center Location, radius float64) (Location, float64) {
	return pointAtDistance(center, randDist(radius), ts(), name())
}
