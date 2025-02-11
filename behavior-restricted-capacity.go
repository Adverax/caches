package cache

type BehaviorRestrictedCapacity[K comparable, V any] struct {
	Behavior[K, V]
	capacity int
	index    Index[K, V]
}

func (that *BehaviorRestrictedCapacity[K, V]) Cleanup(keeper Keeper[K, V]) {
	that.Behavior.Cleanup(keeper)
	keeper.Index().Truncate(func(item Entry[K, V]) bool {
		if keeper.Length() <= that.capacity {
			return false
		}
		keeper.Unset(item.Key())
		return true
	})
}

type behaviorCapacityProlongation[K comparable, V any] struct {
	Behavior[K, V]
}

func (that *behaviorCapacityProlongation[K, V]) Get(keeper Keeper[K, V], entry Entry[K, V]) {
	that.Behavior.Get(keeper, entry)
	keeper.Index().Exclude(entry)
	keeper.Index().Include(entry)
}

func NewRestrictedCapacityBehavior[K comparable, V any](
	behavior Behavior[K, V],
	capacity int,
	prolongation bool,
) Behavior[K, V] {
	if isZeroVal(behavior) {
		behavior = NewDefaultBehavior[K, V]()
	}

	if capacity == 0 {
		capacity = 1000000
	}

	b := &BehaviorRestrictedCapacity[K, V]{
		Behavior: behavior,
		capacity: capacity,
	}

	if !prolongation {
		return b
	}

	return &behaviorCapacityProlongation[K, V]{
		Behavior: b,
	}
}
