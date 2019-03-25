package cloudify_v1_3

import (
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
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
	TargetGroups Groups      `lookup:"targets,TargetGroupNames" json:"-" yaml:"-"`
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

func (self *Policy) Normalize(s *normal.ServiceTemplate) *normal.Policy {
	log.Infof("{normalize} policy: %s", self.Name)

	p := s.NewPolicy(self.Name)

	if types, ok := normal.GetTypes(self.Context.Hierarchy, self.PolicyType); ok {
		p.Types = types
	}

	self.Properties.Normalize(p.Properties, "")

	// TODO: groups

	return p
}

//
// Policies
//

type Policies []*Policy
