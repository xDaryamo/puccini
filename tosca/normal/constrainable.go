package normal

import (
	"github.com/tliron/puccini/tosca"
)

//
// Constrainable
//

type Constrainable interface {
	SetKey(Constrainable)
	SetInformation(*Information)
	AddConstraint(*tosca.FunctionCall)
}

//
// Constrainables
//

type Constrainables map[interface{}]Constrainable

//
// ConstrainableList
//

type ConstrainableList []Constrainable

func (self ConstrainableList) AppendWithKey(key interface{}, value Constrainable) ConstrainableList {
	var constrainableKey Constrainable

	var ok bool
	if constrainableKey, ok = key.(Constrainable); !ok {
		constrainableKey = NewValue(key)
	}

	value.SetKey(constrainableKey)

	return append(self, value)
}
