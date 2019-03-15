package cloudify_v1_3

import (
	"github.com/tliron/puccini/tosca"
)

//
// GroupPolicyTriggerTrigger
//
// [https://docs.cloudify.co/4.5.5/developer/blueprints/spec-groups/]
//

type GroupPolicyTrigger struct {
	*Entity `name:"group policy trigger"`

	PolicyTriggerTypeName *string `read:"type" require:"type"`
	Parameters            Values  `read:"parameters,Value"`

	PolicyTriggerType *PolicyTriggerType `lookup:"type,PolicyTriggerTypeName" json:"-" yaml:"-"`
}

func NewGroupPolicyTrigger(context *tosca.Context) *GroupPolicyTrigger {
	return &GroupPolicyTrigger{
		Entity:     NewEntity(context),
		Parameters: make(Values),
	}
}

// tosca.Reader signature
func ReadGroupPolicyTrigger(context *tosca.Context) interface{} {
	self := NewGroupPolicyTrigger(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))
	return self
}

//
// GroupPolicyTriggers
//

type GroupPolicyTriggers map[string]*GroupPolicyTrigger
