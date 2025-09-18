package tosca_v2_0

import (
	"github.com/tliron/puccini/tosca/parsing"
)

//
// RelationshipType
//
// [TOSCA-v2.0] @ 7.3 Relationship Type
//

type RelationshipType struct {
	*Type `name:"relationship type"`

	PropertyDefinitions  PropertyDefinitions  `read:"properties,PropertyDefinition" inherit:"properties,Parent"`
	AttributeDefinitions AttributeDefinitions `read:"attributes,AttributeDefinition" inherit:"attributes,Parent"`
	InterfaceDefinitions InterfaceDefinitions `read:"interfaces,InterfaceDefinition" inherit:"interfaces,Parent"`

	// TOSCA 2.0: valid_capability_types
	ValidCapabilityTypeNames *[]string `read:"valid_capability_types" inherit:"valid_capability_types,Parent"`

	// TOSCA 2.0: New fields for target and source node type validation
	ValidTargetNodeTypeNames *[]string `read:"valid_target_node_types" inherit:"valid_target_node_types,Parent"`
	ValidSourceNodeTypeNames *[]string `read:"valid_source_node_types" inherit:"valid_source_node_types,Parent"`

	Parent               *RelationshipType `lookup:"derived_from,ParentName" traverse:"ignore" json:"-" yaml:"-"`
	ValidCapabilityTypes CapabilityTypes   `lookup:"valid_capability_types,ValidCapabilityTypeNames" inherit:"valid_capability_types,Parent" traverse:"ignore" json:"-" yaml:"-"`
	ValidTargetNodeTypes NodeTypes         `lookup:"valid_target_node_types,ValidTargetNodeTypeNames" inherit:"valid_target_node_types,Parent" traverse:"ignore" json:"-" yaml:"-"`
	ValidSourceNodeTypes NodeTypes         `lookup:"valid_source_node_types,ValidSourceNodeTypeNames" inherit:"valid_source_node_types,Parent" traverse:"ignore" json:"-" yaml:"-"`
}

func NewRelationshipType(context *parsing.Context) *RelationshipType {
	return &RelationshipType{
		Type:                 NewType(context),
		PropertyDefinitions:  make(PropertyDefinitions),
		AttributeDefinitions: make(AttributeDefinitions),
		InterfaceDefinitions: make(InterfaceDefinitions),
	}
}

// ([parsing.Reader] signature)
func ReadRelationshipType(context *parsing.Context) parsing.EntityPtr {
	self := NewRelationshipType(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))
	return self
}

// ([parsing.Hierarchical] interface)
func (self *RelationshipType) GetParent() parsing.EntityPtr {
	return self.Parent
}

// ([parsing.Inherits] interface)
func (self *RelationshipType) Inherit() {
	logInherit.Debugf("relationship type: %s", self.Name)

	if self.Parent == nil {
		return
	}

	self.PropertyDefinitions.Inherit(self.Parent.PropertyDefinitions)
	self.AttributeDefinitions.Inherit(self.Parent.AttributeDefinitions)
	self.InterfaceDefinitions.Inherit(self.Parent.InterfaceDefinitions)

	// Inherit validation rules for capability types, target node types, and source node types
	// Note: Derived types can only further restrict, not expand these lists
}

//
// RelationshipTypes
//

type RelationshipTypes []*RelationshipType
