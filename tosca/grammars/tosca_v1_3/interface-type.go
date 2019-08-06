package tosca_v1_3

import (
	"github.com/tliron/puccini/tosca"
)

//
// InterfaceType
//
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.7.5
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.6.5
//

type InterfaceType struct {
	*Type `name:"interface type"`

	InputDefinitions     PropertyDefinitions  `read:"inputs,PropertyDefinition" inherit:"inputs,Parent"`
	OperationDefinitions OperationDefinitions `read:"?,OperationDefinition" inherit:"?,Parent"`

	Parent *InterfaceType `lookup:"derived_from,ParentName" json:"-" yaml:"-"`
}

func NewInterfaceType(context *tosca.Context) *InterfaceType {
	return &InterfaceType{
		Type:                 NewType(context),
		InputDefinitions:     make(PropertyDefinitions),
		OperationDefinitions: make(OperationDefinitions),
	}
}

// tosca.Reader signature
func ReadInterfaceType(context *tosca.Context) interface{} {
	self := NewInterfaceType(context)
	context.ReadFields(self)
	return self
}

// tosca.Hierarchical interface
func (self *InterfaceType) GetParent() interface{} {
	return self.Parent
}

// tosca.Inherits interface
func (self *InterfaceType) Inherit() {
	log.Infof("{inherit} interface type: %s", self.Name)

	if self.Parent == nil {
		return
	}

	self.InputDefinitions.Inherit(self.Parent.InputDefinitions)
	self.OperationDefinitions.Inherit(self.Parent.OperationDefinitions)
}

//
// InterfaceTypes
//

type InterfaceTypes []*InterfaceType
