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
		Entity:     NewEntity(context),
		Name:       context.Name,
		Operations: make(OperationAssignments),
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

func (self *InterfaceAssignment) GetDefinitionForRelationshipSource(relationship *RelationshipAssignment) (*InterfaceDefinition, bool) {
	if relationship.RelationshipType == nil {
		return nil, false
	}
	definition, ok := relationship.RelationshipType.SourceInterfaceDefinitions[self.Name]
	return definition, ok
}

func (self *InterfaceAssignment) GetDefinitionForRelationshipTarget(relationship *RelationshipAssignment) (*InterfaceDefinition, bool) {
	if relationship.RelationshipType == nil {
		return nil, false
	}
	definition, ok := relationship.RelationshipType.TargetInterfaceDefinitions[self.Name]
	return definition, ok
}

func (self *InterfaceAssignment) Render(definition *InterfaceDefinition) {
	self.Operations.Render(definition.OperationDefinitions, self.Context)
}

func (self *InterfaceAssignment) Normalize(i *normal.Interface, definition *InterfaceDefinition) {
	log.Debugf("{normalize} interface: %s", self.Name)
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
			assignment.Context.ReportUndeclared("interface")
			delete(self, key)
		}
	}
}

func (self InterfaceAssignments) NormalizeForNodeTemplate(nodeTemplate *NodeTemplate, n *normal.NodeTemplate) {
	for key, intr := range self {
		if definition, ok := intr.GetDefinitionForNodeTemplate(nodeTemplate); ok {
			intr.Normalize(n.NewInterface(key), definition)
		}
	}
}

func (self InterfaceAssignments) NormalizeForRelationshipSource(relationship *RelationshipAssignment, r *normal.Relationship) {
	for key, intr := range self {
		if definition, ok := intr.GetDefinitionForRelationshipSource(relationship); ok {
			i := r.NewInterface(key)
			intr.Normalize(i, definition)
			for _, operation := range i.Operations {
				operation.Host = "SOURCE"
			}
		}
	}
}

func (self InterfaceAssignments) NormalizeForRelationshipTarget(relationship *RelationshipAssignment, r *normal.Relationship) {
	for key, intr := range self {
		if definition, ok := intr.GetDefinitionForRelationshipTarget(relationship); ok {
			i := r.NewInterface(key)
			intr.Normalize(i, definition)
			for _, operation := range i.Operations {
				operation.Host = "TARGET"
			}
		}
	}
}
