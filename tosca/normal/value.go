package normal

import (
	"github.com/tliron/puccini/tosca"
)

//
// Value
//

type Value struct {
	Key         Constrainable `json:"key,omitempty" yaml:"key,omitempty"`
	Description string        `json:"description,omitempty" yaml:"description,omitempty"`
	Constraints FunctionCalls `json:"constraints,omitempty" yaml:"constraints,omitempty"`

	Value interface{} `json:"value" yaml:"value"`
}

func NewValue(value interface{}) *Value {
	return &Value{Value: value}
}

// Constrainable interface
func (self *Value) SetKey(key Constrainable) {
	self.Key = key
}

// Constrainable interface
func (self *Value) SetDescription(description string) {
	self.Description = description
}

// Constrainable interface
func (self *Value) AddConstraint(constraint *tosca.FunctionCall) {
	self.Constraints = append(self.Constraints, NewFunctionCall(constraint))
}
