package tosca_v1_3

import (
	"github.com/tliron/puccini/tosca"
)

//
// CapabilityDefinition
//
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.7.2
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.6.2
//

type CapabilityDefinition struct {
	*Entity `name:"capability definition"`
	Name    string

	Description              *string              `read:"description" inherit:"description,CapabilityType"`
	CapabilityTypeName       *string              `read:"type"` // required only if cannot be inherited
	PropertyDefinitions      PropertyDefinitions  `read:"properties,PropertyDefinition" inherit:"properties,CapabilityType"`
	AttributeDefinitions     AttributeDefinitions `read:"attributes,PropertyDefinition" inherit:"attributes,CapabilityType"`
	ValidSourceNodeTypeNames *[]string            `read:"valid_source_types" inherit:"valid_source_types,CapabilityType"`
	Occurrences              *RangeEntity         `read:"occurrences,RangeEntity"`

	CapabilityType       *CapabilityType `lookup:"type,CapabilityTypeName" json:"-" yaml:"-"`
	ValidSourceNodeTypes NodeTypes       `lookup:"valid_source_types,ValidSourceNodeTypeNames" apply:"valid_source_types,CapabilityType" json:"-" yaml:"-"`

	typeMissingProblemReported bool
}

func NewCapabilityDefinition(context *tosca.Context) *CapabilityDefinition {
	return &CapabilityDefinition{
		Entity:               NewEntity(context),
		Name:                 context.Name,
		PropertyDefinitions:  make(PropertyDefinitions),
		AttributeDefinitions: make(AttributeDefinitions),
	}
}

// tosca.Reader signature
func ReadCapabilityDefinition(context *tosca.Context) interface{} {
	self := NewCapabilityDefinition(context)

	if context.Is("map") {
		// Long notation
		context.ValidateUnsupportedFields(context.ReadFields(self))
	} else if context.ValidateType("map", "string") {
		// Short notation
		self.CapabilityTypeName = context.FieldChild("type", context.Data).ReadString()
	}

	return self
}

// tosca.Mappable interface
func (self *CapabilityDefinition) GetKey() string {
	return self.Name
}

func (self *CapabilityDefinition) Inherit(parentDefinition *CapabilityDefinition) {
	if parentDefinition != nil {
		if ((self.Description == nil) || ((self.CapabilityType != nil) && (self.Description == self.CapabilityType.Description))) && (parentDefinition.Description != nil) {
			self.Description = parentDefinition.Description
		}
		if (self.CapabilityTypeName == nil) && (parentDefinition.CapabilityTypeName != nil) {
			self.CapabilityTypeName = parentDefinition.CapabilityTypeName
		}
		if (self.ValidSourceNodeTypeNames == nil) && (parentDefinition.ValidSourceNodeTypeNames != nil) {
			self.ValidSourceNodeTypeNames = parentDefinition.ValidSourceNodeTypeNames
		}
		if (self.Occurrences == nil) && (parentDefinition.Occurrences != nil) {
			self.Occurrences = parentDefinition.Occurrences
		}
		if (self.CapabilityType == nil) && (parentDefinition.CapabilityType != nil) {
			self.CapabilityType = parentDefinition.CapabilityType
		}
		if (self.ValidSourceNodeTypes == nil) && (parentDefinition.ValidSourceNodeTypes != nil) {
			self.ValidSourceNodeTypes = parentDefinition.ValidSourceNodeTypes
		}

		// Validate type compatibility
		if (self.CapabilityType != nil) && (parentDefinition.CapabilityType != nil) && !self.Context.Hierarchy.IsCompatible(parentDefinition.CapabilityType, self.CapabilityType) {
			self.Context.ReportIncompatibleType(self.CapabilityType.Name, parentDefinition.CapabilityType.Name)
			return
		}

		self.PropertyDefinitions.Inherit(parentDefinition.PropertyDefinitions)
		self.AttributeDefinitions.Inherit(parentDefinition.AttributeDefinitions)
	} else {
		self.PropertyDefinitions.Inherit(nil)
		self.AttributeDefinitions.Inherit(nil)
	}

	if self.CapabilityTypeName == nil {
		// Avoid reporting more than once
		if !self.typeMissingProblemReported {
			self.Context.FieldChild("type", nil).ReportFieldMissing()
			self.typeMissingProblemReported = true
		}
	}
}

//
// CapabilityDefinitions
//

type CapabilityDefinitions map[string]*CapabilityDefinition

func (self CapabilityDefinitions) Inherit(parentDefinitions CapabilityDefinitions) {
	for name, definition := range parentDefinitions {
		if _, ok := self[name]; !ok {
			self[name] = definition
		}
	}

	for name, definition := range self {
		if parentDefinition, ok := parentDefinitions[name]; ok {
			if definition != parentDefinition {
				definition.Inherit(parentDefinition)
			}
		} else {
			definition.Inherit(nil)
		}
	}
}
