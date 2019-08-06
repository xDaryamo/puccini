package tosca_v1_3

import (
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
)

//
// RelationshipAssignment
//
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.8.2.2.3
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.7.2.2.3
//

type RelationshipAssignment struct {
	*Entity `name:"relationship"`

	RelationshipTemplateNameOrTypeName *string              `read:"type"`
	Properties                         Values               `read:"properties,Value"`
	Attributes                         Values               `read:"attributes,AttributeValue"` // missing in spec
	Interfaces                         InterfaceAssignments `read:"interfaces,InterfaceAssignment"`

	RelationshipTemplate *RelationshipTemplate `lookup:"type,RelationshipTemplateNameOrTypeName" json:"-" yaml:"-"`
	RelationshipType     *RelationshipType     `lookup:"type,RelationshipTemplateNameOrTypeName" json:"-" yaml:"-"`
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
		// Long notation
		context.ValidateUnsupportedFields(context.ReadFields(self))
	} else if context.ValidateType("map", "string") {
		// Short notation
		self.RelationshipTemplateNameOrTypeName = context.FieldChild("type", context.Data).ReadString()
	}

	return self
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
	self.Interfaces.NormalizeForRelationship(self, r)
}
