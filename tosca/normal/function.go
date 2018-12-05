package normal

import (
	"github.com/tliron/puccini/tosca"
)

//
// Function
//

type Function struct {
	Function    *tosca.Function `json:"function" yaml:"function"`
	Constraints Functions       `json:"constraints" yaml:"constraints"`
	Description string          `json:"description" yaml:"description"`
}

func NewFunction(function *tosca.Function) *Function {
	return &Function{Function: function}
}

// Constrainable interface
func (self *Function) AddConstraint(constraint *tosca.Function) {
	self.Constraints = append(self.Constraints, NewFunction(constraint))
}

// Constrainable interface
func (self *Function) SetDescription(description string) {
	self.Description = description
}

//
// Functions
//

type Functions []*Function

//
// FunctionsMap
//

type FunctionsMap map[string]Functions

//
// FunctionsMapMap
//

type FunctionsMapMap map[string]FunctionsMap
