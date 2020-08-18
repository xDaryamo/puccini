package cloudify_v1_3

import (
	"github.com/tliron/puccini/tosca"
)

//
// PolicyType
//
// [https://docs.cloudify.co/5.0.5/developer/blueprints/spec-policy-types/]
//

type PolicyType struct {
	*Entity `name:"policy type"`
	Name    string `namespace:""`

	Source              *string             `read:"source" require:""`
	PropertyDefinitions PropertyDefinitions `read:"properties,PropertyDefinition"`
}

func NewPolicyType(context *tosca.Context) *PolicyType {
	return &PolicyType{
		Entity:              NewEntity(context),
		Name:                context.Name,
		PropertyDefinitions: make(PropertyDefinitions),
	}
}

// tosca.Reader signature
func ReadPolicyType(context *tosca.Context) tosca.EntityPtr {
	self := NewPolicyType(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))
	return self
}

var policyTypeRoot *PolicyType

// tosca.Hierarchical interface
func (self *PolicyType) GetParent() tosca.EntityPtr {
	return policyTypeRoot
}

//
// PolicyTypes
//

type PolicyTypes []*PolicyType
