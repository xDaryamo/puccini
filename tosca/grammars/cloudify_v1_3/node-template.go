package cloudify_v1_3

import (
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
)

//
// NodeTemplate
//
// [https://docs.cloudify.co/4.5.5/developer/blueprints/spec-node-templates/]
//

type NodeTemplate struct {
	*Entity `name:"node template"`
	Name    string `namespace:""`

	NodeTypeName  *string                  `read:"type" require:"type"`
	Properties    Values                   `read:"properties,Value"`
	Instances     *NodeTemplateInstances   `read:"instances,NodeTemplateInstances"` // deprecated in Cloudify DSL 1.3
	Interfaces    InterfaceAssignments     `read:"interfaces,InterfaceAssignment"`
	Relationships RelationshipAssignments  `read:"relationships,[]RelationshipAssignment"`
	Capabilities  NodeTemplateCapabilities `read:"capabilities,NodeTemplateCapability"`

	NodeType *NodeType `lookup:"type,NodeTypeName" json:"-" yaml:"-"`
}

func NewNodeTemplate(context *tosca.Context) *NodeTemplate {
	return &NodeTemplate{
		Entity:       NewEntity(context),
		Name:         context.Name,
		Properties:   make(Values),
		Interfaces:   make(InterfaceAssignments),
		Capabilities: make(NodeTemplateCapabilities),
	}
}

// tosca.Reader signature
func ReadNodeTemplate(context *tosca.Context) interface{} {
	self := NewNodeTemplate(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))
	self.Capabilities.Validate(context, self.Instances)
	return self
}

// tosca.Renderable interface
func (self *NodeTemplate) Render() {
	log.Infof("{render} node template: %s", self.Name)

	if self.NodeType == nil {
		return
	}

	self.Properties.RenderProperties(self.NodeType.PropertyDefinitions, "property", self.Context.FieldChild("properties", nil))
	self.Interfaces.Render(self.NodeType.InterfaceDefinitions, self.Context.FieldChild("interfaces", nil))
}

var capabilityTypeName = "cloudify.Node"
var capabilityTypes = normal.NewTypes(capabilityTypeName)

func (self *NodeTemplate) Normalize(s *normal.ServiceTemplate) *normal.NodeTemplate {
	log.Infof("{normalize} node template: %s", self.Name)

	n := s.NewNodeTemplate(self.Name)

	if types, ok := normal.GetTypes(self.Context.Hierarchy, self.NodeType); ok {
		n.Types = types
	}

	self.Properties.Normalize(n.Properties, "")
	self.Interfaces.NormalizeForNodeTemplate(self, n)

	c := n.NewCapability("node")
	c.Types = capabilityTypes
	for _, capability := range self.Capabilities {
		capability.Properties.Normalize(c.Properties, capability.Context.Name+".")
	}

	return n
}

func (self *NodeTemplate) NormalizeRelationships(s *normal.ServiceTemplate) {
	log.Infof("{normalize} node template relationships: %s", self.Name)

	n := s.NodeTemplates[self.Name]
	for _, relationship := range self.Relationships {
		relationship.Normalize(self, s, n)
	}
}

//
// NodeTemplates
//

type NodeTemplates []*NodeTemplate

func (self NodeTemplates) Normalize(s *normal.ServiceTemplate) {
	for _, nodeTemplate := range self {
		s.NodeTemplates[nodeTemplate.Name] = nodeTemplate.Normalize(s)
	}

	// Relationships must be normalized after node templates
	// (because they may reference other node templates)
	for _, nodeTemplate := range self {
		nodeTemplate.NormalizeRelationships(s)
	}
}
