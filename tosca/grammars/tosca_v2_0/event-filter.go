package tosca_v2_0

import (
	"github.com/tliron/puccini/tosca"
)

//
// EventFilter
//
// [TOSCA-v2.0] @ ?
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.6.21
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.17
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.15
//

type EventFilter struct {
	*Entity `name:"event filter" json:"-" yaml:"-"`

	NodeTemplateNameOrTypeName *string `read:"node"`
	RequirementName            *string `read:"requirement"`
	CapabilityName             *string `read:"capability"`

	NodeTemplate *NodeTemplate `lookup:"node,NodeTemplateNameOrTypeName" json:"-" yaml:"-"`
	NodeType     *NodeType     `lookup:"node,NodeTemplateNameOrTypeName" json:"-" yaml:"-"`
}

func NewEventFilter(context *tosca.Context) *EventFilter {
	return &EventFilter{Entity: NewEntity(context)}
}

// tosca.Reader signature
func ReadEventFilter(context *tosca.Context) tosca.EntityPtr {
	self := NewEventFilter(context)
	context.ReadFields(self)
	return self
}
