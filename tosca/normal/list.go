package normal

import (
	"github.com/tliron/puccini/tosca"
)

//
// List
//

type List struct {
	Key         Constrainable `json:"$key,omitempty" yaml:"$key,omitempty"`
	Information *Information  `json:"$information,omitempty" yaml:"$information,omitempty"`
	Constraints FunctionCalls `json:"$constraints,omitempty" yaml:"$constraints,omitempty"`

	EntryConstraints FunctionCalls `json:"$entryConstraints,omitempty" yaml:"$entryConstraints,omitempty"`

	Entries ConstrainableList `json:"$list" yaml:"$list"`
}

func NewList(length int) *List {
	return &List{Entries: make(ConstrainableList, length)}
}

// Constrainable interface
func (self *List) SetKey(key Constrainable) {
	self.Key = key
}

// Constrainable interface
func (self *List) SetInformation(information *Information) {
	self.Information = CopyInformation(information)
}

// Constrainable interface
func (self *List) AddConstraint(functionCall *tosca.FunctionCall) {
	self.Constraints = append(self.Constraints, NewFunctionCall(functionCall))
}

func (self *List) AddEntryConstraint(constraint *tosca.FunctionCall) {
	self.EntryConstraints = append(self.EntryConstraints, NewFunctionCall(constraint))
}

func (self *List) Set(index int, value Constrainable) {
	self.Entries[index] = value
}
