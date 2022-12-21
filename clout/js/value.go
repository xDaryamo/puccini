package js

import (
	"fmt"

	"github.com/tliron/kutil/ard"
)

//
// Value
//

type Value struct {
	Data       any           `json:"data" yaml:"data"` // List, Map, or ard.Value
	Validators Validators    `json:"validators,omitempty" yaml:"validators,omitempty"`
	Converter  *FunctionCall `json:"converter,omitempty" yaml:"converter,omitempty"`

	Notation ard.StringMap `json:"-" yaml:"-"`
}

func (self *ExecutionContext) NewValue(data ard.Value, notation ard.StringMap, meta ard.StringMap) (*Value, error) {
	value := Value{
		Data:     data,
		Notation: notation,
	}

	var err error
	if value.Validators, err = self.NewValidatorsFromMeta(meta); err != nil {
		return nil, err
	}
	if value.Converter, err = self.NewConverter(meta); err != nil {
		return nil, err
	}

	return &value, nil
}

func (self *ExecutionContext) NewListValue(list ard.List, notation ard.StringMap, meta ard.StringMap) (*Value, error) {
	var elementMeta ard.StringMap
	if meta != nil {
		if data, ok := meta["element"]; ok {
			if map_, ok := data.(ard.StringMap); ok {
				elementMeta = map_
			} else {
				return nil, fmt.Errorf("malformed meta \"element\", not a map: %T", data)
			}
		}
	}

	if list_, err := self.NewList(list, elementMeta); err == nil {
		return self.NewValue(list_, notation, meta)
	} else {
		return nil, err
	}
}

func (self *ExecutionContext) NewMapValue(list ard.List, notation ard.StringMap, meta ard.StringMap) (*Value, error) {
	var keyMeta ard.StringMap
	var valueMeta ard.StringMap
	var fieldsMeta ard.StringMap
	if meta != nil {
		if data, ok := meta["key"]; ok {
			if map_, ok := data.(ard.StringMap); ok {
				keyMeta = map_
			} else {
				return nil, fmt.Errorf("malformed meta \"key\", not a map: %T", data)
			}
		}

		if data, ok := meta["value"]; ok {
			if map_, ok := data.(ard.StringMap); ok {
				valueMeta = map_
			} else {
				return nil, fmt.Errorf("malformed meta \"value\", not a map: %T", data)
			}
		}

		if data, ok := meta["fields"]; ok {
			if map_, ok := data.(ard.StringMap); ok {
				fieldsMeta = map_
			} else {
				return nil, fmt.Errorf("malformed meta \"fields\", not a map: %T", data)
			}
		}
	}

	if map_, err := self.NewMap(list, keyMeta, valueMeta, fieldsMeta); err == nil {
		return self.NewValue(map_, notation, meta)
	} else {
		return nil, err
	}
}

// Coercible interface
func (self *Value) Coerce() (ard.Value, error) {
	data := self.Data

	var err error
	switch data_ := data.(type) {
	case List:
		if data, err = data_.Coerce(); err != nil {
			return nil, err
		}

	case Map:
		if data, err = data_.Coerce(); err != nil {
			return nil, err
		}
	}

	if err := self.Validators.Apply(data); err == nil {
		if self.Converter != nil {
			return self.Converter.Convert(data)
		} else {
			return data, nil
		}
	} else {
		return nil, err
	}
}

// Coercible interface
func (self *Value) AddValidators(validators Validators) {
	self.Validators = append(self.Validators, validators...)
}

// Coercible interface
func (self *Value) Unwrap() ard.Value {
	return self.Notation
}
