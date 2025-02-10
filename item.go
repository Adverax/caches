package cache

import "time"

type Item[K comparable, V any] struct {
	id         int64
	key        K
	val        V
	expiration int64
	size       int64
}

func (that *Item[K, V]) ID() int64 {
	return that.id
}

func (that *Item[K, V]) Key() K {
	return that.key
}

func (that *Item[K, V]) Value() V {
	return that.val
}

func (that *Item[K, V]) Expiration() int64 {
	return that.expiration
}

func (that *Item[K, V]) IsExpired() bool {
	return that.expiration != 0 && that.IsExpiredEx(time.Now().UnixNano())
}

func (that *Item[K, V]) IsExpiredEx(now int64) bool {
	return now > that.expiration
}

func (that *Item[K, V]) Size() int64 {
	return that.size
}

func (that *Item[K, V]) SetSize(size int64) {
	that.size = size
}

func (that *Item[K, V]) SetExpiration(expiration int64) {
	that.expiration = expiration
}
