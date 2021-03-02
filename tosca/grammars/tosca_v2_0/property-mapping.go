package tosca_v2_0

import (
	"github.com/tliron/kutil/ard"
	"github.com/tliron/puccini/tosca"
)

//
// PropertyMapping
//
// Attaches to SubstitutionMappings
//
// [TOSCA-v2.0] @ ?
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.8.8
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
func ReadPropertyMapping(context *tosca.Context) tosca.EntityPtr {
	self := NewPropertyMapping(context)
	if context.ValidateType(ard.TypeList) {
		strings := context.ReadStringListFixed(2)
		if strings != nil {
			self.NodeTemplateName = &(*strings)[0]
			self.PropertyName = &(*strings)[1]
		}
	}
	return self
}

// parser.Renderable interface
func (self *PropertyMapping) Render() {
	logRender.Debug("property mapping")

	if (self.NodeTemplate == nil) || (self.PropertyName == nil) {
		return
	}

	name := *self.PropertyName
	self.NodeTemplate.Render()
	if _, ok := self.NodeTemplate.Properties[name]; !ok {
		self.Context.ListChild(1, name).ReportReferenceNotFound("property", self.NodeTemplate)
	}
}

//
// PropertyMappings
//

type PropertyMappings []*PropertyMapping
