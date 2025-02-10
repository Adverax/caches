package cache

import "time"

type Entry[K comparable, V any] interface {
	// ID returns the ID of the item.
	ID() int64
	// Key returns the key of the item.
	Key() K
	// Value returns the value of the item.
	Value() V
	// Expiration returns the expiration of the item.
	Expiration() int64
	// SetExpiration sets the expiration of the item.
	SetExpiration(expiration int64)
	// IsExpired returns true if the item is expired.
	IsExpired() bool
	// IsExpiredEx returns true if the item is expired at the given time.
	IsExpiredEx(now int64) bool
	// Size returns the size of the item.
	Size() int64
	// SetSize sets the size of the item.
	SetSize(size int64)
}

type Keeper[K comparable, V any] interface {
	// Lock locks the cache.
	Lock()
	// Unlock unlocks the cache.
	Unlock()
	// Unset removes an item from the cache.
	Unset(k K)
	// Length returns the number of items in the cache.
	Length() int
	// Index returns the index of the cache.
	Index() Index[K, V]
}

type Index[K comparable, V any] interface {
	// Truncate removes items from the index until the iterator returns false.
	Truncate(iterator func(entry Entry[K, V]) bool)
	// Assert adds an item to the index.
	Assert(entry Entry[K, V])
	// Retract removes an item from the index.
	Retract(entry Entry[K, V])
	// Flush clears the index.
	Flush()
}

type Behavior[K comparable, V any] interface {
	Duration() time.Duration
	Flush(keeper Keeper[K, V])
	Get(keeper Keeper[K, V], entry Entry[K, V])
	Set(keeper Keeper[K, V], oldEntry, newEntry Entry[K, V])
	Cleanup(keeper Keeper[K, V])
	Close()
}
