package cache

type IndexDefault[K comparable, V any] struct{}

func (that *IndexDefault[K, V]) Reset() {}

func (that *IndexDefault[K, V]) Truncate(iterator func(entry Entry[K, V]) bool) {}

func (that *IndexDefault[K, V]) Append(entry Entry[K, V]) {}

func (that *IndexDefault[K, V]) Remove(entry Entry[K, V]) {}

func NewDefaultIndex[K comparable, V any]() *IndexDefault[K, V] {
	return &IndexDefault[K, V]{}
}
