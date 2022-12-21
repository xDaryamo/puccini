package js

import (
	"fmt"

	"github.com/tliron/kutil/ard"
)

//
// Coercible
//

type Coercible interface {
	Coerce() (ard.Value, error)
	AddValidators(Validators)
	Unwrap() ard.Value
}

func (self *ExecutionContext) NewCoercible(data ard.Value, meta ard.StringMap) (Coercible, error) {
	if notation, ok := data.(ard.StringMap); ok {
		if data, ok := notation["$primitive"]; ok {
			return self.NewValue(data, notation, meta)
		} else if data, ok := notation["$list"]; ok {
			if list, ok := asList(data); ok {
				return self.NewListValue(list, notation, meta)
			} else {
				return nil, fmt.Errorf("malformed \"$list\", not a list: %T", data)
			}
		} else if data, ok := notation["$map"]; ok {
			if list, ok := asList(data); ok {
				return self.NewMapValue(list, notation, meta)
			} else {
				return nil, fmt.Errorf("malformed \"$map\", not a list of entries: %T", data)
			}
		} else if data, ok := notation["$functionCall"]; ok {
			if map_, ok := asStringMap(data); ok {
				return self.NewFunctionCall(map_, notation, meta)
			} else {
				return nil, fmt.Errorf("malformed \"$functionCall\", not a map: %T", data)
			}
		} else {
			return nil, fmt.Errorf("not a coercible, doesn't have \"$primitive\", \"$list\", \"$map\", or \"$functionCall\": %+v", data)
		}
	} else {
		return nil, fmt.Errorf("malformed coercible, not a map: %T", data)
	}
}
