package cache

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type CacheShould struct {
	suite.Suite
}

func TestCache(t *testing.T) {
	suite.Run(t, new(CacheShould))
}

func (that *CacheShould) TestPersistent() {
	c := NewCache[string, string](nil, nil)
	c.Set("hello", "word")
	item := c.Get("hello")
	that.NotNil(item)
	that.Equal("word", item.val)
}

func (that *CacheShould) TestWithNotExpired_MustBeAccessable() {
	c := NewCache[string, string](
		NewRestrictedTimeToLifeBehavior[string, string](nil, time.Hour, false),
		NewExpirationIndex[string, string](),
	)
	c.Set("hello1", "word1")
	time.Sleep(time.Microsecond)
	c.Set("hello2", "word2")
	time.Sleep(time.Microsecond)
	c.Set("hello3", "word3")
	time.Sleep(1000 * time.Millisecond)
	item := c.Get("hello1")
	that.NotNil(item)
	item = c.Get("hello2")
	that.NotNil(item)
	item = c.Get("hello3")
	that.NotNil(item)
}

func (that *CacheShould) TestWithExpired_MustBeNotAccessable() {
	c := NewCache[string, string](
		NewRestrictedTimeToLifeBehavior[string, string](nil, 10*time.Millisecond, false),
		NewExpirationIndex[string, string](),
	)
	c.Set("hello1", "word1")
	time.Sleep(time.Microsecond)
	c.Set("hello2", "word2")
	time.Sleep(time.Microsecond)
	c.Set("hello3", "word3")
	time.Sleep(time.Microsecond)
	c.Assign("hello4", "word4", time.Hour)
	time.Sleep(20 * time.Millisecond)
	item := c.Get("hello1")
	that.Nil(item)
	item = c.Get("hello2")
	that.Nil(item)
	item = c.Get("hello3")
	that.Nil(item)
	item = c.Get("hello4")
	that.NotNil(item)
}

func (that *CacheShould) TestWithProlongation_MustBeAccessable() {
	c := NewCache[string, string](
		NewRestrictedTimeToLifeBehavior[string, string](nil, 50*time.Millisecond, true),
		NewExpirationIndex[string, string](),
	)

	c.Set("hello1", "word1")

	for i := 0; i < 10; i++ {
		c.Get("hello1")
		time.Sleep(20 * time.Millisecond)
	}

	item := c.Get("hello1")
	that.NotNil(item)
	expiration := time.Unix(0, item.expiration)
	d := expiration.Sub(time.Now())
	that.True(d > time.Millisecond*20)
}

func (that *CacheShould) TestWithCapacity_MustRemoveItems() {
	c := NewCache[string, string](
		NewRestrictedCapacityBehavior(
			NewRestrictedTimeToLifeBehavior[string, string](nil, time.Hour, false),
			2,
			false,
		),
		NewExpirationIndex[string, string](),
	)

	c.Set("hello1", "word1")
	time.Sleep(10 * time.Millisecond)
	c.Set("hello2", "word2")
	time.Sleep(10 * time.Millisecond)
	c.Set("hello3", "word3")
	that.Nil(c.Get("hello1"))
	that.NotNil(c.Get("hello2"))
	that.NotNil(c.Get("hello3"))
}

func (that *CacheShould) TestWithMaxSize_MustRemoveItems() {
	c := NewCache[string, string](
		NewRestrictedMemorySizeBehavior(
			NewRestrictedTimeToLifeBehavior[string, string](nil, time.Hour, false),
			4,
			func(entry Entry[string, string]) int64 {
				val := entry.Value()
				return int64(len(val))
			},
		),
		NewSerialIndex[string, string](),
	)

	for i := 0; i < 100; i++ {
		key := fmt.Sprintf("%d", i)
		c.Set(key, key)
	}
	that.Equal(2, c.ItemCount())
}
