package cloudify_v1_3

import (
	"github.com/tliron/puccini/tosca"
)

//
// PolicyType
//
// [https://docs.cloudify.co/4.5.5/developer/blueprints/spec-policy-types/]
//

type PolicyType struct {
	*Entity `name:"policy type"`
	Name    string `namespace:""`

	Source              *string             `read:"source" require:"source"`
	PropertyDefinitions PropertyDefinitions `read:"properties,PropertyDefinition"`
}

func NewPolicyType(context *tosca.Context) *PolicyType {
	return &PolicyType{
		Entity:              NewEntity(context),
		Name:                context.Name,
		PropertyDefinitions: make(PropertyDefinitions),
	}
}

var policyTypeRoot *PolicyType

// tosca.Hierarchical interface
func (self *PolicyType) GetParent() interface{} {
	return policyTypeRoot
}

// tosca.Reader signature
func ReadPolicyType(context *tosca.Context) interface{} {
	self := NewPolicyType(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))
	return self
}

//
// PolicyTypes
//

type PolicyTypes []*PolicyType
