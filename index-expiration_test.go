package cache

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"testing"
)

type IndexExpirationShould struct {
	suite.Suite
}

func TestIndexExpiration(t *testing.T) {
	suite.Run(t, new(IndexExpirationShould))
}

func (that *IndexExpirationShould) TestIndexExpiration() {
	index := NewExpirationIndex[string, string]()
	assert.NotNil(that.T(), index)

	item1 := &Item[string, string]{key: "key1", expiration: 1}
	item2 := &Item[string, string]{key: "key2", expiration: 2}
	item3 := &Item[string, string]{key: "key3", expiration: 3}

	index.Include(item1)
	index.Include(item2)
	index.Include(item3)

	assert.Equal(that.T(), 3, index.Len())
	assert.Equal(that.T(), item1, index.Items()[0])
	assert.Equal(that.T(), item2, index.Items()[1])
	assert.Equal(that.T(), item3, index.Items()[2])
}
