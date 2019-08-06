package tosca_v1_3

import (
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
)

//
// NodeTemplate
//
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.8.3
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.7.3
//

type NodeTemplate struct {
	*Entity `name:"node template"`
	Name    string `namespace:""`

	Directives                   *[]string              `read:"directives"`
	CopyNodeTemplateName         *string                `read:"copy"`
	NodeTypeName                 *string                `read:"type" require:"type"`
	Description                  *string                `read:"description" inherit:"description,NodeType"`
	Properties                   Values                 `read:"properties,Value"`
	Attributes                   Values                 `read:"attributes,AttributeValue"`
	Capabilities                 CapabilityAssignments  `read:"capabilities,CapabilityAssignment"`
	Requirements                 RequirementAssignments `read:"requirements,{}RequirementAssignment"`
	RequirementTargetsNodeFilter *NodeFilter            `read:"node_filter,NodeFilter"`
	Interfaces                   InterfaceAssignments   `read:"interfaces,InterfaceAssignment"`
	Artifacts                    Artifacts              `read:"artifacts,Artifact"`

	CopyNodeTemplate *NodeTemplate `lookup:"copy,CopyNodeTemplateName" json:"-" yaml:"-"`
	NodeType         *NodeType     `lookup:"type,NodeTypeName" json:"-" yaml:"-"`
}

func NewNodeTemplate(context *tosca.Context) *NodeTemplate {
	return &NodeTemplate{
		Entity:       NewEntity(context),
		Name:         context.Name,
		Properties:   make(Values),
		Attributes:   make(Values),
		Capabilities: make(CapabilityAssignments),
		Interfaces:   make(InterfaceAssignments),
		Artifacts:    make(Artifacts),
	}
}

// tosca.Reader signature
func ReadNodeTemplate(context *tosca.Context) interface{} {
	self := NewNodeTemplate(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))
	return self
}

func (self *NodeTemplate) GetCapabilitiesOfType(capabilityType *CapabilityType) []*CapabilityAssignment {
	var capabilities []*CapabilityAssignment
	for _, capability := range self.Capabilities {
		definition, ok := capability.GetDefinition(self)
		if ok && (definition.CapabilityType != nil) && self.Context.Hierarchy.IsCompatible(capabilityType, definition.CapabilityType) {
			capabilities = append(capabilities, capability)
		}
	}
	return capabilities
}

func (self *NodeTemplate) GetCapabilityByName(capabilityName string, capabilityType *CapabilityType) (*CapabilityAssignment, bool) {
	if capabilityType == nil {
		return nil, false
	}

	capability, ok := self.Capabilities[capabilityName]
	if !ok {
		return nil, false
	}

	// Make sure it matches the type
	definition, ok := capability.GetDefinition(self)
	if ok && self.Context.Hierarchy.IsCompatible(capabilityType, definition.CapabilityType) {
		return capability, true
	}

	return nil, false
}

// tosca.Renderable interface
func (self *NodeTemplate) Render() {
	log.Infof("{render} node template: %s", self.Name)

	// TODO: copy

	if self.NodeType == nil {
		return
	}

	self.Properties.RenderProperties(self.NodeType.PropertyDefinitions, "property", self.Context.FieldChild("properties", nil))
	self.Attributes.RenderAttributes(self.NodeType.AttributeDefinitions, self.Context.FieldChild("attributes", nil))
	self.Capabilities.Render(self.NodeType.CapabilityDefinitions, self.Context.FieldChild("capabilities", nil))
	self.Requirements.Render(self.NodeType.RequirementDefinitions, self.Context.FieldChild("requirements", nil))
	self.Interfaces.Render(self.NodeType.InterfaceDefinitions, self.Context.FieldChild("interfaces", nil))
	self.Artifacts.Render(self.NodeType.ArtifactDefinitions, self.Context.FieldChild("artifacts", nil))
}

func (self *NodeTemplate) Normalize(s *normal.ServiceTemplate) *normal.NodeTemplate {
	log.Infof("{normalize} node template: %s", self.Name)

	n := s.NewNodeTemplate(self.Name)

	if self.Description != nil {
		n.Description = *self.Description
	}

	if types, ok := normal.GetTypes(self.Context.Hierarchy, self.NodeType); ok {
		n.Types = types
	}

	if self.Directives != nil {
		n.Directives = *self.Directives
	}

	self.Properties.Normalize(n.Properties)
	self.Attributes.Normalize(n.Attributes)
	self.Capabilities.Normalize(self, n)
	self.Interfaces.NormalizeForNodeTemplate(self, n)
	self.Artifacts.Normalize(n)

	return n
}

//
// NodeTemplates
//

type NodeTemplates []*NodeTemplate

func (self NodeTemplates) Normalize(s *normal.ServiceTemplate) {
	for _, nodeTemplate := range self {
		s.NodeTemplates[nodeTemplate.Name] = nodeTemplate.Normalize(s)
	}

	// Requirements must be normalized after node templates
	// (because they may reference other node templates)
	for _, nodeTemplate := range self {
		n := s.NodeTemplates[nodeTemplate.Name]
		nodeTemplate.Requirements.Normalize(nodeTemplate, s, n)
	}
}
