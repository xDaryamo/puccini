package tosca_v2_0

import (
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
)

//
// RelationshipTemplate
//
// [TOSCA-v2.0] @ ?
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.8.4
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.8.4
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.7.4
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.7.4
//

type RelationshipTemplate struct {
	*Entity `name:"relationship template"`
	Name    string `namespace:""`

	CopyRelationshipTemplateName *string              `read:"copy"`
	RelationshipTypeName         *string              `read:"type" require:""`
	Metadata                     Metadata             `read:"metadata,Metadata"` // introduced in TOSCA 1.1
	Description                  *string              `read:"description"`
	Properties                   Values               `read:"properties,Value"`
	Attributes                   Values               `read:"attributes,AttributeValue"`
	Interfaces                   InterfaceAssignments `read:"interfaces,InterfaceAssignment"`

	CopyRelationshipTemplate *RelationshipTemplate `lookup:"copy,CopyRelationshipTemplateName" json:"-" yaml:"-"`
	RelationshipType         *RelationshipType     `lookup:"type,RelationshipTypeName" json:"-" yaml:"-"`

	rendered bool
}

func NewRelationshipTemplate(context *tosca.Context) *RelationshipTemplate {
	return &RelationshipTemplate{
		Entity:     NewEntity(context),
		Properties: make(Values),
		Attributes: make(Values),
		Name:       context.Name,
		Interfaces: make(InterfaceAssignments),
	}
}

// tosca.Reader signature
func ReadRelationshipTemplate(context *tosca.Context) tosca.EntityPtr {
	self := NewRelationshipTemplate(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))
	return self
}

// tosca.PreReadable interface
func (self *RelationshipTemplate) PreRead() {
	CopyTemplate(self.Context)
}

// parser.Renderable interface
func (self *RelationshipTemplate) Render() {
	logRender.Debugf("relationship template: %s", self.Name)

	if self.rendered {
		// Avoid rendering more than once (can happen if we were called from RelationshipAssignment.Render)
		return
	}
	self.rendered = true

	if self.RelationshipType == nil {
		return
	}

	self.Properties.RenderProperties(self.RelationshipType.PropertyDefinitions, "property", self.Context.FieldChild("properties", nil))
	self.Attributes.RenderAttributes(self.RelationshipType.AttributeDefinitions, self.Context.FieldChild("attributes", nil))
	self.Interfaces.Render(self.RelationshipType.InterfaceDefinitions, self.Context.FieldChild("interfaces", nil))
}

func (self *RelationshipTemplate) Normalize(normalRelationship *normal.Relationship) {
	normalRelationship.Metadata = self.Metadata

	if self.Description != nil {
		normalRelationship.Description = *self.Description
	}
}

//
// RelationshipTemplates
//

type RelationshipTemplates []*RelationshipTemplate
