package js

import (
	"github.com/tliron/puccini/ard"
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
	c := make(CoercibleList, len(list))

	for index, data := range list {
		var err error
		if c[index], err = self.NewCoercible(data, site, source, target); err != nil {
			return nil, err
		}
	}

	return c, nil
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

type CoercibleMap map[string]Coercible

func (self *CloutContext) NewCoercibleMap(map_ ard.Map, site interface{}, source interface{}, target interface{}) (CoercibleMap, error) {
	c := make(CoercibleMap)

	for key, data := range map_ {
		var err error
		if c[key], err = self.NewCoercible(data, site, source, target); err != nil {
			return nil, err
		}
	}

	return c, nil
}

func (self CoercibleMap) Coerce() (interface{}, error) {
	value := make(ard.Map)

	for key, coercible := range self {
		var err error
		if value[key], err = coercible.Coerce(); err != nil {
			return nil, err
		}
	}

	return value, nil
}
