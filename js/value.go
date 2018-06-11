package js

import (
	"github.com/tliron/puccini/ard"
	"github.com/tliron/puccini/clout"
)

//
// Value
//

type Value struct {
	Clout       *clout.Clout `json:"-" yaml:"-"`
	Value       interface{}  `json:"value" yaml:"value"`
	Constraints Constraints  `json:"constraints" yaml:"constraints"`
}

func NewValue(data interface{}, site interface{}, source interface{}, target interface{}, c *clout.Clout) (*Value, error) {
	self := Value{Clout: c, Value: data}

	var err error
	if map_, ok := data.(ard.Map); ok {
		if v, ok := map_["value"]; ok {
			self.Value = v
		} else if v, ok := map_["list"]; ok {
			if l, ok := v.(ard.List); ok {
				self.Value, err = NewCoercibleList(l, site, source, target, c)
				if err != nil {
					return nil, err
				}
			} else {
				return &self, nil
			}
		} else if v, ok := map_["map"]; ok {
			if m, ok := v.(ard.Map); ok {
				self.Value, err = NewCoercibleMap(m, site, source, target, c)
				if err != nil {
					return nil, err
				}
			} else {
				return &self, nil
			}
		} else {
			return &self, nil
		}

		self.Constraints, err = NewConstraints(map_, c)
		if err != nil {
			return nil, err
		}
	}

	return &self, nil
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
