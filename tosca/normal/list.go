package normal

import (
	"github.com/tliron/puccini/tosca"
)

//
// List
//

type List struct {
	Key         interface{}   `json:"key,omitempty" yaml:"key,omitempty"`
	Description string        `json:"description" yaml:"description"`
	Constraints FunctionCalls `json:"constraints" yaml:"constraints"`

	List []Constrainable `json:"list" yaml:"list"`
}

func NewList(length int) *List {
	return &List{List: make([]Constrainable, length)}
}

// Constrainable interface
func (self *List) SetKey(key interface{}) {
	self.Key = key
}

// Constrainable interface
func (self *List) SetDescription(description string) {
	self.Description = description
}

// Constrainable interface
func (self *List) AddConstraint(functionCall *tosca.FunctionCall) {
	self.Constraints = append(self.Constraints, NewFunctionCall(functionCall))
}
