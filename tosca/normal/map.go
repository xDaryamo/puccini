package normal

import (
	"github.com/tliron/puccini/tosca"
)

//
// Map
//

type Map struct {
	Key         Constrainable     `json:"$key,omitempty" yaml:"$key,omitempty"`
	Information *ValueInformation `json:"$information,omitempty" yaml:"$information,omitempty"`
	Constraints FunctionCalls     `json:"$constraints,omitempty" yaml:"$constraints,omitempty"`
	Converter   *FunctionCall     `json:"$converter,omitempty" yaml:"$converter,omitempty"`

	KeyConstraints   FunctionCalls `json:"$keyConstraints,omitempty" yaml:"$keyConstraints,omitempty"`
	ValueConstraints FunctionCalls `json:"$valueConstraints,omitempty" yaml:"$valueConstraints,omitempty"`

	Entries ConstrainableList `json:"$map" yaml:"$map"`
}

func NewMap() *Map {
	return new(Map)
}

// Constrainable interface
func (self *Map) SetKey(key Constrainable) {
	self.Key = key
}

// Constrainable interface
func (self *Map) SetInformation(information *ValueInformation) {
	self.Information = CopyValueInformation(information)
}

// Constrainable interface
func (self *Map) AddConstraint(constraint *tosca.FunctionCall) {
	self.Constraints = append(self.Constraints, NewFunctionCall(constraint))
}

// Constrainable interface
func (self *Map) SetConverter(converter *tosca.FunctionCall) {
	self.Converter = NewFunctionCall(converter)
}

func (self *Map) AddKeyConstraint(constraint *tosca.FunctionCall) {
	self.KeyConstraints = append(self.KeyConstraints, NewFunctionCall(constraint))
}

func (self *Map) AddValueConstraint(constraint *tosca.FunctionCall) {
	self.ValueConstraints = append(self.ValueConstraints, NewFunctionCall(constraint))
}

func (self *Map) Put(key interface{}, value Constrainable) {
	self.Entries = self.Entries.AppendWithKey(key, value)
}
