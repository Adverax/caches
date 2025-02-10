package cache

import "time"

type BehaviorWithAutoCleanup[K comparable, V any] struct {
	Behavior[K, V]
	interval time.Duration
	done     chan struct{}
}

func (that *BehaviorWithAutoCleanup[K, V]) Cleanup(keeper *Keeper[K, V]) {
	// nothing
}

func (that *BehaviorWithAutoCleanup[K, V]) Start(keeper Keeper[K, V]) {
	go func() {
		ticker := time.NewTicker(that.interval)
		for {
			select {
			case <-that.done:
				return
			case <-ticker.C:
				that.autoCleanup(keeper)
			}
		}
	}()
}

func (that *BehaviorWithAutoCleanup[K, V]) autoCleanup(keeper Keeper[K, V]) {
	keeper.Lock()
	defer keeper.Unlock()

	that.Behavior.Cleanup(keeper)
}

func (that *BehaviorWithAutoCleanup[K, V]) Close() {
	close(that.done)
}

func NewAutoCleanupBehavior[K comparable, V any](
	behavior Behavior[K, V],
	interval time.Duration,
) *BehaviorWithAutoCleanup[K, V] {
	if isZeroVal(behavior) {
		behavior = NewDefaultBehavior[K, V]()
	}
	if interval == 0 {
		interval = time.Second
	}

	return &BehaviorWithAutoCleanup[K, V]{
		Behavior: behavior,
		interval: interval,
		done:     make(chan struct{}),
	}
}
