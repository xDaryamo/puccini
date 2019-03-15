package cloudify_v1_3

import (
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
)

//
// InterfaceAssignment
//
// [https://docs.cloudify.co/4.5.5/developer/blueprints/spec-interfaces/]
//

type InterfaceAssignment struct {
	*Entity `name:"interface assignment"`
	Name    string

	Operations OperationAssignments `read:"?,OperationAssignment"`
}

func NewInterfaceAssignment(context *tosca.Context) *InterfaceAssignment {
	return &InterfaceAssignment{
		Entity: NewEntity(context),
		Name:   context.Name,
	}
}

// tosca.Reader signature
func ReadInterfaceAssignment(context *tosca.Context) interface{} {
	self := NewInterfaceAssignment(context)
	context.ReadFields(self)
	return self
}

// tosca.Mappable interface
func (self *InterfaceAssignment) GetKey() string {
	return self.Name
}

func (self *InterfaceAssignment) GetDefinitionForNodeTemplate(nodeTemplate *NodeTemplate) (*InterfaceDefinition, bool) {
	if nodeTemplate.NodeType == nil {
		return nil, false
	}
	definition, ok := nodeTemplate.NodeType.InterfaceDefinitions[self.Name]
	return definition, ok
}

func (self *InterfaceAssignment) Normalize(i *normal.Interface, definition *InterfaceDefinition) {
	log.Debugf("{normalize} interface: %s", self.Name)

	for key, operation := range self.Operations {
		i.Operations[key] = operation.Normalize(i)
	}
}

//
// InterfaceAssignments
//

type InterfaceAssignments map[string]*InterfaceAssignment
