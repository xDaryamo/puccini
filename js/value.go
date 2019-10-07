package js

import (
	"github.com/tliron/puccini/ard"
)

//
// Value
//

type Value struct {
	Value       interface{} `json:"value" yaml:"value"`
	Constraints Constraints `json:"constraints" yaml:"constraints"`

	Notation ard.StringMap `json:"-" yaml:"-"`
}

// Coercible interface
func (self *Value) Coerce() (interface{}, error) {
	value := self.Value

	var err error
	switch value.(type) {
	case List:
		if value, err = value.(List).Coerce(); err != nil {
			return nil, err
		}

	case Map:
		if value, err = value.(Map).Coerce(); err != nil {
			return nil, err
		}
	}

	return self.Constraints.Apply(value)
}

// Coercible interface
func (self *Value) SetConstraints(constraints Constraints) {
	self.Constraints = constraints
}

// Coercible interface
func (self *Value) Unwrap() interface{} {
	return self.Notation
}
