package cloudify_v1_3

import (
	"github.com/tliron/puccini/tosca"
)

//
// InterfaceDefinition
//
// [https://docs.cloudify.co/4.5.5/developer/blueprints/spec-interfaces/]
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
func ReadInterfaceDefinition(context *tosca.Context) interface{} {
	self := NewInterfaceDefinition(context)
	context.ReadFields(self)
	return self
}

// tosca.Mappable interface
func (self *InterfaceDefinition) GetKey() string {
	return self.Name
}

//
// InterfaceDefinitions
//

type InterfaceDefinitions map[string]*InterfaceDefinition
