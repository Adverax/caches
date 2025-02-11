package cache

type IndexDefault[K comparable, V any] struct{}

func (that *IndexDefault[K, V]) Reset() {}

func (that *IndexDefault[K, V]) Truncate(iterator func(entry Entry[K, V]) bool) {}

func (that *IndexDefault[K, V]) Include(entry Entry[K, V]) bool { return false }

func (that *IndexDefault[K, V]) Exclude(entry Entry[K, V]) bool { return false }

func NewDefaultIndex[K comparable, V any]() *IndexDefault[K, V] {
	return &IndexDefault[K, V]{}
}
