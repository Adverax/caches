package cache

import "sort"

type IndexSerial[K comparable, V any] struct {
	items []Entry[K, V]
}

func (that *IndexSerial[K, V]) Len() int {
	return len(that.items)
}

func (that *IndexSerial[K, V]) Swap(i, j int) {
	that.items[i], that.items[j] = that.items[j], that.items[i]
}

func (that *IndexSerial[K, V]) Less(i, j int) bool {
	return that.items[i].ID() < that.items[j].ID()
}

func (that *IndexSerial[K, V]) Flush() {}

func (that *IndexSerial[K, V]) Truncate(iterator func(entry Entry[K, V]) bool) {
	for len(that.items) > 0 {
		item := that.items[0]
		if iterator(item) {
			that.items = that.items[1:]
			continue
		}
		return
	}
}

func (that *IndexSerial[K, V]) indexOf(entry Entry[K, V]) int {
	low := 0
	high := len(that.items) - 1

	for low <= high {
		mid := low + (high-low)/2

		if that.items[mid].ID() < entry.ID() {
			low = mid + 1
		} else if that.items[mid].ID() > entry.ID() {
			high = mid - 1
		} else {
			if that.items[mid] == entry {
				return mid
			}

			right := mid + 1
			for right <= high && that.items[right].ID() == entry.ID() {
				if that.items[right] == entry {
					return right
				}
				right++
			}

			left := mid - 1
			for left >= low && that.items[left].ID() == entry.ID() {
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

func (that *IndexSerial[K, V]) Retract(entry Entry[K, V]) {
	index := that.indexOf(entry)
	if index != -1 {
		that.items = append(that.items[:index], that.items[index+1:]...)
	}
}

func (that *IndexSerial[K, V]) Assert(entry Entry[K, V]) {
	index := sort.Search(
		len(that.items),
		func(i int) bool {
			return that.items[i].ID() >= entry.ID()
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

func NewSerialIndex[K comparable, V any]() Index[K, V] {
	return &IndexSerial[K, V]{}
}
