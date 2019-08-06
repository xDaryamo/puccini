package tosca_v1_3

import (
	"github.com/tliron/puccini/tosca"
)

//
// CapabilityType
//
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.7.7
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.6.7
//

type CapabilityType struct {
	*Type `name:"capability type"`

	PropertyDefinitions      PropertyDefinitions  `read:"properties,PropertyDefinition" inherit:"properties,Parent"`
	AttributeDefinitions     AttributeDefinitions `read:"attributes,AttributeDefinition" inherit:"attributes,Parent"`
	ValidSourceNodeTypeNames *[]string            `read:"valid_source_types" inherit:"valid_source_types,Parent"`

	Parent               *CapabilityType `lookup:"derived_from,ParentName" json:"-" yaml:"-"`
	ValidSourceNodeTypes NodeTypes       `lookup:"valid_source_types,ValidSourceNodeTypeNames" inherit:"valid_source_types,Parent" json:"-" yaml:"-"`
}

func NewCapabilityType(context *tosca.Context) *CapabilityType {
	return &CapabilityType{
		Type:                 NewType(context),
		PropertyDefinitions:  make(PropertyDefinitions),
		AttributeDefinitions: make(AttributeDefinitions),
	}
}

// tosca.Reader signature
func ReadCapabilityType(context *tosca.Context) interface{} {
	self := NewCapabilityType(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))
	return self
}

// tosca.Hierarchical interface
func (self *CapabilityType) GetParent() interface{} {
	return self.Parent
}

// tosca.Inherits interface
func (self *CapabilityType) Inherit() {
	log.Infof("{inherit} capability type: %s", self.Name)

	if self.Parent == nil {
		return
	}

	self.PropertyDefinitions.Inherit(self.Parent.PropertyDefinitions)
	self.AttributeDefinitions.Inherit(self.Parent.AttributeDefinitions)
}

//
// CapabilityTypes
//

type CapabilityTypes []*CapabilityType
