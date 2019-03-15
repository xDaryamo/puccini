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

	Data        interface{}
	Constraints Constraints
	Description *string
}

func NewValue(context *tosca.Context) *Value {
	return &Value{
		Entity: NewEntity(context),
		Name:   context.Name,
		Data:   context.Data,
	}
}

// tosca.Reader signature
func ReadValue(context *tosca.Context) interface{} {
	ToFunctions(context)
	return NewValue(context)
}

// tosca.Mappable interface
func (self *Value) GetKey() string {
	return self.Name
}

func (self *Value) Normalize() normal.Constrainable {
	var constrainable normal.Constrainable

	if list, ok := self.Data.(ard.List); ok {
		l := normal.NewConstrainableList(len(list))
		for index, value := range list {
			l.List[index] = NewValue(self.Context.ListChild(index, value)).Normalize()
		}
		constrainable = l
	} else if map_, ok := self.Data.(ard.Map); ok {
		m := normal.NewConstrainableMap()
		for key, value := range map_ {
			m.Map[key] = NewValue(self.Context.MapChild(key, value)).Normalize()
		}
		constrainable = m
	} else if function, ok := self.Data.(*tosca.Function); ok {
		NormalizeFunctionArguments(function, self.Context)
		constrainable = normal.NewFunction(function)
	} else {
		constrainable = normal.NewValue(self.Data)
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
