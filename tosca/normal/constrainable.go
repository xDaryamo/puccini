package normal

import (
	"github.com/tliron/puccini/tosca"
)

//
// Constrainable
//

type Constrainable interface {
	SetKey(Constrainable)
	SetInformation(*ValueInformation)
	AddConstraint(*tosca.FunctionCall)
	SetConverter(*tosca.FunctionCall)
}

//
// Constrainables
//

type Constrainables map[any]Constrainable

//
// ConstrainableList
//

type ConstrainableList []Constrainable

func (self ConstrainableList) AppendWithKey(key any, value Constrainable) ConstrainableList {
	var constrainableKey Constrainable

	var ok bool
	if constrainableKey, ok = key.(Constrainable); !ok {
		constrainableKey = NewValue(key)
	}

	value.SetKey(constrainableKey)

	return append(self, value)
}
