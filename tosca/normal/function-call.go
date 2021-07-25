package normal

import (
	"github.com/tliron/puccini/tosca"
)

//
// FunctionCall
//

type FunctionCall struct {
	Key         Constrainable     `json:"$key,omitempty" yaml:"$key,omitempty"`
	Information *ValueInformation `json:"$information,omitempty" yaml:"$information,omitempty"`
	Constraints FunctionCalls     `json:"$constraints,omitempty" yaml:"$constraints,omitempty"`
	Converter   *FunctionCall     `json:"$converter,omitempty" yaml:"$converter,omitempty"`

	FunctionCall *tosca.FunctionCall `json:"$functionCall" yaml:"$functionCall"`
}

func NewFunctionCall(functionCall *tosca.FunctionCall) *FunctionCall {
	return &FunctionCall{FunctionCall: functionCall}
}

// Constrainable interface
func (self *FunctionCall) SetKey(key Constrainable) {
	self.Key = key
}

// Constrainable interface
func (self *FunctionCall) SetInformation(information *ValueInformation) {
	self.Information = CopyValueInformation(information)
}

// Constrainable interface
func (self *FunctionCall) AddConstraint(constraint *tosca.FunctionCall) {
	self.Constraints = append(self.Constraints, NewFunctionCall(constraint))
}

// Constrainable interface
func (self *FunctionCall) SetConverter(converter *tosca.FunctionCall) {
	self.Converter = NewFunctionCall(converter)
}

//
// FunctionCalls
//

type FunctionCalls []*FunctionCall

//
// FunctionCallMap
//

type FunctionCallMap map[string]FunctionCalls

//
// FunctionCallMapMap
//

type FunctionCallMapMap map[string]FunctionCallMap
