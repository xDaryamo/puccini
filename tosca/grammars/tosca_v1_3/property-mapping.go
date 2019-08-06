package tosca_v1_3

import (
	"github.com/tliron/puccini/tosca"
)

//
// PropertyMapping
//
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.8.8
//

type PropertyMapping struct {
	*Entity `name:"property mapping"`

	NodeTemplateName *string `require:"0"`
	PropertyName     *string `require:"1"`

	NodeTemplate *NodeTemplate `lookup:"0,NodeTemplateName" json:"-" yaml:"-"`
}

func NewPropertyMapping(context *tosca.Context) *PropertyMapping {
	return &PropertyMapping{Entity: NewEntity(context)}
}

// tosca.Reader signature
func ReadPropertyMapping(context *tosca.Context) interface{} {
	self := NewPropertyMapping(context)
	if context.ValidateType("list") {
		strings := context.ReadStringListFixed(2)
		if strings != nil {
			self.NodeTemplateName = &(*strings)[0]
			self.PropertyName = &(*strings)[1]
		}
	}
	return self
}

// tosca.Renderable interface
func (self *PropertyMapping) Render() {
	log.Info("{render} property mapping")

	if (self.NodeTemplate == nil) || (self.NodeTemplate.NodeType == nil) || (self.PropertyName == nil) {
		return
	}

	name := *self.PropertyName
	if _, ok := self.NodeTemplate.NodeType.PropertyDefinitions[name]; !ok {
		self.Context.ListChild(1, name).ReportReferenceNotFound("property", self.NodeTemplate)
	}
}

//
// PropertyMappings
//

type PropertyMappings []*PropertyMapping
