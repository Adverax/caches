package cache

import "time"

type DefaultBehavior[K comparable, V any] struct {
	index Index[K, V]
}

func (that *DefaultBehavior[K, V]) Duration() time.Duration {
	return 0
}

func (that *DefaultBehavior[K, V]) Close() {}

func (that *DefaultBehavior[K, V]) Flush(keeper Keeper[K, V]) {
	keeper.Index().Flush()
}

func (that *DefaultBehavior[K, V]) Get(keeper Keeper[K, V], entry Entry[K, V]) {}

func (that *DefaultBehavior[K, V]) Set(keeper Keeper[K, V], oldEntry, newEntry Entry[K, V]) {
	if !isZeroVal(oldEntry) {
		keeper.Index().Retract(oldEntry)
	}
	if !isZeroVal(newEntry) {
		keeper.Index().Assert(newEntry)
	}
}

func (that *DefaultBehavior[K, V]) Cleanup(keeper Keeper[K, V]) {}

func NewDefaultBehavior[K comparable, V any]() *DefaultBehavior[K, V] {
	return &DefaultBehavior[K, V]{}
}
