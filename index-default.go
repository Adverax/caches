package cache

type IndexDefault[K comparable, V any] struct{}

func (that *IndexDefault[K, V]) Flush() {}

func (that *IndexDefault[K, V]) Truncate(iterator func(entry Entry[K, V]) bool) {}

func (that *IndexDefault[K, V]) Assert(entry Entry[K, V]) {}

func (that *IndexDefault[K, V]) Retract(entry Entry[K, V]) {}

func NewDefaultIndex[K comparable, V any]() *IndexDefault[K, V] {
	return &IndexDefault[K, V]{}
}
