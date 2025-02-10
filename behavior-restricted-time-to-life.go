package cache

import "time"

type BehaviorRestrictedTimeToLife[K comparable, V any] struct {
	Behavior[K, V]
	duration time.Duration
}

func (that *BehaviorRestrictedTimeToLife[K, V]) Duration() time.Duration {
	return that.duration
}

func (that *BehaviorRestrictedTimeToLife[K, V]) Cleanup(keeper Keeper[K, V]) {
	that.Behavior.Cleanup(keeper)
	now := time.Now().UnixNano()
	keeper.Index().Truncate(
		func(entry Entry[K, V]) bool {
			if !entry.IsExpiredEx(now) {
				return false
			}
			keeper.Unset(entry.Key())
			return true
		},
	)
}

type behaviorExpirationProlongation[K comparable, V any] struct {
	Behavior[K, V]
	prolongation time.Duration
}

func (that *behaviorExpirationProlongation[K, V]) Get(keeper Keeper[K, V], entry Entry[K, V]) {
	that.Behavior.Get(keeper, entry)
	index := keeper.Index()
	index.Remove(entry)
	entry.SetExpiration(time.Now().Add(that.prolongation).UnixNano())
	index.Append(entry)
}

func NewRestrictedTimeToLifeBehavior[K comparable, V any](
	behavior Behavior[K, V],
	duration time.Duration,
	prolongation bool,
) Behavior[K, V] {
	if isZeroVal(behavior) {
		behavior = NewDefaultBehavior[K, V]()
	}

	if duration == 0 {
		duration = 5 * time.Minute
	}

	b := &BehaviorRestrictedTimeToLife[K, V]{
		Behavior: behavior,
		duration: duration,
	}

	if !prolongation {
		return b
	}

	return &behaviorExpirationProlongation[K, V]{
		Behavior:     b,
		prolongation: duration,
	}
}
