package normal

import (
	"github.com/tliron/puccini/tosca"
)

//
// Map
//

type Map struct {
	Key              Constrainable `json:"key,omitempty" yaml:"key,omitempty"`
	Description      string        `json:"description,omitempty" yaml:"description,omitempty"`
	Constraints      FunctionCalls `json:"constraints,omitempty" yaml:"constraints,omitempty"`
	KeyDescription   string        `json:"keyDescription,omitempty" yaml:"keyDescription,omitempty"`
	KeyConstraints   FunctionCalls `json:"keyConstraints,omitempty" yaml:"keyConstraints,omitempty"`
	ValueDescription string        `json:"valueDescription,omitempty" yaml:"valueDescription,omitempty"`
	ValueConstraints FunctionCalls `json:"valueConstraints,omitempty" yaml:"valueConstraints,omitempty"`

	Entries ConstrainableList `json:"map" yaml:"map"`
}

func NewMap() *Map {
	return &Map{}
}

// Constrainable interface
func (self *Map) SetKey(key Constrainable) {
	self.Key = key
}

// Constrainable interface
func (self *Map) SetDescription(description string) {
	self.Description = description
}

// Constrainable interface
func (self *Map) AddConstraint(constraint *tosca.FunctionCall) {
	self.Constraints = append(self.Constraints, NewFunctionCall(constraint))
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
