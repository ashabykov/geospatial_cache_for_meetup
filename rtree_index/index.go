// based on impl: https://github.com/tidwall/rtree

package rtree_index

import (
	"github.com/tidwall/rtree"

	"github.com/ashabykov/geospatial_cache_for_meetup/location"
)

type Index struct {
	base rtree.RTreeG[string]
}

func NewIndex() *Index {
	return &Index{
		base: rtree.RTreeG[string]{},
	}
}

func (idx *Index) Size() int {
	return idx.base.Len()
}

func (idx *Index) Add(location location.Location) {
	idx.base.Insert(
		location.List(),
		location.List(),
		location.Key(),
	)
}

func (idx *Index) Remove(location location.Location) {
	idx.base.Delete(
		location.List(),
		location.List(),
		location.Key(),
	)
}

func (idx *Index) Nearby(
	target location.Location,
	radius float64,
	limit int,
) []location.Name {

	result := make([]location.Name, 0, limit)

	idx.base.Nearby(
		CosineDistance(target, nil),
		func(min, max [2]float64, name string, dist float64) bool {
			// filter by limit
			if len(result) == limit {
				// if we reached to the limit
				// we do must stop to iterate
				return false
			}

			// filter by radius
			if dist > radius {
				// we must check
				// until reach to the limit
				// do not stop to iterate
				return true
			}

			result = append(
				result,
				location.Name(name),
			)
			return true
		},
	)
	return result
}

func CosineDistance(
	targ location.Location,
	itemDist func(min, max [2]float64, data string) float64,
) (dist func(min, max [2]float64, data string, item bool) float64) {
	return func(min, max [2]float64, data string, item bool) (dist float64) {
		if item && itemDist != nil {
			return itemDist(min, max, data)
		}
		return targ.CosineDistance(
			location.Location{
				Lon: location.Longitude(max[0]),
				Lat: location.Latitude(max[1]),
			},
		)
	}
}
