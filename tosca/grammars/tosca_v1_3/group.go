package tosca_v1_3

import (
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
)

//
// Group
//
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.8.5
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.7.5
//

type Group struct {
	*Entity `name:"group"`
	Name    string `namespace:""`

	GroupTypeName           *string              `read:"type" require:"type"`
	Description             *string              `read:"description" inherit:"description,GroupType"`
	Properties              Values               `read:"properties,Value"`
	Interfaces              InterfaceAssignments `read:"interfaces,InterfaceAssignment"`
	MemberNodeTemplateNames *[]string            `read:"members"`

	GroupType           *GroupType      `lookup:"type,GroupTypeName" json:"-" yaml:"-"`
	MemberNodeTemplates []*NodeTemplate `lookup:"members,MemberNodeTemplateNames" json:"-" yaml:"-"`
}

func NewGroup(context *tosca.Context) *Group {
	return &Group{
		Entity:     NewEntity(context),
		Name:       context.Name,
		Properties: make(Values),
		Interfaces: make(InterfaceAssignments),
	}
}

// tosca.Reader signature
func ReadGroup(context *tosca.Context) interface{} {
	self := NewGroup(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))
	return self
}

// tosca.Renderable interface
func (self *Group) Render() {
	log.Infof("{render} group: %s", self.Name)

	if self.GroupType == nil {
		return
	}

	self.Properties.RenderProperties(self.GroupType.PropertyDefinitions, "property", self.Context.FieldChild("properties", nil))
	self.Interfaces.Render(self.GroupType.InterfaceDefinitions, self.Context.FieldChild("interfaces", nil))

	// Validate members
	for index, nodeTemplate := range self.MemberNodeTemplates {
		compatible := false
		for _, nodeType := range self.GroupType.MemberNodeTypes {
			if self.Context.Hierarchy.IsCompatible(nodeType, nodeTemplate.NodeType) {
				compatible = true
				break
			}
		}
		if !compatible {
			childContext := self.Context.FieldChild("members", nil).ListChild(index, nil)
			childContext.ReportIncompatible(nodeTemplate.Name, "group", "member")
		}
	}
}

func (self *Group) Normalize(s *normal.ServiceTemplate) *normal.Group {
	log.Infof("{normalize} group: %s", self.Name)

	g := s.NewGroup(self.Name)

	if self.Description != nil {
		g.Description = *self.Description
	}

	if types, ok := normal.GetTypes(self.Context.Hierarchy, self.GroupType); ok {
		g.Types = types
	}

	self.Properties.Normalize(g.Properties)
	self.Interfaces.NormalizeForGroup(self, g)

	for _, nodeTemplate := range self.MemberNodeTemplates {
		if n, ok := s.NodeTemplates[nodeTemplate.Name]; ok {
			g.Members = append(g.Members, n)
		}
	}

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
