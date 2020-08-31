package js

import (
	"fmt"

	"github.com/tliron/kutil/ard"
)

//
// Constraints
//

type Constraints []*FunctionCall

func (self *CloutContext) NewConstraints(list ard.List, functionCallContext FunctionCallContext) (Constraints, error) {
	constraints := make(Constraints, len(list))

	for index, element := range list {
		if coercible, err := self.NewCoercible(element, functionCallContext); err == nil {
			var ok bool
			if constraints[index], ok = coercible.(*FunctionCall); !ok {
				return nil, fmt.Errorf("malformed constraint, not a function call: %+v", element)
			}
		} else {
			return nil, err
		}
	}

	return constraints, nil
}

func (self *CloutContext) NewConstraintsFromNotation(notation ard.StringMap, name string, functionCallContext FunctionCallContext) (Constraints, error) {
	if data, ok := notation[name]; ok {
		if list, ok := data.(ard.List); ok {
			return self.NewConstraints(list, functionCallContext)
		} else {
			return nil, fmt.Errorf("malformed %q, not a list: %T", name, data)
		}
	} else {
		return nil, nil
	}
}

func (self Constraints) Validate(value interface{}) (bool, error) {
	if coercible, ok := value.(Coercible); ok {
		var err error
		if value, err = coercible.Coerce(); err != nil {
			return false, err
		}
	}

	for _, constraint := range self {
		if valid, err := constraint.Validate(value, false); err == nil {
			if !valid {
				return false, nil
			}
		} else {
			return false, err
		}
	}

	return true, nil
}

func (self Constraints) Apply(value interface{}) (interface{}, error) {
	if coercible, ok := value.(Coercible); ok {
		var err error
		if value, err = coercible.Coerce(); err != nil {
			return nil, err
		}
	}

	for _, constraint := range self {
		if _, err := constraint.Validate(value, true); err != nil {
			return nil, err
		}
	}

	return value, nil
}
