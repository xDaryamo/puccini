package tosca_v1_3

import (
	"github.com/tliron/puccini/tosca"
)

//
// RelationshipType
//
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.7.10
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.6.10
//

type RelationshipType struct {
	*Type `name:"relationship type"`

	PropertyDefinitions            PropertyDefinitions  `read:"properties,PropertyDefinition" inherit:"properties,Parent"`
	AttributeDefinitions           AttributeDefinitions `read:"attributes,AttributeDefinition" inherit:"attributes,Parent"`
	InterfaceDefinitions           InterfaceDefinitions `read:"interfaces,InterfaceDefinition" inherit:"interfaces,Parent"`
	ValidTargetCapabilityTypeNames *[]string            `read:"valid_target_types" inherit:"valid_target_types,Parent"`

	Parent                     *RelationshipType `lookup:"derived_from,ParentName" json:"-" yaml:"-"`
	ValidTargetCapabilityTypes CapabilityTypes   `lookup:"valid_target_types,ValidTargetCapabilityTypeNames" inherit:"valid_target_types,Parent" json:"-" yaml:"-"`
}

func NewRelationshipType(context *tosca.Context) *RelationshipType {
	return &RelationshipType{
		Type:                 NewType(context),
		PropertyDefinitions:  make(PropertyDefinitions),
		AttributeDefinitions: make(AttributeDefinitions),
		InterfaceDefinitions: make(InterfaceDefinitions),
	}
}

// tosca.Reader signature
func ReadRelationshipType(context *tosca.Context) interface{} {
	self := NewRelationshipType(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))
	return self
}

// tosca.Hierarchical interface
func (self *RelationshipType) GetParent() interface{} {
	return self.Parent
}

// tosca.Inherits interface
func (self *RelationshipType) Inherit() {
	log.Infof("{inherit} relationship type: %s", self.Name)

	if self.Parent == nil {
		return
	}

	self.PropertyDefinitions.Inherit(self.Parent.PropertyDefinitions)
	self.AttributeDefinitions.Inherit(self.Parent.AttributeDefinitions)
	self.InterfaceDefinitions.Inherit(self.Parent.InterfaceDefinitions)
}

//
// RelationshipTypes
//

type RelationshipTypes []*RelationshipType
