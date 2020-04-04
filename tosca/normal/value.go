package normal

import (
	"github.com/tliron/puccini/ard"
	"github.com/tliron/puccini/tosca"
)

//
// Value
//

type Value struct {
	Key         Constrainable `json:"$key,omitempty" yaml:"$key,omitempty"`
	Description string        `json:"$description,omitempty" yaml:"$description,omitempty"`
	Constraints FunctionCalls `json:"$constraints,omitempty" yaml:"$constraints,omitempty"`

	Value ard.Value `json:"$value" yaml:"$value"`
	Type  string    `json:"$type,omitempty" yaml:"$type,omitempty"`
}

func NewValue(value ard.Value) *Value {
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
