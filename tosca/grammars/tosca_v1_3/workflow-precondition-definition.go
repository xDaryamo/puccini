package tosca_v1_3

import (
	"github.com/tliron/puccini/tosca"
)

//
// WorkflowPreconditionDefinition
//
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.22
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.20
//

type WorkflowPreconditionDefinition struct {
	*Entity `name:"workflow precondition definition"`

	TargetNodeTemplateOrGroupName *string          `read:"target" require:"target"`
	TargetNodeRequirementName     *string          `read:"target_relationship"`
	ConditionClauses              ConditionClauses `read:"condition,[]ConditionClause"`

	TargetNodeTemplate *NodeTemplate `lookup:"target,TargetNodeTemplateOrGroupName" json:"-" yaml:"-"`
	TargetGroup        *Group        `lookup:"target,TargetNodeTemplateOrGroupName" json:"-" yaml:"-"`
}

func NewWorkflowPreconditionDefinition(context *tosca.Context) *WorkflowPreconditionDefinition {
	return &WorkflowPreconditionDefinition{Entity: NewEntity(context)}
}

// tosca.Reader signature
func ReadWorkflowPreconditionDefinition(context *tosca.Context) interface{} {
	self := NewWorkflowPreconditionDefinition(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))
	return self
}

//
// WorkflowPreconditionDefinitions
//

type WorkflowPreconditionDefinitions []*WorkflowPreconditionDefinition
