package cache

import (
	"reflect"
)

// isZeroVal check if any type is its zero value
func isZeroVal(x interface{}) bool {
	return x == nil || reflect.DeepEqual(x, reflect.Zero(reflect.TypeOf(x)).Interface())
}
