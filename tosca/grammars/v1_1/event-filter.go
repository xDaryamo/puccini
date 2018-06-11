package v1_1

import (
	"github.com/tliron/puccini/tosca"
)

//
// EventFilter
//

type EventFilter struct {
	*Entity `name:"event filter" json:"-" yaml:"-"`

	NodeTemplateNameOrNodeTypeName *string `read:"node"`
	RequirementName                *string `read:"requirement"`
	CapabilityName                 *string `read:"capability"`

	NodeTemplate *NodeTemplate `lookup:"node,NodeTemplateNameOrNodeTypeName" json:"-" yaml:"-"`
	NodeType     *NodeType     `lookup:"node,NodeTemplateNameOrNodeTypeName" json:"-" yaml:"-"`
}

func NewEventFilter(context *tosca.Context) *EventFilter {
	return &EventFilter{Entity: NewEntity(context)}
}

// tosca.Reader signature
func ReadEventFilter(context *tosca.Context) interface{} {
	self := NewEventFilter(context)
	context.ReadFields(self, Readers)
	return self
}

func init() {
	Readers["EventFilter"] = ReadEventFilter
}
