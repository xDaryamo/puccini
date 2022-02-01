package tosca_v2_0

import (
	"github.com/tliron/puccini/tosca"
)

//
// InterfaceType
//
// [TOSCA-v2.0] @ ?
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.7.5
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.7.5
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.6.5
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.6.4
//

type InterfaceType struct {
	*Type `name:"interface type"`

	InputDefinitions        PropertyDefinitions     `read:"inputs,PropertyDefinition" inherit:"inputs,Parent"`
	OperationDefinitions    OperationDefinitions    `read:"operations,OperationDefinition" inherit:"operations,Parent"`
	NotificationDefinitions NotificationDefinitions `read:"notifications,NotificationDefinition" inherit:"notifications,Parent"` // introduced in TOSCA 1.3

	Parent *InterfaceType `lookup:"derived_from,ParentName" json:"-" yaml:"-"`
}

func NewInterfaceType(context *tosca.Context) *InterfaceType {
	return &InterfaceType{
		Type:                    NewType(context),
		InputDefinitions:        make(PropertyDefinitions),
		OperationDefinitions:    make(OperationDefinitions),
		NotificationDefinitions: make(NotificationDefinitions),
	}
}

// tosca.Reader signature
func ReadInterfaceType(context *tosca.Context) tosca.EntityPtr {
	self := NewInterfaceType(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))
	return self
}

// tosca.Hierarchical interface
func (self *InterfaceType) GetParent() tosca.EntityPtr {
	return self.Parent
}

// tosca.Inherits interface
func (self *InterfaceType) Inherit() {
	logInherit.Debugf("interface type: %s", self.Name)

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
