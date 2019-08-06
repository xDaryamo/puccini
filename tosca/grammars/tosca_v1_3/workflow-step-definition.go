package tosca_v1_3

import (
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
)

//
// WorkflowStepDefinition
//
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.23
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.21
//

type WorkflowStepDefinition struct {
	*Entity `name:"workflow step definition"`
	Name    string

	TargetNodeTemplateOrGroupName *string                     `read:"target" require:"target"`
	TargetNodeRequirementName     *string                     `read:"target_relationship"`
	OperationHost                 *string                     `read:"operation_host"`
	FilterConstraintClauses       ConstraintClauses           `read:"filter,[]ConstraintClause"`
	ActivityDefinitions           WorkflowActivityDefinitions `read:"activities,[]WorkflowActivityDefinition" require:"activities"`
	OnSuccessStepNames            *[]string                   `read:"on_success"`
	OnFailureStepNames            *[]string                   `read:"on_failure"`

	TargetNodeTemplate *NodeTemplate             `lookup:"target,TargetNodeTemplateOrGroupName" json:"-" yaml:"-"`
	TargetGroup        *Group                    `lookup:"target,TargetNodeTemplateOrGroupName" json:"-" yaml:"-"`
	OnSuccessSteps     []*WorkflowStepDefinition // custom lookup
	OnFailureSteps     []*WorkflowStepDefinition // custom lookup
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
	context.ValidateUnsupportedFields(context.ReadFields(self))
	return self
}

// tosca.Mappable interface
func (self *WorkflowStepDefinition) GetKey() string {
	return self.Name
}

func (self *WorkflowStepDefinition) Render(definitions WorkflowStepDefinitions) {
	if self.OnSuccessStepNames != nil {
		for index, name := range *self.OnSuccessStepNames {
			if s, ok := definitions[name]; ok {
				self.OnSuccessSteps = append(self.OnSuccessSteps, s)
			} else {
				self.Context.ListChild(index, name).ReportUnknown("step")
			}
		}
	}

	if self.OnFailureStepNames != nil {
		for index, name := range *self.OnFailureStepNames {
			if s, ok := definitions[name]; ok {
				self.OnFailureSteps = append(self.OnFailureSteps, s)
			} else {
				self.Context.ListChild(index, name).ReportUnknown("step")
			}
		}
	}

	for _, activity := range self.ActivityDefinitions {
		activity.Render(self)
	}

	// TODO: validate OperationHost
}

func (self *WorkflowStepDefinition) Normalize(w *normal.Workflow, s *normal.ServiceTemplate) *normal.WorkflowStep {
	log.Infof("{normalize} workflow step: %s", self.Name)

	st := w.NewStep(self.Name)

	if self.TargetNodeTemplate != nil {
		if n, ok := s.NodeTemplates[self.TargetNodeTemplate.Name]; ok {
			st.TargetNodeTemplate = n
		}
	} else if self.TargetGroup != nil {
		if g, ok := s.Groups[self.TargetGroup.Name]; ok {
			st.TargetGroup = g
		}
	}

	for _, activity := range self.ActivityDefinitions {
		activity.Normalize(st, s)
	}

	return st
}

func (self *WorkflowStepDefinition) NormalizeNext(st *normal.WorkflowStep, w *normal.Workflow) {
	for _, next := range self.OnSuccessSteps {
		st.OnSuccessSteps = append(st.OnSuccessSteps, w.Steps[next.Name])
	}

	for _, next := range self.OnFailureSteps {
		st.OnFailureSteps = append(st.OnFailureSteps, w.Steps[next.Name])
	}

	if self.OperationHost != nil {
		st.Host = *self.OperationHost
	}
}

//
// WorkflowStepDefinitions
//

type WorkflowStepDefinitions map[string]*WorkflowStepDefinition

func (self WorkflowStepDefinitions) Render() {
	for _, step := range self {
		step.Render(self)
	}
}

func (self WorkflowStepDefinitions) Normalize(w *normal.Workflow, s *normal.ServiceTemplate) {
	steps := make(normal.WorkflowSteps)
	for name, step := range self {
		steps[name] = step.Normalize(w, s)
	}
	for name, step := range self {
		step.NormalizeNext(steps[name], w)
	}
}
