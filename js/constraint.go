package js

import (
	"fmt"

	"github.com/tliron/puccini/ard"
)

type Constraints []*Function

func (self *CloutContext) NewConstraints(map_ ard.Map) (Constraints, error) {
	v, ok := map_["constraints"]
	if !ok {
		return nil, nil
	}

	list, ok := v.(ard.List)
	if !ok {
		return nil, fmt.Errorf("malformed \"constraints\"")
	}

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

func (self Constraints) Apply(value interface{}) (interface{}, error) {
	// Coerce value
	if coercible, ok := value.(Coercible); ok {
		var err error
		value, err = coercible.Coerce()
		if err != nil {
			return false, err
		}
	}

	for _, constraint := range self {
		_, err := constraint.Validate(value)
		if err != nil {
			return nil, err
		}
	}

	return value, nil
}
