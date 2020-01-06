package hot

import (
	"github.com/tliron/puccini/ard"
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
	"github.com/tliron/yamlkeys"
)

//
// Value
//

type Value struct {
	*Entity `name:"value"`
	Name    string

	Constraints Constraints
	Description *string
}

func NewValue(context *tosca.Context) *Value {
	return &Value{
		Entity: NewEntity(context),
		Name:   context.Name,
	}
}

// tosca.Reader signature
func ReadValue(context *tosca.Context) interface{} {
	ToFunctionCalls(context)
	return NewValue(context)
}

// tosca.Mappable interface
func (self *Value) GetKey() string {
	return self.Name
}

func (self *Value) Normalize() normal.Constrainable {
	var constrainable normal.Constrainable

	switch data := self.Context.Data.(type) {
	case ard.List:
		list := normal.NewList(len(data))
		for index, value := range data {
			list.Set(index, NewValue(self.Context.ListChild(index, value)).Normalize())
		}
		constrainable = list

	case ard.Map:
		map_ := normal.NewMap()
		for key, value := range data {
			if _, ok := key.(string); !ok {
				// HOT does not support complex keys
				self.Context.MapChild(key, yamlkeys.KeyData(key)).ReportValueWrongType("string")
			}
			name := yamlkeys.KeyString(key)
			map_.Put(name, NewValue(self.Context.MapChild(name, value)).Normalize())
		}
		constrainable = map_

	case *tosca.FunctionCall:
		NormalizeFunctionCallArguments(data, self.Context)
		constrainable = normal.NewFunctionCall(data)

	default:
		constrainable = normal.NewValue(data)
	}

	self.Constraints.Normalize(self.Context, constrainable)

	if self.Description != nil {
		constrainable.SetDescription(*self.Description)
	}

	return constrainable
}

//
// Values
//

type Values map[string]*Value

func (self Values) Normalize(c normal.Constrainables) {
	for key, value := range self {
		c[key] = value.Normalize()
	}
}
