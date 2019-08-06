package tosca_v1_3

import (
	"math"

	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
)

//
// CapabilityAssignment
//
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.8.1
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.7.1
//

type CapabilityAssignment struct {
	*Entity `name:"capability"`
	Name    string

	Properties Values `read:"properties,Value"`
	Attributes Values `read:"attributes,AttributeValue"`
}

func NewCapabilityAssignment(context *tosca.Context) *CapabilityAssignment {
	return &CapabilityAssignment{
		Entity:     NewEntity(context),
		Name:       context.Name,
		Properties: make(Values),
		Attributes: make(Values),
	}
}

// tosca.Reader signature
func ReadCapabilityAssignment(context *tosca.Context) interface{} {
	self := NewCapabilityAssignment(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))
	return self
}

// tosca.Mappable interface
func (self *CapabilityAssignment) GetKey() string {
	return self.Name
}

func (self *CapabilityAssignment) GetDefinition(nodeTemplate *NodeTemplate) (*CapabilityDefinition, bool) {
	if nodeTemplate.NodeType == nil {
		return nil, false
	}
	definition, ok := nodeTemplate.NodeType.CapabilityDefinitions[self.Name]
	return definition, ok
}

func (self *CapabilityAssignment) Normalize(n *normal.NodeTemplate, definition *CapabilityDefinition) *normal.Capability {
	log.Debugf("{normalize} capability: %s", self.Name)

	c := n.NewCapability(self.Name)

	if definition.Description != nil {
		c.Description = *definition.Description
	}

	if definition.Occurrences != nil {
		c.MinRelationshipCount = definition.Occurrences.Range.Lower
		c.MaxRelationshipCount = definition.Occurrences.Range.Upper
	} else {
		// Default occurrences is [ 0, UNBOUNDED ]
		c.MinRelationshipCount = 0
		c.MaxRelationshipCount = math.MaxUint64
	}

	if types, ok := normal.GetTypes(self.Context.Hierarchy, definition.CapabilityType); ok {
		c.Types = types
	}

	self.Properties.Normalize(c.Properties)
	self.Attributes.Normalize(c.Attributes)

	return c
}

//
// CapabilityAssignments
//

type CapabilityAssignments map[string]*CapabilityAssignment

func (self *CapabilityAssignment) Render(definition *CapabilityDefinition) {
	self.Properties.RenderProperties(definition.PropertyDefinitions, "property", self.Context.FieldChild("properties", nil))
	self.Attributes.RenderAttributes(definition.AttributeDefinitions, self.Context.FieldChild("attributes", nil))
}

func (self CapabilityAssignments) Render(definitions CapabilityDefinitions, context *tosca.Context) {
	for key, definition := range definitions {
		assignment, ok := self[key]
		if !ok {
			assignment = NewCapabilityAssignment(context.MapChild(key, nil))
			self[key] = assignment
		}
		assignment.Render(definition)
	}

	for key, assignment := range self {
		_, ok := definitions[key]
		if !ok {
			assignment.Context.ReportUndefined("capability")
			delete(self, key)
		}
	}
}

func (self CapabilityAssignments) Normalize(nodeTemplate *NodeTemplate, n *normal.NodeTemplate) {
	for key, capability := range self {
		if definition, ok := capability.GetDefinition(nodeTemplate); ok {
			n.Capabilities[key] = capability.Normalize(n, definition)
		}
	}
}
