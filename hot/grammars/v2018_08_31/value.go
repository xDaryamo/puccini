package v2018_08_31

import (
	"github.com/tliron/puccini/tosca"
)

//
// Value
//

type Value struct {
	*Entity `name:"value"`
	Name    string

	Data interface{}
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
	/*if function, ok := GetFunction(context); ok {
		context.Data = function
	}*/
	return NewValue(context)
}

// tosca.Mappable interface
func (self *Value) GetKey() string {
	return self.Name
}

//
// Values
//

type Values map[string]*Value
