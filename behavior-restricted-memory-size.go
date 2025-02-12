package cache

type BehaviorRestrictedMemorySize[K comparable, V any] struct {
	Behavior[K, V]
	size    int64
	maxSize int64
	sizeOf  func(item Entry[K, V]) int64
}

func (that *BehaviorRestrictedMemorySize[K, V]) Cleanup(keeper Keeper[K, V]) {
	that.Behavior.Cleanup(keeper)
	keeper.Index().Truncate(
		func(item Entry[K, V]) bool {
			if that.size <= that.maxSize {
				return false
			}
			that.size -= item.Size()
			keeper.Unset(item.Key())
			return true
		},
	)
}

func (that *BehaviorRestrictedMemorySize[K, V]) Set(keeper Keeper[K, V], oldEntry, newEntry Entry[K, V]) {
	if !isZeroVal(oldEntry) {
		that.size -= oldEntry.Size()
	}
	if !isZeroVal(newEntry) {
		newEntry.SetSize(that.sizeOf(newEntry))
		that.size += newEntry.Size()
	}
	that.Behavior.Set(keeper, oldEntry, newEntry)
}

func NewRestrictedMemorySizeBehavior[K comparable, V any](
	behavior Behavior[K, V],
	maxSize int64,
	sizeOf func(item Entry[K, V]) int64,
) Behavior[K, V] {
	if isZeroVal(behavior) {
		behavior = NewDefaultBehavior[K, V]()
	}

	if maxSize == 0 {
		maxSize = 1000000
	}

	return &BehaviorRestrictedMemorySize[K, V]{
		Behavior: behavior,
		maxSize:  maxSize,
		sizeOf:   sizeOf,
	}
}
