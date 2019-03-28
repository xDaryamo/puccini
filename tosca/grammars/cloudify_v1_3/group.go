package cloudify_v1_3

import (
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
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

	MemberNodeTemplates NodeTemplates `lookup:"members,MemberNodeTemplateNames" json:"-" yaml:"-"`
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

var groupTypeName = "cloudify.Group"
var groupTypes = normal.NewTypes(groupTypeName)

func (self *Group) Normalize(s *normal.ServiceTemplate) *normal.Group {
	log.Infof("{normalize} group: %s", self.Name)

	g := s.NewGroup(self.Name)
	g.Types = groupTypes

	for _, nodeTemplate := range self.MemberNodeTemplates {
		if n, ok := s.NodeTemplates[nodeTemplate.Name]; ok {
			g.Members = append(g.Members, n)
		}
	}

	// TODO: normalize policies
	// TODO: normalize triggers in policies

	return g
}

//
// Groups
//

type Groups []*Group

func (self Groups) Normalize(s *normal.ServiceTemplate) {
	for _, group := range self {
		s.Groups[group.Name] = group.Normalize(s)
	}
}
