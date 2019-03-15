package cloudify_v1_3

import (
	"github.com/tliron/puccini/tosca"
)

//
// Policy
//
// [https://docs.cloudify.co/4.5.5/developer/blueprints/spec-policies/]
//

type Policy struct {
	*Entity `name:"policy"`
	Name    string `namespace:""`

	PolicyTypeName   *string   `read:"type" require:"type"`
	Properties       Values    `read:"properties,Value"`
	TargetGroupNames *[]string `read:"targets" require:"targets"`

	PolicyType   *PolicyType `lookup:"type,PolicyTypeName" json:"-" yaml:"-"`
	TargetGroups []*Group    `lookup:"targets,TargetGroupNames" json:"-" yaml:"-"`
}

func NewPolicy(context *tosca.Context) *Policy {
	return &Policy{
		Entity:     NewEntity(context),
		Name:       context.Name,
		Properties: make(Values),
	}
}

// tosca.Reader signature
func ReadPolicy(context *tosca.Context) interface{} {
	self := NewPolicy(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))
	return self
}
