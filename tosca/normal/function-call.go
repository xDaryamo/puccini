package normal

import (
	"github.com/tliron/puccini/tosca"
)

//
// FunctionCall
//

type FunctionCall struct {
	Key         interface{}   `json:"key,omitempty" yaml:"key,omitempty"`
	Description string        `json:"description" yaml:"description"`
	Constraints FunctionCalls `json:"constraints" yaml:"constraints"`

	FunctionCall *tosca.FunctionCall `json:"functionCall" yaml:"functionCall"`
}

func NewFunctionCall(functionCall *tosca.FunctionCall) *FunctionCall {
	return &FunctionCall{FunctionCall: functionCall}
}

// Constrainable interface
func (self *FunctionCall) SetKey(key interface{}) {
	self.Key = key
}

// Constrainable interface
func (self *FunctionCall) SetDescription(description string) {
	self.Description = description
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
