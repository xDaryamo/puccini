package normal

import (
	"github.com/tliron/puccini/tosca"
)

//
// FunctionCall
//

type FunctionCall struct {
	FunctionCall *tosca.FunctionCall `json:"functionCall" yaml:"functionCall"`
	Constraints  FunctionCalls       `json:"constraints" yaml:"constraints"`
	Description  string              `json:"description" yaml:"description"`
}

func NewFunctionCall(functionCall *tosca.FunctionCall) *FunctionCall {
	return &FunctionCall{FunctionCall: functionCall}
}

// Constrainable interface
func (self *FunctionCall) AddConstraint(constraint *tosca.FunctionCall) {
	self.Constraints = append(self.Constraints, NewFunctionCall(constraint))
}

// Constrainable interface
func (self *FunctionCall) SetDescription(description string) {
	self.Description = description
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
