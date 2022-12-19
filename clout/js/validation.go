package js

import (
	"fmt"

	"github.com/tliron/kutil/ard"
)

//
// Validators
//

type Validators []*FunctionCall

func (self *CloutContext) NewValidators(list ard.List, meta ard.StringMap, functionCallContext FunctionCallContext) (Validators, error) {
	validators := make(Validators, len(list))

	for index, element := range list {
		if value, err := self.NewCoercible(element, nil, functionCallContext); err == nil {
			var ok bool
			if validators[index], ok = value.(*FunctionCall); !ok {
				return nil, fmt.Errorf("malformed validator, not a function call: %+v", element)
			}
		} else {
			return nil, err
		}
	}

	return validators, nil
}

func (self *CloutContext) NewValidatorsFromMeta(meta ard.StringMap, functionCallContext FunctionCallContext) (Validators, error) {
	if meta != nil {
		if data, ok := meta["validators"]; ok {
			if list, ok := data.(ard.List); ok {
				return self.NewValidators(list, nil, functionCallContext)
			} else {
				return nil, fmt.Errorf("malformed \"validators\", not a list: %T", data)
			}
		} else {
			return nil, nil
		}
	} else {
		return nil, nil
	}
}

// Called from JavaScript
func (self Validators) IsValid(value any) (bool, error) {
	if value_, ok := value.(Coercible); ok {
		var err error
		if value, err = value_.Coerce(); err != nil {
			return false, err
		}
	}

	for _, validator := range self {
		if valid, err := validator.Validate(value, false); err == nil {
			if !valid {
				return false, nil
			}
		} else {
			return false, err
		}
	}

	return true, nil
}

func (self Validators) Apply(value any) error {
	if value_, ok := value.(Coercible); ok {
		var err error
		if value, err = value_.Coerce(); err != nil {
			return err
		}
	}

	for _, validator := range self {
		if _, err := validator.Validate(value, true); err != nil {
			return err
		}
	}

	return nil
}
