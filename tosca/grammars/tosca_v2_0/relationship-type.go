package tosca_v2_0

import (
	"github.com/tliron/puccini/tosca/parsing"
)

//
// RelationshipType
//
// [TOSCA-v2.0] @ 7.3 Relationship Type
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.7.10
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.7.10
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.6.10
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.6.9
//

type RelationshipType struct {
	*Type `name:"relationship type"`

	PropertyDefinitions  PropertyDefinitions  `read:"properties,PropertyDefinition" inherit:"properties,Parent"`
	AttributeDefinitions AttributeDefinitions `read:"attributes,AttributeDefinition" inherit:"attributes,Parent"`
	InterfaceDefinitions InterfaceDefinitions `read:"interfaces,InterfaceDefinition" inherit:"interfaces,Parent"`

	// TOSCA 2.0: valid_capability_types
	ValidCapabilityTypeNames *[]string `read:"valid_capability_types" inherit:"valid_capability_types,Parent"`

	// TOSCA 1.3 backward compatibility: valid_target_types (maps to valid_capability_types in TOSCA 2.0)
	ValidTargetCapabilityTypeNames *[]string `read:"valid_target_types" inherit:"valid_target_types,Parent"`

	// TOSCA 2.0: New fields for target and source node type validation
	ValidTargetNodeTypeNames *[]string `read:"valid_target_node_types" inherit:"valid_target_node_types,Parent"`
	ValidSourceNodeTypeNames *[]string `read:"valid_source_node_types" inherit:"valid_source_node_types,Parent"`

	Parent                     *RelationshipType `lookup:"derived_from,ParentName" traverse:"ignore" json:"-" yaml:"-"`
	ValidCapabilityTypes       CapabilityTypes   `lookup:"valid_capability_types,ValidCapabilityTypeNames" inherit:"valid_capability_types,Parent" traverse:"ignore" json:"-" yaml:"-"`
	ValidTargetCapabilityTypes CapabilityTypes   `lookup:"valid_target_types,ValidTargetCapabilityTypeNames" inherit:"valid_target_types,Parent" traverse:"ignore" json:"-" yaml:"-"`
	ValidTargetNodeTypes       NodeTypes         `lookup:"valid_target_node_types,ValidTargetNodeTypeNames" inherit:"valid_target_node_types,Parent" traverse:"ignore" json:"-" yaml:"-"`
	ValidSourceNodeTypes       NodeTypes         `lookup:"valid_source_node_types,ValidSourceNodeTypeNames" inherit:"valid_source_node_types,Parent" traverse:"ignore" json:"-" yaml:"-"`
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

	// TOSCA 2.0 backward compatibility: Map valid_target_types to valid_capability_types if needed
	// If TOSCA 1.3 style valid_target_types is used but valid_capability_types is not specified,
	// copy the valid_target_types to valid_capability_types for unified processing
	if (self.ValidTargetCapabilityTypeNames != nil) && (self.ValidCapabilityTypeNames == nil) {
		self.ValidCapabilityTypeNames = self.ValidTargetCapabilityTypeNames
	}

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

	// TOSCA 2.0: Inherit validation rules for capability types, target node types, and source node types
	// Note: Derived types can only further restrict, not expand these lists

	// Handle backward compatibility for valid_target_types -> valid_capability_types mapping in parent
	if (self.Parent.ValidTargetCapabilityTypeNames != nil) && (self.Parent.ValidCapabilityTypeNames == nil) {
		if (self.ValidTargetCapabilityTypeNames != nil) && (self.ValidCapabilityTypeNames == nil) {
			self.ValidCapabilityTypeNames = self.ValidTargetCapabilityTypeNames
		}
	}
}

//
// RelationshipTypes
//

type RelationshipTypes []*RelationshipType
