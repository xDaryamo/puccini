package tosca_v2_0

import (
	"github.com/tliron/go-ard"
	"github.com/tliron/puccini/tosca/parsing"
)

//
// CapabilityDefinition
//
// [TOSCA-v2.0] @ 8.2
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.7.2
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.7.2
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.6.2
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.6.1
//

type CapabilityDefinition struct {
	*Entity `name:"capability definition"`
	Name    string

	Description                *string              `read:"description"`
	Metadata                   Metadata             `read:"metadata,Metadata"`
	CapabilityTypeName         *string              `read:"type"` // mandatory only if cannot be inherited
	PropertyDefinitions        PropertyDefinitions  `read:"properties,PropertyDefinition" inherit:"properties,CapabilityType"`
	AttributeDefinitions       AttributeDefinitions `read:"attributes,AttributeDefinition" inherit:"attributes,CapabilityType"`
	ValidSourceNodeTypeNames   *[]string            `read:"valid_source_node_types" inherit:"valid_source_node_types,CapabilityType"`
	ValidRelationshipTypeNames *[]string            `read:"valid_relationship_types" inherit:"valid_relationship_types,CapabilityType"`

	CapabilityType         *CapabilityType   `lookup:"type,CapabilityTypeName" traverse:"ignore" json:"-" yaml:"-"`
	ValidSourceNodeTypes   NodeTypes         `lookup:"valid_source_node_types,ValidSourceNodeTypeNames" traverse:"ignore" json:"-" yaml:"-"`
	ValidRelationshipTypes RelationshipTypes `lookup:"valid_relationship_types,ValidRelationshipTypeNames" traverse:"ignore" json:"-" yaml:"-"`

	typeMissingProblemReported bool
}

func NewCapabilityDefinition(context *parsing.Context) *CapabilityDefinition {
	return &CapabilityDefinition{
		Entity:               NewEntity(context),
		Name:                 context.Name,
		PropertyDefinitions:  make(PropertyDefinitions),
		AttributeDefinitions: make(AttributeDefinitions),
	}
}

// ([parsing.Reader] signature)
func ReadCapabilityDefinition(context *parsing.Context) parsing.EntityPtr {
	self := NewCapabilityDefinition(context)

	if context.Is(ard.TypeMap) {
		// Long notation
		context.ValidateUnsupportedFields(context.ReadFields(self))
	} else if context.ValidateType(ard.TypeMap, ard.TypeString) {
		// Short notation
		self.CapabilityTypeName = context.FieldChild("type", context.Data).ReadString()
	}

	return self
}

// ([parsing.Mappable] interface)
func (self *CapabilityDefinition) GetKey() string {
	return self.Name
}

func (self *CapabilityDefinition) Inherit(parentDefinition *CapabilityDefinition) {
	logInherit.Debugf("capability definition: %s", self.Name)

	// Validate type compatibility
	if (self.CapabilityType != nil) && (parentDefinition.CapabilityType != nil) && !self.Context.Hierarchy.IsCompatible(parentDefinition.CapabilityType, self.CapabilityType) {
		self.Context.ReportIncompatibleType(self.CapabilityType, parentDefinition.CapabilityType)
		return
	}

	if ((self.Description == nil) || ((self.CapabilityType != nil) && (self.Description == self.CapabilityType.Description))) && (parentDefinition.Description != nil) {
		self.Description = parentDefinition.Description
	}
	
	// For capability refinement: inherit type if not specified
	if (self.CapabilityTypeName == nil) && (parentDefinition.CapabilityTypeName != nil) {
		self.CapabilityTypeName = parentDefinition.CapabilityTypeName
		// Mark this as inherited so render() knows not to complain
		self.typeMissingProblemReported = true
	}
	
	if (self.ValidSourceNodeTypeNames == nil) && (parentDefinition.ValidSourceNodeTypeNames != nil) {
		self.ValidSourceNodeTypeNames = parentDefinition.ValidSourceNodeTypeNames
	}
	if (self.ValidRelationshipTypeNames == nil) && (parentDefinition.ValidRelationshipTypeNames != nil) {
		self.ValidRelationshipTypeNames = parentDefinition.ValidRelationshipTypeNames
	}
	if (self.CapabilityType == nil) && (parentDefinition.CapabilityType != nil) {
		self.CapabilityType = parentDefinition.CapabilityType
	}
	if (self.ValidSourceNodeTypes == nil) && (parentDefinition.ValidSourceNodeTypes != nil) {
		self.ValidSourceNodeTypes = parentDefinition.ValidSourceNodeTypes
	}
	if (self.ValidRelationshipTypes == nil) && (parentDefinition.ValidRelationshipTypes != nil) {
		self.ValidRelationshipTypes = parentDefinition.ValidRelationshipTypes
	}

	self.PropertyDefinitions.Inherit(parentDefinition.PropertyDefinitions)
	self.AttributeDefinitions.Inherit(parentDefinition.AttributeDefinitions)
}

// ([parsing.Renderable] interface)
func (self *CapabilityDefinition) Render() {
	// Avoid rendering more than once
	self.renderOnce.Do(self.render)
}

func (self *CapabilityDefinition) render() {
	logRender.Debugf("capability definition: %s", self.Name)

	if self.CapabilityTypeName == nil {
		// Avoid reporting more than once
		if !self.typeMissingProblemReported {
			self.Context.FieldChild("type", nil).ReportKeynameMissing()
			self.typeMissingProblemReported = true
		}
	} else {
		// If we have a capability type, inherit property definitions from it
		if self.CapabilityType != nil {
			// Inherit property definitions from capability type
			for name, propDef := range self.CapabilityType.PropertyDefinitions {
				if existingPropDef, exists := self.PropertyDefinitions[name]; exists {
					// If property exists but has no type, inherit from capability type
					if existingPropDef.DataTypeName == nil && propDef.DataTypeName != nil {
						existingPropDef.DataTypeName = propDef.DataTypeName
						existingPropDef.DataType = propDef.DataType
					}
				}
			}
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
		}
	}
}
