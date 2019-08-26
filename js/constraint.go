package js

import (
	"errors"

	"github.com/tliron/puccini/ard"
)

type Constraints []*FunctionCall

func (self *CloutContext) NewConstraints(list ard.List, site interface{}, source interface{}, target interface{}) (Constraints, error) {
	constraints := make(Constraints, len(list))

	for index, element := range list {
		var err error
		if constraints[index], err = self.NewFunctionCall(element, site, source, target); err != nil {
			return nil, err
		}
	}

	return constraints, nil
}

func (self *CloutContext) NewConstraintsForValue(map_ ard.Map, site interface{}, source interface{}, target interface{}) (Constraints, error) {
	v, ok := map_["constraints"]
	if !ok {
		return nil, nil
	}

	list, ok := v.(ard.List)
	if !ok {
		return nil, errors.New("malformed \"constraints\"")
	}

	return self.NewConstraints(list, site, source, target)
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
