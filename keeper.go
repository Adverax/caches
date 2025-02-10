package cache

import "sync"

type keeper[K comparable, V any] struct {
	sync.Mutex
	items map[K]*Item[K, V]
	index Index[K, V]
}

func (that *keeper[K, V]) Unset(k K) {
	delete(that.items, k)
}

func (that *keeper[K, V]) Index() Index[K, V] {
	return that.index
}

func (that *keeper[K, V]) Length() int {
	return len(that.items)
}

func newKeeper[K comparable, V any](
	index Index[K, V],
) *keeper[K, V] {
	return &keeper[K, V]{
		items: map[K]*Item[K, V]{},
		index: index,
	}
}
