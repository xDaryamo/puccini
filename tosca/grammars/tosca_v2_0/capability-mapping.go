package tosca_v2_0

import (
	"github.com/tliron/puccini/tosca/parsing"
)

//
// CapabilityMapping
//
// Attaches to SubstitutionMappings
//
// [TOSCA-v2.0] @ 15.4
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.8.10
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.8.9
// [TOSCA-Simple-Profile-YAML-v1.1] @ 2.10, 2.11
// [TOSCA-Simple-Profile-YAML-v1.0] @ 2.10, 2.11
//

type CapabilityMapping struct {
	*Entity `name:"capability mapping"`
	Name    string

	NodeTemplateName             *string
	NodeTemplateCapabilityName   *string

	NodeTemplate *NodeTemplate         `traverse:"ignore" json:"-" yaml:"-"`
	Capability   *CapabilityAssignment `traverse:"ignore" json:"-" yaml:"-"`
}

func NewCapabilityMapping(context *parsing.Context) *CapabilityMapping {
	return &CapabilityMapping{
		Entity: NewEntity(context),
		Name:   context.Name,
	}
}

// ([parsing.Reader] signature)
func ReadCapabilityMapping(context *parsing.Context) parsing.EntityPtr {
	self := NewCapabilityMapping(context)

	if strings := context.ReadStringListFixed(2); strings != nil {
		self.NodeTemplateName = &(*strings)[0]
		self.NodeTemplateCapabilityName = &(*strings)[1]
	}

	return self
}

// ([parsing.Mappable] interface)
func (self *CapabilityMapping) GetKey() string {
	return self.Name
}

func (self *CapabilityMapping) GetCapabilityDefinition() (*CapabilityDefinition, bool) {
	if (self.Capability != nil) && (self.NodeTemplate != nil) {
		return self.Capability.GetDefinition(self.NodeTemplate)
	} else {
		return nil, false
	}
}

// ([parsing.Renderable] interface)
func (self *CapabilityMapping) Render() {
	// Avoid rendering more than once
	self.renderOnce.Do(self.render)
}

func (self *CapabilityMapping) render() {
	logRender.Debug("capability mapping")

	if (self.NodeTemplateName == nil) || (self.NodeTemplateCapabilityName == nil) {
		self.Context.ReportValueMalformed("capability mapping", "must specify both node template name and node template capability name")
		return
	}

	nodeTemplateName := *self.NodeTemplateName
	nodeTemplateCapabilityName := *self.NodeTemplateCapabilityName

	// Validate node template name is not empty
	if nodeTemplateName == "" {
		self.Context.ListChild(0, nodeTemplateName).ReportValueMalformed("node template name", "cannot be empty")
		return
	}

	// Validate node template capability name is not empty
	if nodeTemplateCapabilityName == "" {
		self.Context.ListChild(1, nodeTemplateCapabilityName).ReportValueMalformed("node template capability name", "cannot be empty")
		return
	}

	// Look up the node template within the substituting service template
	if nodeTemplate, ok := self.Context.Namespace.LookupForType(nodeTemplateName, nodeTemplatePtrType); ok {
		self.NodeTemplate = nodeTemplate.(*NodeTemplate)

		// Ensure the node template is rendered
		self.NodeTemplate.Render()

		// Look up the capability definition within the node template
		if capability, ok := self.NodeTemplate.Capabilities[nodeTemplateCapabilityName]; ok {
			self.Capability = capability
		} else {
			self.Context.ListChild(1, nodeTemplateCapabilityName).ReportReferenceNotFound("capability", self.NodeTemplate)
		}
	} else {
		self.Context.ListChild(0, nodeTemplateName).ReportUnknown("node template")
	}
}

//
// CapabilityMappings
//

type CapabilityMappings map[string]*CapabilityMapping
