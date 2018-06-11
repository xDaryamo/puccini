package js

import (
	"github.com/tliron/puccini/ard"
	"github.com/tliron/puccini/clout"
)

//
// Coercible
//

type Coercible interface {
	Coerce() (interface{}, error)
}

func NewCoercible(data interface{}, site interface{}, source interface{}, target interface{}, c *clout.Clout) (Coercible, error) {
	function, err := NewFunction(data, site, source, target, c)
	if err == nil {
		return function, nil
	}
	return NewValue(data, site, source, target, c)
}

//
// CoercibleList
//

type CoercibleList []Coercible

func NewCoercibleList(list ard.List, site interface{}, source interface{}, target interface{}, c *clout.Clout) (CoercibleList, error) {
	self := make(CoercibleList, len(list))
	for index, data := range list {
		var err error
		self[index], err = NewCoercible(data, site, source, target, c)
		if err != nil {
			return nil, err
		}
	}
	return self, nil
}

// Coercible interface
func (self CoercibleList) Coerce() (interface{}, error) {
	value := make(ard.List, len(self))
	for index, coercible := range self {
		var err error
		value[index], err = coercible.Coerce()
		if err != nil {
			return nil, err
		}
	}
	return value, nil
}

//
// CoercibleMap
//

type CoercibleMap map[string]Coercible

func NewCoercibleMap(map_ ard.Map, site interface{}, source interface{}, target interface{}, c *clout.Clout) (CoercibleMap, error) {
	self := make(CoercibleMap)
	for key, data := range map_ {
		var err error
		self[key], err = NewCoercible(data, site, source, target, c)
		if err != nil {
			return nil, err
		}
	}
	return self, nil
}

// Coercible interface
func (self CoercibleMap) Coerce() (interface{}, error) {
	value := make(ard.Map)
	for key, coercible := range self {
		var err error
		value[key], err = coercible.Coerce()
		if err != nil {
			return nil, err
		}
	}
	return value, nil
}
