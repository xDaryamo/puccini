package normal

import (
	"github.com/tliron/puccini/ard"
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/yamlkeys"
)

//
// Map
//

type Map struct {
	Key         interface{}   `json:"key,omitempty" yaml:"key,omitempty"`
	Constraints FunctionCalls `json:"constraints" yaml:"constraints"`
	Description string        `json:"description" yaml:"description"`

	Map MarshalableConstrainableMap `json:"map" yaml:"map"`
}

func NewMap() *Map {
	return &Map{Map: make(MarshalableConstrainableMap)}
}

// Constrainable interface
func (self *Map) SetKey(key interface{}) {
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

// For access in JavaScript
func (self Map) Object(name string) map[string]interface{} {
	// JavaScript requires keys to be strings, so we would lose complex keys
	o := make(ard.StringMap)
	for key, constrainable := range self.Map {
		o[yamlkeys.KeyString(key)] = constrainable
	}
	return o
}
