package normal

import (
	"github.com/tliron/puccini/tosca"
)

//
// Constrainable
//

type Constrainable interface {
	AddConstraint(*tosca.Function)
	SetDescription(string)
}

//
// Constrainables
//

type Constrainables map[string]Constrainable

//
// ConstrainableList
//

type ConstrainableList struct {
	List        []Constrainable `json:"list" yaml:"list"`
	Constraints Functions       `json:"constraints" yaml:"constraints"`
	Description string          `json:"description" yaml:"description"`
}

func NewConstrainableList(length int) *ConstrainableList {
	return &ConstrainableList{List: make([]Constrainable, length)}
}

// Constrainable interface
func (self *ConstrainableList) AddConstraint(constraint *tosca.Function) {
	self.Constraints = append(self.Constraints, NewFunction(constraint))
}

// Constrainable interface
func (self *ConstrainableList) SetDescription(description string) {
	self.Description = description
}

//
// ConstrainableMap
//

type ConstrainableMap struct {
	Map         Constrainables `json:"map" yaml:"map"`
	Constraints Functions      `json:"constraints" yaml:"constraints"`
	Description string         `json:"description" yaml:"description"`
}

func NewConstrainableMap() *ConstrainableMap {
	return &ConstrainableMap{Map: make(Constrainables)}
}

// Constrainable interface
func (self *ConstrainableMap) AddConstraint(constraint *tosca.Function) {
	self.Constraints = append(self.Constraints, NewFunction(constraint))
}

// Constrainable interface
func (self *ConstrainableMap) SetDescription(description string) {
	self.Description = description
}

// For access in JavaScript
func (self ConstrainableMap) Object(name string) map[string]interface{} {
	o := make(map[string]interface{})
	for key, coercible := range self.Map {
		o[key] = coercible
	}
	return o
}
