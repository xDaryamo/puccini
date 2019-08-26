package normal

import (
	"github.com/tliron/puccini/tosca"
)

//
// Value
//

type Value struct {
	Value       interface{}   `json:"value" yaml:"value"` // can be list or map of Coercibles
	Constraints FunctionCalls `json:"constraints" yaml:"constraints"`
	Description string        `json:"description" yaml:"description"`
}

func NewValue(value interface{}) *Value {
	return &Value{Value: value}
}

// Constrainable interface
func (self *Value) AddConstraint(constraint *tosca.FunctionCall) {
	self.Constraints = append(self.Constraints, NewFunctionCall(constraint))
}

// Constrainable interface
func (self *Value) SetDescription(description string) {
	self.Description = description
}
