package normal

import (
	"github.com/tliron/puccini/tosca"
)

//
// Value
//

type Value struct {
	Key         interface{}   `json:"key,omitempty" yaml:"key,omitempty"`
	Constraints FunctionCalls `json:"constraints" yaml:"constraints"`
	Description string        `json:"description" yaml:"description"`

	Value interface{} `json:"value" yaml:"value"` // can be ConstrainableList or ConstrainableMap
}

func NewValue(value interface{}) *Value {
	return &Value{Value: value}
}

// Constrainable interface
func (self *Value) SetKey(key interface{}) {
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
