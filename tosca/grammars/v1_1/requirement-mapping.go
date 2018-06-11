package v1_1

import (
	"github.com/tliron/puccini/ard"
	"github.com/tliron/puccini/tosca"
)

//
// RequirementMapping
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
	if context.ValidateType("list") {
		list := context.Data.(ard.List)
		if len(list) != 2 {
			context.Report("must be 2")
			return self
		}

		nodeTemplateNameContext := context.ListChild(0, list[0])
		requirementNameContext := context.ListChild(1, list[1])

		self.NodeTemplateName = nodeTemplateNameContext.ReadString()
		self.RequirementName = requirementNameContext.ReadString()
	}
	return self
}

func init() {
	Readers["RequirementMapping"] = ReadRequirementMapping
}
