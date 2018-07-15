package v1_1

import (
	"github.com/tliron/puccini/tosca"
)

//
// WorkflowPreconditionDefinition
//

type WorkflowPreconditionDefinition struct {
	*Entity `name:"workflow precondition definition"`

	TargetNodeTemplateOrGroupName *string            `read:"target" require:"target"`
	TargetNodeRequirementName     *string            `read:"target_relationship"`
	ConditionClauses              []*ConditionClause `read:"condition,[]ConditionClause"`

	TargetNodeTemplate *NodeTemplate `lookup:"target,TargetNodeTemplateOrGroupName" json:"-" yaml:"-"`
	TargetGroup        *Group        `lookup:"target,TargetNodeTemplateOrGroupName" json:"-" yaml:"-"`
}

func NewWorkflowPreconditionDefinition(context *tosca.Context) *WorkflowPreconditionDefinition {
	return &WorkflowPreconditionDefinition{Entity: NewEntity(context)}
}

// tosca.Reader signature
func ReadWorkflowPreconditionDefinition(context *tosca.Context) interface{} {
	self := NewWorkflowPreconditionDefinition(context)
	context.ValidateUnsupportedFields(context.ReadFields(self, Readers))
	return self
}
