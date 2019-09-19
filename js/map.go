package js

import (
	"fmt"

	"github.com/tliron/puccini/ard"
)

//
// Map
//

type Map []MapEntry

func (self *CloutContext) NewMap(list ard.List, keyConstraints Constraints, valueConstraints Constraints, site interface{}, source interface{}, target interface{}) (Map, error) {
	var map_ Map

	for _, data := range list {
		if entry, err := self.NewMapEntry(data, keyConstraints, valueConstraints, site, source, target); err == nil {
			map_ = append(map_, entry)
		} else {
			return nil, err
		}
	}

	return map_, nil
}

func (self Map) Coerce() (interface{}, error) {
	value := make(ard.StringMap)

	for _, entry := range self {
		if k, v, err := entry.Coerce(); err == nil {
			value[k] = v
		} else {
			return nil, err
		}
	}

	return value, nil
}

//
// MapEntry
//

type MapEntry struct {
	Key   Coercible `json:"key" yaml:"key"`
	Value Coercible `json:"value" yaml:"value"`
}

func (self *CloutContext) NewMapEntry(data interface{}, keyConstraints Constraints, valueConstraints Constraints, site interface{}, source interface{}, target interface{}) (MapEntry, error) {
	var entry MapEntry

	if map_, ok := data.(ard.StringMap); ok {
		if key, ok := map_["key"]; ok {
			var err error
			if entry.Key, err = self.NewCoercible(key, site, source, target); err == nil {
				if entry.Value, err = self.NewCoercible(map_, site, source, target); err == nil {
					entry.Key.SetConstraints(keyConstraints)
					entry.Value.SetConstraints(valueConstraints)
					return entry, nil
				} else {
					return entry, err
				}
			} else {
				return entry, err
			}
		} else {
			return entry, fmt.Errorf("malformed map entry, no \"key\": %v", map_)
		}
	} else {
		return entry, fmt.Errorf("malformed map entry, not a map: %T", data)
	}
}

func (self MapEntry) Coerce() (string, interface{}, error) {
	if key, err := self.Key.Coerce(); err == nil {
		if value, err := self.Value.Coerce(); err == nil {
			return fmt.Sprintf("%v", key), value, nil
		} else {
			return "", nil, err
		}
	} else {
		return "", nil, err
	}
}
