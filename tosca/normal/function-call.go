package normal

import (
	"github.com/tliron/puccini/tosca"
)

//
// FunctionCall
//

type FunctionCall struct {
	Key         Constrainable `json:"$key,omitempty" yaml:"$key,omitempty"`
	Information *Information  `json:"$information,omitempty" yaml:"$information,omitempty"`
	Constraints FunctionCalls `json:"$constraints,omitempty" yaml:"$constraints,omitempty"`

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
func (self *FunctionCall) SetInformation(information *Information) {
	self.Information = CopyInformation(information)
}

// Constrainable interface
func (self *FunctionCall) AddConstraint(constraint *tosca.FunctionCall) {
	self.Constraints = append(self.Constraints, NewFunctionCall(constraint))
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
