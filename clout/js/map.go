package js

import (
	"fmt"

	"github.com/tliron/kutil/ard"
	"github.com/tliron/kutil/util"
)

//
// Map
//

// TODO: support for fields

type Map []MapEntry

func (self *CloutContext) NewMap(list ard.List, keyMeta ard.StringMap, valueMeta ard.StringMap, functionCallContext FunctionCallContext) (Map, error) {
	map_ := make(Map, len(list))

	for index, data := range list {
		var err error
		if map_[index], err = self.NewMapEntry(data, keyMeta, valueMeta, functionCallContext); err != nil {
			return nil, err
		}
	}

	return map_, nil
}

func (self Map) Coerce() (ard.Value, error) {
	map_ := make(ard.StringMap)

	for _, entry := range self {
		if key, value, err := entry.Coerce(); err == nil {
			// Key should be a string value
			if _, ok := map_[key]; ok {
				if keyFunctionCall, ok := entry.Key.(*FunctionCall); ok {
					// ????
					if arguments, err := keyFunctionCall.CoerceArguments(); err == nil {
						return nil, keyFunctionCall.NewErrorf(arguments, "duplicate map key %q during coercion", key)
					} else {
						return nil, err
					}
				} else {
					return nil, fmt.Errorf("duplicate map key during coercion: %s", key)
				}
			}

			map_[key] = value
		} else {
			return nil, err
		}
	}

	return map_, nil
}

//
// MapEntry
//

type MapEntry struct {
	Key   Coercible `json:"key" yaml:"key"`
	Value Coercible `json:"value" yaml:"value"`
}

func (self *CloutContext) NewMapEntry(data any, keyMeta ard.StringMap, valueMeta ard.StringMap, functionCallContext FunctionCallContext) (MapEntry, error) {
	var entry MapEntry

	if map_, ok := data.(ard.StringMap); ok {
		if key, ok := map_["$key"]; ok {
			var err error
			if entry.Key, err = self.NewCoercible(key, keyMeta, functionCallContext); err == nil {
				if entry.Value, err = self.NewCoercible(map_, valueMeta, functionCallContext); err == nil {
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

func (self MapEntry) Coerce() (string, any, error) {
	if key, err := self.Key.Coerce(); err == nil {
		if value, err := self.Value.Coerce(); err == nil {
			return util.ToString(key), value, nil
		} else {
			return "", nil, err
		}
	} else {
		return "", nil, err
	}
}
