package js

import (
	"fmt"

	"github.com/tliron/puccini/ard"
	"github.com/tliron/yamlkeys"
)

//
// Coercible
//

type Coercible interface {
	Coerce() (interface{}, error)
	Unwrap() interface{}
}

func (self *CloutContext) NewCoercible(data interface{}, site interface{}, source interface{}, target interface{}) (Coercible, error) {
	if functionCall, err := self.NewFunctionCall(data, site, source, target); err == nil {
		return functionCall, nil
	}

	return self.NewValue(data, site, source, target)
}

//
// CoercibleList
//

type CoercibleList []Coercible

func (self *CloutContext) NewCoercibleList(list ard.List, site interface{}, source interface{}, target interface{}) (CoercibleList, error) {
	coercible := make(CoercibleList, len(list))

	for index, data := range list {
		var err error
		if coercible[index], err = self.NewCoercible(data, site, source, target); err != nil {
			return nil, err
		}
	}

	return coercible, nil
}

func (self CoercibleList) Coerce() (interface{}, error) {
	value := make(ard.List, len(self))

	for index, coercible := range self {
		var err error
		if value[index], err = coercible.Coerce(); err != nil {
			return nil, err
		}
	}

	return value, nil
}

//
// CoercibleMap
//

type CoercibleMap []CoercibleMapEntry

type CoercibleMapEntry struct {
	Key       interface{}
	Coercible Coercible
}

func (self *CloutContext) NewCoercibleMap(list ard.List, site interface{}, source interface{}, target interface{}) (CoercibleMap, error) {
	var coercible CoercibleMap

	for _, data := range list {
		if map_, ok := data.(ard.StringMap); ok {
			if key, ok := map_["key"]; ok {
				if c, err := self.NewCoercible(data, site, source, target); err == nil {
					coercible = append(coercible, CoercibleMapEntry{key, c})
				} else {
					return nil, err
				}
			} else {
				return nil, fmt.Errorf("map entry does not have \"key\": %v", data)
			}
		} else {
			return nil, fmt.Errorf("map entry is not a map: %v", data)
		}
	}

	return coercible, nil
}

func (self CoercibleMap) Coerce() (interface{}, error) {
	value := make(ard.StringMap)

	for _, entry := range self {
		var err error
		// TODO: support complex keys?
		if value[yamlkeys.KeyString(entry.Key)], err = entry.Coercible.Coerce(); err != nil {
			return nil, err
		}
	}

	return value, nil
}
