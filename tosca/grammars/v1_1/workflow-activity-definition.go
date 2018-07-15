package v1_1

import (
	"github.com/tliron/puccini/ard"
	"github.com/tliron/puccini/tosca"
)

//
// WorkflowActivityDefinition
//

type WorkflowActivityDefinition struct {
	*Entity `name:"workflow activity definition"`

	DelegateWorkflowDefinitionName *string
	InlineWorkflowDefinitionName   *string
	SetNodeState                   *string
	CallOperationSpec              *string

	DelegateWorkflowDefinition *WorkflowDefinition `lookup:"delegate,DelegateWorkflowDefinitionName" json:"-" yaml:"-"`
	InlineWorkflowDefinition   *WorkflowDefinition `lookup:"inline,InlineWorkflowDefinitionName" json:"-" yaml:"-"`
}

func NewWorkflowActivityDefinition(context *tosca.Context) *WorkflowActivityDefinition {
	return &WorkflowActivityDefinition{Entity: NewEntity(context)}
}

// tosca.Reader signature
func ReadWorkflowActivityDefinition(context *tosca.Context) interface{} {
	self := NewWorkflowActivityDefinition(context)
	if context.ValidateType("map") {
		map_ := context.Data.(ard.Map)
		if len(map_) != 1 {
			context.ReportValueMalformed("workflow activity definition", "map length not 1")
			return self
		}

		for operator, value := range map_ {
			childContext := context.FieldChild(operator, value)

			switch operator {
			case "delegate":
				self.DelegateWorkflowDefinitionName = childContext.ReadString()
			case "inline":
				self.InlineWorkflowDefinitionName = childContext.ReadString()
			case "set_state":
				self.SetNodeState = childContext.ReadString()
			case "call_operation":
				self.CallOperationSpec = childContext.ReadString()
			default:
				context.ReportValueMalformed("workflow activity definition", "unsupported operator")
				return self
			}

			// We have only one key
			break
		}
	}
	return self
}
