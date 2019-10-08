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

func (self *CloutContext) NewCoercible(data interface{}, functionCallContext FunctionCallContext) (Coercible, error) {
	if notation, ok := data.(ard.StringMap); ok {
		if data, ok := notation["value"]; ok {
			return self.NewValue(data, notation, functionCallContext)
		} else if data, ok := notation["list"]; ok {
			if list, ok := data.(ard.List); ok {
				return self.NewValueForList(list, notation, functionCallContext)
			} else {
				return nil, fmt.Errorf("malformed \"list\", not a list: %T", data)
			}
		} else if data, ok := notation["map"]; ok {
			if list, ok := data.(ard.List); ok {
				return self.NewValueForMap(list, notation, functionCallContext)
			} else {
				return nil, fmt.Errorf("malformed \"map\", not a list: %T", data)
			}
		} else if data, ok := notation["functionCall"]; ok {
			if map_, ok := data.(ard.StringMap); ok {
				return self.NewFunctionCall(map_, notation, functionCallContext)
			} else {
				return nil, fmt.Errorf("malformed \"functionCall\", not a map: %T", data)
			}
		} else {
			return nil, fmt.Errorf("not a coercible, doesn't have \"value\", \"list\", \"map\", or \"functionCall\": %v", data)
		}
	} else {
		return nil, fmt.Errorf("malformed coercible, not a map: %T", data)
	}
}
