package tosca_v1_3

import (
	"github.com/tliron/puccini/tosca"
)

//
// RequirementMapping
//
// Attaches to SubstitutionMappings
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
func ReadRequirementMapping(context *tosca.Context) interface{} {
	self := NewRequirementMapping(context)

	if strings := context.ReadStringListFixed(2); strings != nil {
		self.NodeTemplateName = &(*strings)[0]
		self.RequirementName = &(*strings)[1]
	}

	return self
}

// tosca.Renderable interface
func (self *RequirementMapping) Render() {
	log.Info("{render} requirement mapping")

	if (self.NodeTemplate == nil) || (self.RequirementName == nil) {
		return
	}

	name := *self.RequirementName
	found := false
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
