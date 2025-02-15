// based on impl: https://github.com/wangjia184/sortedset

package sorted_set

import (
	"sync"

	"github.com/wangjia184/sortedset"

	"github.com/ashabykov/geospatial_cache_for_meetup/location"
)

type Index struct {
	mu sync.Mutex

	base *sortedset.SortedSet
}

func New() *Index {
	return &Index{
		base: sortedset.New(),
	}
}

func (index *Index) Add(location location.Location) {

	index.mu.Lock()

	defer index.mu.Unlock()

	index.base.AddOrUpdate(
		location.Key(),
		sortedset.SCORE(location.Ts.Int64()),
		location.Name,
	)
}

func (index *Index) Read(from, to location.Timestamp) []location.Name {
	items := index.base.GetByScoreRange(
		sortedset.SCORE(from.Int64()),
		sortedset.SCORE(to.Int64()),
		nil,
	)

	ret := make([]location.Name, 0, len(items))
	for _, item := range items {
		ret = append(ret, item.Value.(location.Name))
	}
	return ret
}

func (index *Index) Remove(location location.Location) {
	index.mu.Lock()

	defer index.mu.Unlock()

	index.base.Remove(location.Key())
}

func (index *Index) Len() int {
	return index.base.GetCount()
}
