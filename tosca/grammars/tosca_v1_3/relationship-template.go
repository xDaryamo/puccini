package tosca_v1_3

import (
	"github.com/tliron/puccini/tosca"
)

//
// RelationshipTemplate
//
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.8.4
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.7.4
//

type RelationshipTemplate struct {
	*Entity `name:"relationship template"`
	Name    string `namespace:""`

	CopyRelationshipTemplateName *string              `read:"copy"`
	RelationshipTypeName         *string              `read:"type" require:"type"`
	Description                  *string              `read:"description" inherit:"description,RelationshipType"`
	Properties                   Values               `read:"properties,Value"`
	Attributes                   Values               `read:"attributes,AttributeValue"`
	Interfaces                   InterfaceAssignments `read:"interfaces,InterfaceAssignment"`

	CopyRelationshipTemplate *RelationshipTemplate `lookup:"copy,CopyRelationshipTemplateName" json:"-" yaml:"-"`
	RelationshipType         *RelationshipType     `lookup:"type,RelationshipTypeName" json:"-" yaml:"-"`
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
func ReadRelationshipTemplate(context *tosca.Context) interface{} {
	self := NewRelationshipTemplate(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))
	return self
}

// tosca.Renderable interface
func (self *RelationshipTemplate) Render() {
	log.Infof("{render} relationship template: %s", self.Name)

	// TODO: copy
	if self.RelationshipType == nil {
		return
	}

	self.Properties.RenderProperties(self.RelationshipType.PropertyDefinitions, "property", self.Context.FieldChild("properties", nil))
	self.Attributes.RenderAttributes(self.RelationshipType.AttributeDefinitions, self.Context.FieldChild("attributes", nil))
	self.Interfaces.Render(self.RelationshipType.InterfaceDefinitions, self.Context.FieldChild("interfaces", nil))
}

//
// RelationshipTemplates
//

type RelationshipTemplates []*RelationshipTemplate
