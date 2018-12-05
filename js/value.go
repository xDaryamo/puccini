package js

import (
	"github.com/tliron/puccini/ard"
)

//
// Value
//

type Value struct {
	Context     *CloutContext `json:"-" yaml:"-"`
	Value       interface{}   `json:"value" yaml:"value"`
	Constraints Constraints   `json:"constraints" yaml:"constraints"`
}

func (self *CloutContext) NewValue(data interface{}, site interface{}, source interface{}, target interface{}) (*Value, error) {
	c := Value{
		Context: self,
		Value:   data,
	}

	var err error
	if map_, ok := data.(ard.Map); ok {
		if v, ok := map_["value"]; ok {
			c.Value = v
		} else if v, ok := map_["list"]; ok {
			if l, ok := v.(ard.List); ok {
				c.Value, err = self.NewCoercibleList(l, site, source, target)
				if err != nil {
					return nil, err
				}
			} else {
				return &c, nil
			}
		} else if v, ok := map_["map"]; ok {
			if m, ok := v.(ard.Map); ok {
				c.Value, err = self.NewCoercibleMap(m, site, source, target)
				if err != nil {
					return nil, err
				}
			} else {
				return &c, nil
			}
		} else {
			return &c, nil
		}

		c.Constraints, err = self.NewConstraintsForValue(map_)
		if err != nil {
			return nil, err
		}
	}

	return &c, nil
}

// Coercible interface
func (self *Value) Coerce() (interface{}, error) {
	r := self.Value

	// Embedded Coercible (either a CoercibleList or a CoercibleMap)
	if c, ok := r.(Coercible); ok {
		var err error
		r, err = c.Coerce()
		if err != nil {
			return nil, err
		}
	}

	return self.Constraints.Apply(r)
}
