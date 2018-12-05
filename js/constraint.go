package js

import (
	"fmt"

	"github.com/tliron/puccini/ard"
)

type Constraints []*Function

func (self *CloutContext) NewConstraints(list ard.List) (Constraints, error) {
	constraints := make(Constraints, len(list))
	for index, element := range list {
		var err error
		constraints[index], err = self.NewFunction(element, nil, nil, nil)
		if err != nil {
			return nil, err
		}
	}

	return constraints, nil
}

func (self *CloutContext) NewConstraintsForValue(map_ ard.Map) (Constraints, error) {
	v, ok := map_["constraints"]
	if !ok {
		return nil, nil
	}

	list, ok := v.(ard.List)
	if !ok {
		return nil, fmt.Errorf("malformed \"constraints\"")
	}

	return self.NewConstraints(list)
}

func (self Constraints) Validate(value interface{}) (bool, error) {
	// Coerce value
	if coercible, ok := value.(Coercible); ok {
		var err error
		value, err = coercible.Coerce()
		if err != nil {
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
	// Coerce value
	if coercible, ok := value.(Coercible); ok {
		var err error
		value, err = coercible.Coerce()
		if err != nil {
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
