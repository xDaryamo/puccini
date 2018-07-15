package v1_1

import (
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
)

//
// InterfaceAssignment
//

type InterfaceAssignment struct {
	*Entity `name:"interface" json:"-" yaml:"-"`
	Name    string

	Inputs     Values               `read:"inputs,Value"`
	Operations OperationAssignments `read:"?,OperationAssignment"`
}

func NewInterfaceAssignment(context *tosca.Context) *InterfaceAssignment {
	return &InterfaceAssignment{
		Entity:     NewEntity(context),
		Name:       context.Name,
		Inputs:     make(Values),
		Operations: make(OperationAssignments),
	}
}

// tosca.Reader signature
func ReadInterfaceAssignment(context *tosca.Context) interface{} {
	self := NewInterfaceAssignment(context)
	context.ReadFields(self, Readers)
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

func (self *InterfaceAssignment) GetDefinitionForGroup(group *Group) (*InterfaceDefinition, bool) {
	if group.GroupType == nil {
		return nil, false
	}
	definition, ok := group.GroupType.InterfaceDefinitions[self.Name]
	return definition, ok
}

func (self *InterfaceAssignment) GetDefinitionForRelationship(relationship *RelationshipAssignment) (*InterfaceDefinition, bool) {
	if relationship.RelationshipType == nil {
		return nil, false
	}
	definition, ok := relationship.RelationshipType.InterfaceDefinitions[self.Name]
	return definition, ok
}

func (self *InterfaceAssignment) Render(definition *InterfaceDefinition) {
	self.Inputs.RenderProperties(definition.InputDefinitions, "input", self.Context.FieldChild("inputs", nil))
	self.Operations.Render(definition.OperationDefinitions, self.Context)
}

func (self *InterfaceAssignment) Normalize(i *normal.Interface, definition *InterfaceDefinition) {
	log.Debugf("{normalize} interface: %s", self.Name)

	if (definition.InterfaceType != nil) && (definition.InterfaceType.Description != nil) {
		i.Description = *definition.InterfaceType.Description
	}

	if types, ok := normal.GetTypes(self.Context.Hierarchy, definition.InterfaceType); ok {
		i.Types = types
	}

	self.Inputs.Normalize(i.Inputs)

	for key, operation := range self.Operations {
		i.Operations[key] = operation.Normalize(i)
	}
}

//
// InterfaceAssignments
//

type InterfaceAssignments map[string]*InterfaceAssignment

func (self InterfaceAssignments) Render(definitions InterfaceDefinitions, context *tosca.Context) {
	for key, definition := range definitions {
		assignment, ok := self[key]
		if !ok {
			assignment = NewInterfaceAssignment(context.MapChild(key, nil))
			self[key] = assignment
		}
		assignment.Render(definition)
	}

	for key, assignment := range self {
		_, ok := definitions[key]
		if !ok {
			assignment.Context.ReportUndefined("interface")
			delete(self, key)
		}
	}
}
