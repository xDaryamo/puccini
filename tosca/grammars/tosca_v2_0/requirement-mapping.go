package tosca_v2_0

import (
	"github.com/tliron/puccini/tosca"
)

//
// RequirementMapping
//
// Attaches to SubstitutionMappings
//
// [TOSCA-v2.0] @ ?
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.8.8
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.8.8
// [TOSCA-Simple-Profile-YAML-v1.1] @ 2.10, 2.11
// [TOSCA-Simple-Profile-YAML-v1.0] @ 2.10, 2.11
//

type RequirementMapping struct {
	*Entity `name:"requirement mapping"`

	NodeTemplateName *string `require:"0"`
	RequirementName  *string `require:"1"`

	NodeTemplate *NodeTemplate `lookup:"0,NodeTemplateName" json:"-" yaml:"-"`
}

func NewRequirementMapping(context *tosca.Context) *RequirementMapping {
	return &RequirementMapping{Entity: NewEntity(context)}
}

// tosca.Reader signature
func ReadRequirementMapping(context *tosca.Context) tosca.EntityPtr {
	self := NewRequirementMapping(context)

	if strings := context.ReadStringListFixed(2); strings != nil {
		self.NodeTemplateName = &(*strings)[0]
		self.RequirementName = &(*strings)[1]
	}

	return self
}

// parser.Renderable interface
func (self *RequirementMapping) Render() {
	logRender.Debug("requirement mapping")

	if (self.NodeTemplate == nil) || (self.RequirementName == nil) {
		return
	}

	name := *self.RequirementName
	found := false
	self.NodeTemplate.Render()
	for _, requirement := range self.NodeTemplate.Requirements {
		if requirement.Name == name {
			if found {
				self.Context.ListChild(1, name).ReportReferenceAmbiguous("requirement", self.NodeTemplate)
				break
			} else {
				found = true
			}
		}
	}

	if !found {
		self.Context.ListChild(1, name).ReportReferenceNotFound("requirement", self.NodeTemplate)
	}
}

//
// RequirementMappings
//

type RequirementMappings []*RequirementMapping
