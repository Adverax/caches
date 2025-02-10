package cache

import "sort"

type IndexExpiration[K comparable, V any] struct {
	items []Entry[K, V]
}

func (that *IndexExpiration[K, V]) Len() int {
	return len(that.items)
}

func (that *IndexExpiration[K, V]) Swap(i, j int) {
	that.items[i], that.items[j] = that.items[j], that.items[i]
}

func (that *IndexExpiration[K, V]) Less(i, j int) bool {
	return that.items[i].Expiration() < that.items[j].Expiration()
}

func (that *IndexExpiration[K, V]) Reset() {
	that.items = nil
}

func (that *IndexExpiration[K, V]) Truncate(iterator func(entry Entry[K, V]) bool) {
	for len(that.items) > 0 {
		item := that.items[0]
		if iterator(item) {
			that.items = that.items[1:]
			continue
		}
		return
	}
}

func (that *IndexExpiration[K, V]) indexOf(entry Entry[K, V]) int {
	low := 0
	high := len(that.items) - 1

	expiration := entry.Expiration()
	for low <= high {
		mid := low + (high-low)/2

		if that.items[mid].Expiration() < expiration {
			low = mid + 1
		} else if that.items[mid].Expiration() > expiration {
			high = mid - 1
		} else {
			if that.items[mid] == entry {
				return mid
			}

			right := mid + 1
			for right <= high && that.items[right].Expiration() == expiration {
				if that.items[right] == entry {
					return right
				}
				right++
			}

			left := mid - 1
			for left >= low && that.items[left].Expiration() == expiration {
				if that.items[left] == entry {
					return left
				}
				left--
			}

			return -1
		}
	}

	return -1
}

func (that *IndexExpiration[K, V]) Remove(entry Entry[K, V]) {
	index := that.indexOf(entry)
	if index != -1 {
		that.items = append(that.items[:index], that.items[index+1:]...)
	}
}

func (that *IndexExpiration[K, V]) Append(entry Entry[K, V]) {
	index := sort.Search(
		len(that.items),
		func(i int) bool {
			return that.items[i].Expiration() >= entry.Expiration()
		},
	)

	if index == len(that.items) {
		that.items = append(that.items, entry)
	} else {
		that.items = append(that.items, nil)
		copy(that.items[index+1:], that.items[index:])
		that.items[index] = entry
	}
}

func NewExpirationIndex[K comparable, V any]() *IndexExpiration[K, V] {
	return &IndexExpiration[K, V]{
		items: make([]Entry[K, V], 0),
	}
}
