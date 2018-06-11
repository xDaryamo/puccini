package v1_1

import (
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
)

//
// RelationshipAssignment
//

type RelationshipAssignment struct {
	*Entity `name:"relationship"`

	RelationshipTemplateOrRelationshipTypeName *string              `read:"type"`
	Properties                                 Values               `read:"properties,Value"`
	Attributes                                 Values               `read:"attributes,Value"`
	Interfaces                                 InterfaceAssignments `read:"interfaces,InterfaceAssignment"`

	RelationshipTemplate *RelationshipTemplate `lookup:"type,RelationshipTemplateOrRelationshipTypeName" json:"-" yaml:"-"`
	RelationshipType     *RelationshipType     `lookup:"type,RelationshipTemplateOrRelationshipTypeName" json:"-" yaml:"-"`
}

func NewRelationshipAssignment(context *tosca.Context) *RelationshipAssignment {
	return &RelationshipAssignment{
		Entity:     NewEntity(context),
		Properties: make(Values),
		Attributes: make(Values),
		Interfaces: make(InterfaceAssignments),
	}
}

// tosca.Reader signature
func ReadRelationshipAssignment(context *tosca.Context) interface{} {
	self := NewRelationshipAssignment(context)
	if context.Is("map") {
		context.ValidateUnsupportedFields(context.ReadFields(self, Readers))
	} else if context.ValidateType("map", "string") {
		self.RelationshipTemplateOrRelationshipTypeName = context.ReadString()
	}
	return self
}

func init() {
	Readers["RelationshipAssignment"] = ReadRelationshipAssignment
}

func (self *RelationshipAssignment) Render(definition *RelationshipDefinition) {
	// TODO: could be relationship template

	// We will consider the "interfaces" at the definition to take priority over those at the type
	self.Interfaces.Render(definition.InterfaceDefinitions, self.Context.FieldChild("interfaces", nil))

	if self.RelationshipType != nil {
		self.Properties.RenderProperties(self.RelationshipType.PropertyDefinitions, "property", self.Context.FieldChild("properties", nil))
		self.Attributes.RenderAttributes(self.RelationshipType.AttributeDefinitions, self.Context.FieldChild("attributes", nil))
		self.Interfaces.Render(self.RelationshipType.InterfaceDefinitions, self.Context.FieldChild("interfaces", nil))
	}
}

func (self *RelationshipAssignment) Normalize(r *normal.Relationship) {
	if (self.RelationshipTemplate != nil) && (self.RelationshipTemplate.Description != nil) {
		r.Description = *self.RelationshipTemplate.Description
	} else if (self.RelationshipType != nil) && (self.RelationshipType.Description != nil) {
		r.Description = *self.RelationshipType.Description
	}

	if types, ok := normal.GetTypes(self.Context.Hierarchy, self.RelationshipType); ok {
		r.Types = types
	}

	self.Properties.Normalize(r.Properties)
	self.Attributes.Normalize(r.Attributes)

	for key, intr := range self.Interfaces {
		if definition, ok := intr.GetDefinitionForRelationship(self); ok {
			i := r.NewInterface(key)
			intr.Normalize(i, definition)
			r.Interfaces[key] = i
		}
	}
}
