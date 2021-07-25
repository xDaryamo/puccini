package normal

import (
	"github.com/tliron/kutil/ard"
	"github.com/tliron/puccini/tosca"
)

//
// Value
//

type Value struct {
	Key         Constrainable     `json:"$key,omitempty" yaml:"$key,omitempty"`
	Information *ValueInformation `json:"$information,omitempty" yaml:"$information,omitempty"`
	Constraints FunctionCalls     `json:"$constraints,omitempty" yaml:"$constraints,omitempty"`
	Converter   *FunctionCall     `json:"$converter,omitempty" yaml:"$converter,omitempty"`

	Value ard.Value `json:"$value" yaml:"$value"`
}

func NewValue(value ard.Value) *Value {
	return &Value{Value: value}
}

// Constrainable interface
func (self *Value) SetKey(key Constrainable) {
	self.Key = key
}

// Constrainable interface
func (self *Value) SetInformation(information *ValueInformation) {
	self.Information = CopyValueInformation(information)
}

// Constrainable interface
func (self *Value) AddConstraint(constraint *tosca.FunctionCall) {
	self.Constraints = append(self.Constraints, NewFunctionCall(constraint))
}

// Constrainable interface
func (self *Value) SetConverter(converter *tosca.FunctionCall) {
	self.Converter = NewFunctionCall(converter)
}
