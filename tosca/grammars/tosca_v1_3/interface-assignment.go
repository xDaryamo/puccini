package tosca_v1_3

import (
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
)

//
// InterfaceAssignment
//
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.16
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.14
//

type InterfaceAssignment struct {
	*Entity `name:"interface" json:"-" yaml:"-"`
	Name    string

	Inputs        Values                  `read:"inputs,Value"`
	Operations    OperationAssignments    `read:"operations,OperationAssignment"`
	Notifications NotificationAssignments `read:"notifications,NotificationAssignment"`
}

func NewInterfaceAssignment(context *tosca.Context) *InterfaceAssignment {
	return &InterfaceAssignment{
		Entity:        NewEntity(context),
		Name:          context.Name,
		Inputs:        make(Values),
		Operations:    make(OperationAssignments),
		Notifications: make(NotificationAssignments),
	}
}

// tosca.Reader signature
func ReadInterfaceAssignment(context *tosca.Context) interface{} {
	self := NewInterfaceAssignment(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))
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
	self.Operations.Normalize(i)
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

	// TODO: notifications
}

func (self InterfaceAssignments) NormalizeForNodeTemplate(nodeTemplate *NodeTemplate, n *normal.NodeTemplate) {
	for key, intr := range self {
		if definition, ok := intr.GetDefinitionForNodeTemplate(nodeTemplate); ok {
			intr.Normalize(n.NewInterface(key), definition)
		}
	}
}

func (self InterfaceAssignments) NormalizeForGroup(group *Group, g *normal.Group) {
	for key, intr := range self {
		if definition, ok := intr.GetDefinitionForGroup(group); ok {
			intr.Normalize(g.NewInterface(key), definition)
		}
	}
}

func (self InterfaceAssignments) NormalizeForRelationship(relationship *RelationshipAssignment, r *normal.Relationship) {
	for key, intr := range self {
		if definition, ok := intr.GetDefinitionForRelationship(relationship); ok {
			intr.Normalize(r.NewInterface(key), definition)
		}
	}
}
