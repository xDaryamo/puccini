package v1_1

import (
	"github.com/tliron/puccini/ard"
	"github.com/tliron/puccini/tosca"
)

//
// CapabilityMapping
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
	if context.ValidateType("list") {
		list := context.Data.(ard.List)
		if len(list) != 2 {
			context.Report("must be 2")
			return self
		}

		nodeTemplateNameContext := context.ListChild(0, list[0])
		capabilityNameContext := context.ListChild(1, list[1])

		self.NodeTemplateName = nodeTemplateNameContext.ReadString()
		self.CapabilityName = capabilityNameContext.ReadString()
	}
	return self
}

func init() {
	Readers["CapabilityMapping"] = ReadCapabilityMapping
}
