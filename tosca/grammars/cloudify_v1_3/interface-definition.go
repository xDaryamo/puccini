package cloudify_v1_3

import (
	"github.com/tliron/puccini/tosca"
)

//
// InterfaceDefinition
//
// [https://docs.cloudify.co/5.0.5/developer/blueprints/spec-interfaces/]
//

type InterfaceDefinition struct {
	*Entity `name:"interface definition"`
	Name    string

	OperationDefinitions OperationDefinitions `read:"?,OperationDefinition"`
}

func NewInterfaceDefinition(context *tosca.Context) *InterfaceDefinition {
	return &InterfaceDefinition{
		Entity:               NewEntity(context),
		Name:                 context.Name,
		OperationDefinitions: make(OperationDefinitions),
	}
}

// tosca.Reader signature
func ReadInterfaceDefinition(context *tosca.Context) tosca.EntityPtr {
	self := NewInterfaceDefinition(context)
	context.ReadFields(self)
	return self
}

// tosca.Mappable interface
func (self *InterfaceDefinition) GetKey() string {
	return self.Name
}

func (self *InterfaceDefinition) Inherit(parentDefinition *InterfaceDefinition) {
	logInherit.Debugf("interface definition: %s", self.Name)

	self.OperationDefinitions.Inherit(parentDefinition.OperationDefinitions)
}

//
// InterfaceDefinitions
//

type InterfaceDefinitions map[string]*InterfaceDefinition

func (self InterfaceDefinitions) Inherit(parentDefinitions InterfaceDefinitions) {
	for name, definition := range parentDefinitions {
		if _, ok := self[name]; !ok {
			self[name] = definition
		}
	}

	for name, definition := range self {
		if parentDefinition, ok := parentDefinitions[name]; ok {
			if definition != parentDefinition {
				definition.Inherit(parentDefinition)
			}
		}
	}
}
