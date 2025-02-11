package cache

import (
	"github.com/adverax/containers/indices"
)

type serialComparator[K comparable, V any] struct{}

func (that *serialComparator[K, V]) Less(a, b Entry[K, V]) bool {
	return a.ID() < b.ID()
}

func (that *serialComparator[K, V]) Greater(a, b Entry[K, V]) bool {
	return a.ID() > b.ID()
}

func (that *serialComparator[K, V]) Equal(a, b Entry[K, V]) bool {
	return a.ID() == b.ID()
}

type IndexSerial[K comparable, V any] struct {
	*indicies.Sorted[Entry[K, V]]
}

func NewSerialIndex[K comparable, V any]() *IndexSerial[K, V] {
	return &IndexSerial[K, V]{
		Sorted: indicies.NewSorted[Entry[K, V]](
			&serialComparator[K, V]{},
		),
	}
}
