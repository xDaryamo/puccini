package cloudify_v1_3

import (
	"github.com/tliron/puccini/tosca"
)

//
// Group
//
// [https://docs.cloudify.co/4.5.5/developer/blueprints/spec-groups/]
//

type Group struct {
	*Entity `name:"group"`
	Name    string `namespace:""`

	MemberNodeTemplateNames *[]string     `read:"members" require:"members"`
	Policies                GroupPolicies `read:"policies,GroupPolicy"`

	MemberNodeTemplates []*NodeTemplate `members:"type,MemberNodeTemplateNames" json:"-" yaml:"-"`
}

func NewGroup(context *tosca.Context) *Group {
	return &Group{
		Entity:   NewEntity(context),
		Name:     context.Name,
		Policies: make(GroupPolicies),
	}
}

// tosca.Reader signature
func ReadGroup(context *tosca.Context) interface{} {
	self := NewGroup(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))
	return self
}
