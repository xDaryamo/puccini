package util

import (
	"github.com/tliron/kutil/ard"
)

func Put(key string, value ard.Value, entity ard.Value, names ...string) bool {
	if map_, ok := ard.NewNode(entity).Get(names...).StringMap(); ok {
		map_[key] = value
		return true
	} else {
		return false
	}
}
