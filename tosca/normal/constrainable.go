package normal

import (
	"encoding/json"

	"github.com/tliron/puccini/tosca"
)

//
// Constrainable
//

type Constrainable interface {
	SetKey(interface{})
	SetDescription(string)
	AddConstraint(*tosca.FunctionCall)
}

//
// Constrainables
//

type Constrainables map[interface{}]Constrainable

//
// MarshalableConstrainableMap
//

type MarshalableConstrainableMap map[interface{}]Constrainable

func (self MarshalableConstrainableMap) Marshalable() interface{} {
	var slice []Constrainable
	for key, constrainable := range self {
		constrainable.SetKey(key)
		slice = append(slice, constrainable)
	}
	return slice
}

// json.Marshaler interface
func (self MarshalableConstrainableMap) MarshalJSON() ([]byte, error) {
	return json.Marshal(self.Marshalable())
}

// yaml.Marshaler interface
func (self MarshalableConstrainableMap) MarshalYAML() (interface{}, error) {
	return self.Marshalable(), nil
}
