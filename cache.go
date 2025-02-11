package cache

import (
	"errors"
	"github.com/adverax/containers/collections"
	"time"
)

type Cache[K comparable, V any] struct {
	keeper *keeper[K, V]
	Behavior[K, V]
	counter int64
}

func (that *Cache[K, V]) Close() {
	that.Behavior.Close()
}

func (that *Cache[K, V]) Set(k K, v V) {
	that.Assign(k, v, that.Behavior.Duration())
}

func (that *Cache[K, V]) Assign(k K, v V, d time.Duration) {
	that.keeper.Lock()
	defer that.keeper.Unlock()

	that.Behavior.Cleanup(that.keeper)

	oldItem := that.get(k)
	newItem := that.set(k, v, d)
	that.Behavior.Set(that.keeper, oldItem, newItem)
}

func (that *Cache[K, V]) set(k K, v V, d time.Duration) *Item[K, V] {
	var expiration int64
	if d > 0 {
		expiration = time.Now().Add(d).UnixNano()
	}

	that.counter++
	item := &Item[K, V]{
		id:         that.counter,
		key:        k,
		val:        v,
		expiration: expiration,
	}
	that.keeper.items[k] = item
	return item
}

func (that *Cache[K, V]) Add(k K, v V) error {
	return that.Append(k, v, that.Behavior.Duration())
}

func (that *Cache[K, V]) Append(k K, v V, d time.Duration) error {
	that.keeper.Lock()
	defer that.keeper.Unlock()

	that.Behavior.Cleanup(that.keeper)

	item := that.get(k)
	if item != nil {
		return ErrDuplicate
	}

	item = that.set(k, v, d)
	that.Behavior.Set(that.keeper, nil, item)
	return nil
}

func (that *Cache[K, V]) Replace(k K, v V, d time.Duration) error {
	that.keeper.Lock()
	defer that.keeper.Unlock()

	that.Behavior.Cleanup(that.keeper)

	oldItem := that.get(k)
	if oldItem == nil {
		return collections.ErrNoMatch
	}

	newItem := that.set(k, v, d)
	that.Behavior.Set(that.keeper, oldItem, newItem)
	return nil
}

func (that *Cache[K, V]) Get(k K) *Item[K, V] {
	that.keeper.Lock()
	defer that.keeper.Unlock()

	that.Behavior.Cleanup(that.keeper)

	item := that.get(k)

	if item != nil && item.IsExpired() {
		delete(that.keeper.items, k)
		that.Behavior.Set(that.keeper, item, nil)
		return nil
	}

	return item
}

func (that *Cache[K, V]) get(k K) *Item[K, V] {
	item, found := that.keeper.items[k]
	if !found {
		return nil
	}

	that.Behavior.Get(that.keeper, item)
	return item
}

func (that *Cache[K, T]) Delete(k K) {
	that.keeper.Lock()
	defer that.keeper.Unlock()

	that.Behavior.Cleanup(that.keeper)

	item := that.get(k)
	if item != nil {
		delete(that.keeper.items, item.key)
		that.Behavior.Set(that.keeper, item, nil)
	}
}

func (that *Cache[K, V]) ItemCount() int {
	that.keeper.Lock()
	defer that.keeper.Unlock()

	that.Behavior.Cleanup(that.keeper)

	return len(that.keeper.items)
}

func (that *Cache[K, V]) Items() map[K]*Item[K, V] {
	that.keeper.Lock()
	defer that.keeper.Unlock()

	that.Behavior.Cleanup(that.keeper)

	m := make(map[K]*Item[K, V], len(that.keeper.items))
	for k, item := range that.keeper.items {
		m[k] = item
	}

	return m
}

func (that *Cache[K, V]) Reset() {
	that.keeper.Lock()
	defer that.keeper.Unlock()

	that.Behavior.Reset(that.keeper)
	that.keeper.items = map[K]*Item[K, V]{}
}

func NewCache[K comparable, V any](
	behavior Behavior[K, V],
	index Index[K, V],
) *Cache[K, V] {
	if isZeroVal(behavior) {
		behavior = NewDefaultBehavior[K, V]()
	}

	if isZeroVal(index) {
		index = NewDefaultIndex[K, V]()
	}

	return &Cache[K, V]{
		keeper:   newKeeper(index),
		Behavior: behavior,
	}
}

var (
	ErrDuplicate = errors.New("duplicate key")
)
