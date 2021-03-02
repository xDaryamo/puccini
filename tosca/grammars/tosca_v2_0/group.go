package tosca_v2_0

import (
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
)

//
// Group
//
// [TOSCA-v2.0] @ ?
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.8.5
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.8.5
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.7.5
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.7.5
//

type Group struct {
	*Entity `name:"group"`
	Name    string `namespace:""`

	GroupTypeName           *string              `read:"type" require:""`
	Metadata                Metadata             `read:"metadata,Metadata"` // introduced in TOSCA 1.1
	Description             *string              `read:"description"`
	Properties              Values               `read:"properties,Value"`
	Interfaces              InterfaceAssignments // removed in TOSCA 1.3
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
func ReadGroup(context *tosca.Context) tosca.EntityPtr {
	self := NewGroup(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))
	return self
}

// parser.Renderable interface
func (self *Group) Render() {
	logRender.Debugf("group: %s", self.Name)

	if self.GroupType == nil {
		return
	}

	self.Properties.RenderProperties(self.GroupType.PropertyDefinitions, "property", self.Context.FieldChild("properties", nil))
	self.Interfaces.Render(self.GroupType.InterfaceDefinitions, self.Context.FieldChild("interfaces", nil))

	// Validate members
	if len(self.GroupType.MemberNodeTypes) > 0 {
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
}

func (self *Group) Normalize(normalServiceTemplate *normal.ServiceTemplate) *normal.Group {
	logNormalize.Debugf("group: %s", self.Name)

	normalGroup := normalServiceTemplate.NewGroup(self.Name)

	normalGroup.Metadata = self.Metadata

	if self.Description != nil {
		normalGroup.Description = *self.Description
	}

	if types, ok := normal.GetTypes(self.Context.Hierarchy, self.GroupType); ok {
		normalGroup.Types = types
	}

	self.Properties.Normalize(normalGroup.Properties)
	self.Interfaces.NormalizeForGroup(self, normalGroup)

	for _, nodeTemplate := range self.MemberNodeTemplates {
		if normalNodeTemplate, ok := normalServiceTemplate.NodeTemplates[nodeTemplate.Name]; ok {
			normalGroup.Members = append(normalGroup.Members, normalNodeTemplate)
		}
	}

	return normalGroup
}

//
// Groups
//

type Groups []*Group

func (self Groups) Normalize(normalServiceTemplate *normal.ServiceTemplate) {
	for _, group := range self {
		normalServiceTemplate.Groups[group.Name] = group.Normalize(normalServiceTemplate)
	}
}
