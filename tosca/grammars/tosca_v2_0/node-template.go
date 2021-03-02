package tosca_v2_0

import (
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
)

//
// NodeTemplate
//
// [TOSCA-v2.0] @ ?
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.8.3
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.8.3
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.7.3
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.7.3
//

type NodeTemplate struct {
	*Entity `name:"node template"`
	Name    string `namespace:""`

	Directives                   *[]string              `read:"directives"`
	CopyNodeTemplateName         *string                `read:"copy"`
	NodeTypeName                 *string                `read:"type" require:""`
	Metadata                     Metadata               `read:"metadata,Metadata"` // introduced in TOSCA 1.1
	Description                  *string                `read:"description"`
	Properties                   Values                 `read:"properties,Value"`
	Attributes                   Values                 `read:"attributes,AttributeValue"`
	Capabilities                 CapabilityAssignments  `read:"capabilities,CapabilityAssignment"`
	Requirements                 RequirementAssignments `read:"requirements,{}RequirementAssignment"`
	RequirementTargetsNodeFilter *NodeFilter            `read:"node_filter,NodeFilter"`
	Interfaces                   InterfaceAssignments   `read:"interfaces,InterfaceAssignment"`
	Artifacts                    Artifacts              `read:"artifacts,Artifact"`

	CopyNodeTemplate *NodeTemplate `lookup:"copy,CopyNodeTemplateName" json:"-" yaml:"-"`
	NodeType         *NodeType     `lookup:"type,NodeTypeName" json:"-" yaml:"-"`

	rendered bool
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
func ReadNodeTemplate(context *tosca.Context) tosca.EntityPtr {
	self := NewNodeTemplate(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))
	return self
}

// tosca.PreReadable interface
func (self *NodeTemplate) PreRead() {
	CopyTemplate(self.Context)
}

// parser.Renderable interface
func (self *NodeTemplate) Render() {
	logRender.Debugf("node template: %s", self.Name)

	if self.rendered {
		// Avoid rendering more than once (can happen if we were called from PropertyMapping etc. Render)
		return
	}
	self.rendered = true

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

func (self *NodeTemplate) Normalize(normalServiceTemplate *normal.ServiceTemplate) *normal.NodeTemplate {
	logNormalize.Debugf("node template: %s", self.Name)

	normalNodeTemplate := normalServiceTemplate.NewNodeTemplate(self.Name)

	normalNodeTemplate.Metadata = self.Metadata

	if self.Description != nil {
		normalNodeTemplate.Description = *self.Description
	}

	if types, ok := normal.GetTypes(self.Context.Hierarchy, self.NodeType); ok {
		normalNodeTemplate.Types = types
	}

	if self.Directives != nil {
		normalNodeTemplate.Directives = *self.Directives
	}

	self.Properties.Normalize(normalNodeTemplate.Properties)
	self.Attributes.Normalize(normalNodeTemplate.Attributes)
	self.Capabilities.Normalize(self, normalNodeTemplate)
	self.Interfaces.NormalizeForNodeTemplate(self, normalNodeTemplate)
	self.Artifacts.Normalize(normalNodeTemplate)

	return normalNodeTemplate
}

//
// NodeTemplates
//

type NodeTemplates []*NodeTemplate

func (self NodeTemplates) Normalize(normalServiceTemplate *normal.ServiceTemplate) {
	for _, nodeTemplate := range self {
		normalServiceTemplate.NodeTemplates[nodeTemplate.Name] = nodeTemplate.Normalize(normalServiceTemplate)
	}

	// Requirements must be normalized after node templates
	// (because they may reference other node templates)
	for _, nodeTemplate := range self {
		if normalNodeTemplate, ok := normalServiceTemplate.NodeTemplates[nodeTemplate.Name]; ok {
			nodeTemplate.Requirements.Normalize(nodeTemplate, normalNodeTemplate)
		}
	}
}
