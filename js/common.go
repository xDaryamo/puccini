package js

import (
	"fmt"
	"strings"

	"github.com/op/go-logging"
	"github.com/tliron/puccini/ard"
)

var log = logging.MustGetLogger("js")

func SetMapNested(map_ ard.Map, key string, value string) error {
	path := strings.Split(key, ".")
	last := len(path) - 1
	if last == -1 {
		return fmt.Errorf("empty key")
	}
	if last > 0 {
		for _, p := range path[:last] {
			o, ok := map_[p]
			if ok {
				map_, ok = o.(ard.Map)
				if !ok {
					return fmt.Errorf("bad nested map structure")
				}
			} else {
				m := make(ard.Map)
				map_[p] = m
				map_ = m
			}
		}
	}
	map_[path[last]] = value
	return nil
}
