package js

import (
	"github.com/tliron/puccini/ard"
)

//
// Value
//

type Value struct {
	Context     *CloutContext `json:"-" yaml:"-"`
	Data        interface{}   `json:"-" yaml:"-"`
	Value       interface{}   `json:"value" yaml:"value"`
	Constraints Constraints   `json:"constraints" yaml:"constraints"`
}

func (self *CloutContext) NewValue(data interface{}, site interface{}, source interface{}, target interface{}) (*Value, error) {
	c := Value{
		Context: self,
		Data:    data,
		Value:   data,
	}

	if map_, ok := data.(ard.StringMap); ok {
		var err error
		if v, ok := map_["value"]; ok {
			c.Value = v
		} else if v, ok := map_["list"]; ok {
			if l, ok := v.(ard.List); ok {
				// Embedded CoercibleList
				if c.Value, err = self.NewCoercibleList(l, site, source, target); err != nil {
					return nil, err
				}
			}
		} else if v, ok := map_["map"]; ok {
			if l, ok := v.(ard.List); ok {
				// Embedded CoercibleMap
				if c.Value, err = self.NewCoercibleMap(l, site, source, target); err != nil {
					return nil, err
				}
			}
		}

		if c.Constraints, err = self.NewConstraintsForValue(map_, site, source, target); err != nil {
			return nil, err
		}
	}

	return &c, nil
}

// Coercible interface
func (self *Value) Coerce() (interface{}, error) {
	r := self.Value

	var err error
	if coercibleList, ok := r.(CoercibleList); ok {
		// Embedded CoercibleList
		if r, err = coercibleList.Coerce(); err != nil {
			return nil, err
		}
	} else if coercibleMap, ok := r.(CoercibleMap); ok {
		// Embedded CoercibleMap
		if r, err = coercibleMap.Coerce(); err != nil {
			return nil, err
		}
	}

	return self.Constraints.Apply(r)
}

// Coercible interface
func (self *Value) Unwrap() interface{} {
	return self.Data
}
