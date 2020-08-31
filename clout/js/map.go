package js

import (
	"fmt"

	"github.com/tliron/kutil/ard"
)

//
// Map
//

type Map []MapEntry

func (self *CloutContext) NewMap(list ard.List, keyConstraints Constraints, valueConstraints Constraints, functionCallContext FunctionCallContext) (Map, error) {
	var map_ Map

	for _, data := range list {
		if entry, err := self.NewMapEntry(data, keyConstraints, valueConstraints, functionCallContext); err == nil {
			map_ = append(map_, entry)
		} else {
			return nil, err
		}
	}

	return map_, nil
}

func (self Map) Coerce() (ard.Value, error) {
	value := make(ard.StringMap)

	for _, entry := range self {
		if k, v, err := entry.Coerce(); err == nil {
			if _, ok := value[k]; ok {
				// TODO: include location information, as with function calls?
				return nil, fmt.Errorf("duplicate map key during coercion: %s", k)
			}
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
	Key   Coercible `json:"$key" yaml:"$key"`
	Value Coercible `json:"$value" yaml:"$value"`
}

func (self *CloutContext) NewMapEntry(data interface{}, keyConstraints Constraints, valueConstraints Constraints, functionCallContext FunctionCallContext) (MapEntry, error) {
	var entry MapEntry

	if map_, ok := data.(ard.StringMap); ok {
		if key, ok := map_["$key"]; ok {
			var err error
			if entry.Key, err = self.NewCoercible(key, functionCallContext); err == nil {
				if entry.Value, err = self.NewCoercible(map_, functionCallContext); err == nil {
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
			return entry, fmt.Errorf("malformed map entry, no \"$key\": %+v", map_)
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
