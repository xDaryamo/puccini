package js

import (
	"github.com/tliron/kutil/ard"
)

//
// Value
//

type Value struct {
	Notation ard.StringMap `json:"-" yaml:"-"`

	Data        interface{} `json:"data" yaml:"data"` // List, Map, or ard.Value
	Constraints Constraints `json:"constraints,omitempty" yaml:"constraints,omitempty"`
}

func (self *CloutContext) NewValue(data ard.Value, notation ard.StringMap, functionCallContext FunctionCallContext) (*Value, error) {
	value := Value{
		Data:     data,
		Notation: notation,
	}

	var err error
	if value.Constraints, err = self.NewConstraintsFromNotation(notation, "$constraints", functionCallContext); err != nil {
		return nil, err
	}

	return &value, nil
}

func (self *CloutContext) NewValueForList(list ard.List, notation ard.StringMap, functionCallContext FunctionCallContext) (*Value, error) {
	if entryConstraints, err := self.NewConstraintsFromNotation(notation, "$entryConstraints", functionCallContext); err == nil {
		if list_, err := self.NewList(list, entryConstraints, functionCallContext); err == nil {
			return self.NewValue(list_, notation, functionCallContext)
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

func (self *CloutContext) NewValueForMap(list ard.List, notation ard.StringMap, functionCallContext FunctionCallContext) (*Value, error) {
	if keyConstraints, err := self.NewConstraintsFromNotation(notation, "$keyConstraints", functionCallContext); err == nil {
		if valueConstraints, err := self.NewConstraintsFromNotation(notation, "$valueConstraints", functionCallContext); err == nil {
			if map_, err := self.NewMap(list, keyConstraints, valueConstraints, functionCallContext); err == nil {
				return self.NewValue(map_, notation, functionCallContext)
			} else {
				return nil, err
			}
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

// Coercible interface
func (self *Value) Coerce() (ard.Value, error) {
	data := self.Data

	var err error
	switch data_ := data.(type) {
	case List:
		if data, err = data_.Coerce(); err != nil {
			return nil, err
		}

	case Map:
		if data, err = data_.Coerce(); err != nil {
			return nil, err
		}
	}

	return self.Constraints.Apply(data)
}

// Coercible interface
func (self *Value) SetConstraints(constraints Constraints) {
	self.Constraints = constraints
}

// Coercible interface
func (self *Value) Unwrap() ard.Value {
	return self.Notation
}
