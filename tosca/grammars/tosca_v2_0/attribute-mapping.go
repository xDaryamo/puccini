package tosca_v2_0

import (
	"github.com/tliron/puccini/tosca"
)

//
// AttributeMapping
//
// Attaches to SubstitutionMappings
//
// [TOSCA-v2.0] @ ?
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.6.15
//

type AttributeMapping struct {
	*Entity `name:"attribute mapping"`

	NodeTemplateName *string `require:"0"`
	AttributeName    *string `require:"1"`

	NodeTemplate *NodeTemplate `lookup:"0,NodeTemplateName" json:"-" yaml:"-"`
}

func NewAttributeMapping(context *tosca.Context) *AttributeMapping {
	return &AttributeMapping{Entity: NewEntity(context)}
}

// tosca.Reader signature
func ReadAttributeMapping(context *tosca.Context) tosca.EntityPtr {
	self := NewAttributeMapping(context)

	if strings := context.ReadStringListFixed(2); strings != nil {
		self.NodeTemplateName = &(*strings)[0]
		self.AttributeName = &(*strings)[1]
	}

	return self
}

// parser.Renderable interface
func (self *AttributeMapping) Render() {
	logRender.Debug("attribute mapping")

	if (self.NodeTemplate == nil) || (self.AttributeName == nil) {
		return
	}

	name := *self.AttributeName
	self.NodeTemplate.Render()
	if _, ok := self.NodeTemplate.Attributes[name]; !ok {
		self.Context.ListChild(1, name).ReportReferenceNotFound("attribute", self.NodeTemplate)
	}
}

//
// AttributeMappings
//

type AttributeMappings []*AttributeMapping
