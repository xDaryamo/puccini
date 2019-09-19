package js

import (
	"fmt"

	"github.com/tliron/puccini/ard"
)

//
// Coercible
//

type Coercible interface {
	Coerce() (interface{}, error)
	SetConstraints(Constraints)
	Unwrap() interface{}
}

func (self *CloutContext) NewCoercible(data interface{}, site interface{}, source interface{}, target interface{}) (Coercible, error) {
	var coercible Coercible

	if notation, ok := data.(ard.StringMap); ok {
		var err error
		if data, ok := notation["functionCall"]; ok {
			if map_, ok := data.(ard.StringMap); ok {
				if functionCall, err := self.NewFunctionCall(map_, site, source, target); err == nil {
					functionCall.Notation = notation
					coercible = functionCall
				} else {
					return nil, err
				}
			} else {
				return nil, fmt.Errorf("malformed \"functionCall\", not a map: %T", data)
			}
		} else {
			value := Value{
				Context:  self,
				Notation: notation,
			}

			if data, ok := notation["value"]; ok {
				value.Value = data
			} else if data, ok := notation["list"]; ok {
				if list, ok := data.(ard.List); ok {
					var entryConstraints Constraints
					if entryConstraints, err = self.NewConstraintsForValue(notation, "entryConstraints", site, source, target); err != nil {
						return nil, err
					}

					if value.Value, err = self.NewList(list, entryConstraints, site, source, target); err != nil {
						return nil, err
					}
				} else {
					return nil, fmt.Errorf("malformed \"list\", not a list: %T", data)
				}
			} else if data, ok := notation["map"]; ok {
				if list, ok := data.(ard.List); ok {
					var keyConstraints Constraints
					if keyConstraints, err = self.NewConstraintsForValue(notation, "keyConstraints", site, source, target); err != nil {
						return nil, err
					}

					var valueConstraints Constraints
					if valueConstraints, err = self.NewConstraintsForValue(notation, "valueConstraints", site, source, target); err != nil {
						return nil, err
					}

					if value.Value, err = self.NewMap(list, keyConstraints, valueConstraints, site, source, target); err != nil {
						return nil, err
					}
				} else {
					return nil, fmt.Errorf("malformed \"map\", not a list: %T", data)
				}
			} else {
				return nil, fmt.Errorf("not a coercible, doesn't have \"value\", \"list\", \"map\", or \"functionCall\": %v", data)
			}

			coercible = &value
		}

		if constraints, err := self.NewConstraintsForValue(notation, "constraints", site, source, target); err == nil {
			coercible.SetConstraints(constraints)
			return coercible, nil
		} else {
			return nil, err
		}
	} else {
		return nil, fmt.Errorf("malformed coercible, not a map: %T", data)
	}
}
