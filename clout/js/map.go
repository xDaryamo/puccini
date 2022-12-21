package js

import (
	"fmt"

	"github.com/tliron/kutil/ard"
	"github.com/tliron/kutil/util"
)

//
// Map
//

type Map []MapEntry

func (self *ExecutionContext) NewMap(list ard.List, keyMeta ard.StringMap, valueMeta ard.StringMap, fieldsMeta ard.StringMap) (Map, error) {
	map_ := make(Map, len(list))

	if fieldsMeta != nil {
		if (keyMeta != nil) || (valueMeta != nil) {
			return nil, fmt.Errorf("has \"fields\" meta in addition to \"key\" or \"value\" meta")
		}
	}

	for index, data := range list {
		var err error
		if map_[index], err = self.NewMapEntry(data, keyMeta, valueMeta, fieldsMeta); err != nil {
			return nil, err
		}
	}

	return map_, nil
}

func (self Map) Coerce() (ard.Value, error) {
	map_ := make(ard.StringMap)

	for _, entry := range self {
		if key, value, err := entry.Coerce(); err == nil {
			if _, ok := map_[key]; ok {
				return nil, fmt.Errorf("duplicate map key during coercion: %s", key)
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

func (self *ExecutionContext) NewMapEntry(data any, keyMeta ard.StringMap, valueMeta ard.StringMap, fieldsMeta ard.StringMap) (MapEntry, error) {
	var entry MapEntry

	if map_, ok := data.(ard.StringMap); ok {
		if key, ok := map_["$key"]; ok {
			var err error
			if entry.Key, err = self.NewCoercible(key, keyMeta); err == nil {
				if fieldsMeta != nil {
					// Find field meta
					if key, err := entry.Key.Coerce(); err == nil {
						if key_, ok := key.(string); ok {
							if valueMeta_, ok := fieldsMeta[key_]; ok {
								if valueMeta, ok = valueMeta_.(ard.StringMap); !ok {
									return entry, fmt.Errorf("malformed meta \"fields\", not a map: %T", valueMeta_)
								}
							} else {
								return entry, fmt.Errorf("malformed meta \"fields\", field not found: %q", key)
							}
						} else {
							return entry, fmt.Errorf("malformed field name, not a string: %T", key)
						}
					} else {
						return entry, err
					}
				}

				if entry.Value, err = self.NewCoercible(map_, valueMeta); err == nil {
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
