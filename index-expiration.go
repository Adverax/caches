package cache

import (
	"github.com/adverax/containers/indices"
)

type expirationComparator[K comparable, V any] struct{}

func (that *expirationComparator[K, V]) Less(a, b Entry[K, V]) bool {
	return a.Expiration() < b.Expiration()
}

func (that *expirationComparator[K, V]) Greater(a, b Entry[K, V]) bool {
	return a.Expiration() > b.Expiration()
}

func (that *expirationComparator[K, V]) Equal(a, b Entry[K, V]) bool {
	return a.Expiration() == b.Expiration()
}

type IndexExpiration[K comparable, V any] struct {
	*indicies.Sorted[Entry[K, V]]
}

func NewExpirationIndex[K comparable, V any]() *IndexExpiration[K, V] {
	return &IndexExpiration[K, V]{
		Sorted: indicies.NewSorted[Entry[K, V]](
			&expirationComparator[K, V]{},
		),
	}
}
