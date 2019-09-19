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

	switch self.Context.Data.(type) {
	case ard.List:
		list := self.Context.Data.(ard.List)
		l := normal.NewList(len(list))
		for index, value := range list {
			l.Set(index, NewValue(self.Context.ListChild(index, value)).Normalize())
		}
		constrainable = l

	case ard.Map:
		m := normal.NewMap()
		for key, value := range self.Context.Data.(ard.Map) {
			if _, ok := key.(string); !ok {
				// HOT does not support complex keys
				self.Context.MapChild(key, yamlkeys.KeyData(key)).ReportValueWrongType("string")
			}
			name := yamlkeys.KeyString(key)
			m.Put(name, NewValue(self.Context.MapChild(name, value)).Normalize())
		}
		constrainable = m

	case *tosca.FunctionCall:
		functionCall := self.Context.Data.(*tosca.FunctionCall)
		NormalizeFunctionCallArguments(functionCall, self.Context)
		constrainable = normal.NewFunctionCall(functionCall)

	default:
		constrainable = normal.NewValue(self.Context.Data)
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
