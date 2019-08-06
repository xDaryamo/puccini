package tosca_v1_3

import (
	"github.com/tliron/puccini/tosca"
)

//
// CapabilityMapping
//
// Attaches to SubstitutionMappings
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
func ReadCapabilityMapping(context *tosca.Context) interface{} {
	self := NewCapabilityMapping(context)

	if strings := context.ReadStringListFixed(2); strings != nil {
		self.NodeTemplateName = &(*strings)[0]
		self.CapabilityName = &(*strings)[1]
	}

	return self
}

// tosca.Renderable interface
func (self *CapabilityMapping) Render() {
	log.Info("{render} capability mapping")

	if (self.NodeTemplate == nil) || (self.CapabilityName == nil) {
		return
	}

	name := *self.CapabilityName
	if _, ok := self.NodeTemplate.Capabilities[name]; !ok {
		self.Context.ListChild(1, name).ReportReferenceNotFound("capability", self.NodeTemplate)
	}
}

//
// CapabilityMappings
//

type CapabilityMappings []*CapabilityMapping
