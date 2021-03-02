package tosca_v2_0

import (
	"github.com/tliron/puccini/tosca"
)

//
// CapabilityMapping
//
// Attaches to SubstitutionMappings
//
// [TOSCA-v2.0] @ ?
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.8.10
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.8.9
// [TOSCA-Simple-Profile-YAML-v1.1] @ 2.10, 2.11
// [TOSCA-Simple-Profile-YAML-v1.0] @ 2.10, 2.11
//

type CapabilityMapping struct {
	*Entity `name:"capability mapping"`

	NodeTemplateName *string `require:"0"`
	CapabilityName   *string `require:"1"`

	NodeTemplate *NodeTemplate `lookup:"0,NodeTemplateName" json:"-" yaml:"-"`
}

func NewCapabilityMapping(context *tosca.Context) *CapabilityMapping {
	return &CapabilityMapping{Entity: NewEntity(context)}
}

// tosca.Reader signature
func ReadCapabilityMapping(context *tosca.Context) tosca.EntityPtr {
	self := NewCapabilityMapping(context)

	if strings := context.ReadStringListFixed(2); strings != nil {
		self.NodeTemplateName = &(*strings)[0]
		self.CapabilityName = &(*strings)[1]
	}

	return self
}

// parser.Renderable interface
func (self *CapabilityMapping) Render() {
	logRender.Debug("capability mapping")

	if (self.NodeTemplate == nil) || (self.CapabilityName == nil) {
		return
	}

	name := *self.CapabilityName
	self.NodeTemplate.Render()
	if _, ok := self.NodeTemplate.Capabilities[name]; !ok {
		log.Debugf("%s", self.NodeTemplate.Capabilities)
		self.Context.ListChild(1, name).ReportReferenceNotFound("capability", self.NodeTemplate)
	}
}

//
// CapabilityMappings
//

type CapabilityMappings []*CapabilityMapping
