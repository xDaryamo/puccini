package cloudify_v1_3

import (
	"github.com/tliron/puccini/tosca"
)

//
// RelationshipType
//
// [https://docs.cloudify.co/4.5.5/developer/blueprints/spec-relationships/]
//

type RelationshipType struct {
	*Type `name:"relationship type"`

	SourceInterfaces InterfaceDefinitions `read:"source_interfaces,InterfaceDefinition" inherit:"source_interfaces,Parent"`
	TargetInterfaces InterfaceDefinitions `read:"target_interfaces,InterfaceDefinition" inherit:"target_interfaces,Parent"`
	Properties       Values               `read:"properties,Value"`

	Parent *RelationshipType `lookup:"derived_from,ParentName" json:"-" yaml:"-"`
}

func NewRelationshipType(context *tosca.Context) *RelationshipType {
	return &RelationshipType{
		Type:             NewType(context),
		SourceInterfaces: make(InterfaceDefinitions),
		TargetInterfaces: make(InterfaceDefinitions),
		Properties:       make(Values),
	}
}

// tosca.Reader signature
func ReadRelationshipType(context *tosca.Context) interface{} {
	self := NewRelationshipType(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))
	ValidateRelationshipProperties(context, self.Properties)
	return self
}
