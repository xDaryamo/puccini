package cloudify_v1_3

import (
	"github.com/tliron/puccini/tosca"
)

//
// GroupPolicy
//
// [https://docs.cloudify.co/5.0.5/developer/blueprints/spec-groups/]
//

type GroupPolicy struct {
	*Entity `name:"group policy"`

	PolicyTypeName *string             `read:"type" require:""`
	Properties     Values              `read:"properties,Value"`
	Triggers       GroupPolicyTriggers `read:"triggers,GroupPolicyTrigger"`

	PolicyType *PolicyType `lookup:"type,PolicyTypeName" json:"-" yaml:"-"`
}

func NewGroupPolicy(context *tosca.Context) *GroupPolicy {
	return &GroupPolicy{
		Entity:     NewEntity(context),
		Properties: make(Values),
		Triggers:   make(GroupPolicyTriggers),
	}
}

// tosca.Reader signature
func ReadGroupPolicy(context *tosca.Context) tosca.EntityPtr {
	self := NewGroupPolicy(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))
	return self
}

//
// GroupPolicies
//

type GroupPolicies map[string]*GroupPolicy
