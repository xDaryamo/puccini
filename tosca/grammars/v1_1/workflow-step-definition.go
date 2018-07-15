package v1_1

import (
	"github.com/tliron/puccini/tosca"
)

//
// WorkflowStepDefinition
//

type WorkflowStepDefinition struct {
	*Entity `name:"workflow step definition"`
	Name    string

	TargetNodeTemplateOrGroupName *string                       `read:"target" require:"target"`
	TargetNodeRequirementName     *string                       `read:"target_relationship"`
	FilterConstraintClauses       ConstraintClauses             `read:"filter,[]ConstraintClause"`
	ActivityDefinitions           []*WorkflowActivityDefinition `read:"activities,[]WorkflowActivityDefinition" require:"activities"`
	OnSuccessStepNames            *[]string                     `read:"on_success"`
	OnFailureStepNames            *[]string                     `read:"on_failure"`

	TargetNodeTemplate *NodeTemplate `lookup:"target,TargetNodeTemplateOrGroupName" json:"-" yaml:"-"`
	TargetGroup        *Group        `lookup:"target,TargetNodeTemplateOrGroupName" json:"-" yaml:"-"`
}

func NewWorkflowStepDefinition(context *tosca.Context) *WorkflowStepDefinition {
	return &WorkflowStepDefinition{
		Entity: NewEntity(context),
		Name:   context.Name,
	}
}

// tosca.Reader signature
func ReadWorkflowStepDefinition(context *tosca.Context) interface{} {
	self := NewWorkflowStepDefinition(context)
	// TODO: "operation_host"
	context.ValidateUnsupportedFields(context.ReadFields(self, Readers))
	return self
}

// tosca.Mappable interface
func (self *WorkflowStepDefinition) GetKey() string {
	return self.Name
}

//
// WorkflowStepDefinitions
//

type WorkflowStepDefinitions map[string]*WorkflowStepDefinition

func (self WorkflowStepDefinitions) Render(context *tosca.Context) {
	for _, step := range self {
		if step.OnSuccessStepNames != nil {
			for index, name := range *step.OnSuccessStepNames {
				if _, ok := self[name]; !ok {
					context.ListChild(index, name).ReportUnknown("step")
				}
			}
		}
	}
}
