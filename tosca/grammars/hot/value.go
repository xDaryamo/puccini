package hot

import (
	"github.com/tliron/puccini/ard"
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
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

	if list, ok := self.Context.Data.(ard.List); ok {
		l := normal.NewConstrainableList(len(list))
		for index, value := range list {
			l.List[index] = NewValue(self.Context.ListChild(index, value)).Normalize()
		}
		constrainable = l
	} else if map_, ok := self.Context.Data.(ard.Map); ok {
		m := normal.NewConstrainableMap()
		for key, value := range map_ {
			m.Map[key] = NewValue(self.Context.MapChild(key, value)).Normalize()
		}
		constrainable = m
	} else if functionCall, ok := self.Context.Data.(*tosca.FunctionCall); ok {
		NormalizeFunctionCallArguments(functionCall, self.Context)
		constrainable = normal.NewFunctionCall(functionCall)
	} else {
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
