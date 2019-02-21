package hot

import (
	"strings"

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
	if function, ok := GetFunction(context); ok {
		context.Data = function
	}
	return NewValue(context)
}

// tosca.Mappable interface
func (self *Value) GetKey() string {
	return self.Name
}

func (self *Value) RenderParameter(parameter *Parameter) {
	if (self.Data == nil) && (parameter.Default != nil) {
		self.Data = parameter.Default
	}

	self.Constraints = parameter.Constraints
	self.Description = parameter.Description
}

func (self *Value) Normalize() normal.Constrainable {
	var constrainable normal.Constrainable

	if function, ok := self.Data.(*tosca.Function); ok {
		NormalizeFunctionArguments(function, self.Context)
		constrainable = normal.NewFunction(function)
	} else {
		constrainable = normal.NewValue(self.Data)
	}

	if self.Description != nil {
		constrainable.SetDescription(*self.Description)
	}

	return constrainable
}

func (self *Value) Fix(type_ string) {
	switch type_ {
	case "boolean":
		switch v := self.Data.(type) {
		case string:
			switch v {
			case "t", "true", "on", "y", "yes", "1":
				self.Data = true
			case "f", "false", "off", "n", "no", "0":
				self.Data = false
			}
		case int:
			switch v {
			case 1:
				self.Data = true
			case 0:
				self.Data = false
			}
		}
	case "comma_delimited_list":
		switch v := self.Data.(type) {
		case string:
			var list ard.List
			for _, s := range strings.Split(v, ",") {
				list = append(list, s)
			}
			self.Data = list
		}
	}
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
