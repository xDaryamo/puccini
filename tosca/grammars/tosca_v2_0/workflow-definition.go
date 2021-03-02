package tosca_v2_0

import (
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
)

//
// WorkflowDefinition
//
// [TOSCA-v2.0] @ ?
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.8.7
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.8.7
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.7.7
//

type WorkflowDefinition struct {
	*Entity `name:"workflow definition"`
	Name    string `namespace:""`

	Metadata                Metadata                        `read:"metadata,Metadata"`
	Description             *string                         `read:"description"`
	InputDefinitions        PropertyDefinitions             `read:"inputs,PropertyDefinition"`
	PreconditionDefinitions WorkflowPreconditionDefinitions `read:"preconditions,WorkflowPreconditionDefinition"`
	StepDefinitions         WorkflowStepDefinitions         `read:"steps,WorkflowStepDefinition"`
}

func NewWorkflowDefinition(context *tosca.Context) *WorkflowDefinition {
	return &WorkflowDefinition{
		Entity:           NewEntity(context),
		Name:             context.Name,
		InputDefinitions: make(PropertyDefinitions),
		StepDefinitions:  make(WorkflowStepDefinitions),
	}
}

// tosca.Reader signature
func ReadWorkflowDefinition(context *tosca.Context) tosca.EntityPtr {
	self := NewWorkflowDefinition(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))
	return self
}

// tosca.Mappable interface
func (self *WorkflowDefinition) GetKey() string {
	return self.Name
}

// parser.Renderable interface
func (self *WorkflowDefinition) Render() {
	logRender.Debugf("workflow definition: %s", self.Name)

	self.StepDefinitions.Render()
}

func (self *WorkflowDefinition) Normalize(normalServiceTemplate *normal.ServiceTemplate) *normal.Workflow {
	logNormalize.Debugf("workflow definition: %s", self.Name)

	normalWorkflow := normalServiceTemplate.NewWorkflow(self.Name)

	if self.Description != nil {
		normalWorkflow.Description = *self.Description
	}

	// TODO: support property definitions
	//self.InputDefinitions.Normalize(w.Inputs)

	self.StepDefinitions.Normalize(normalWorkflow)

	return normalWorkflow
}

//
// WorkflowDefinitions
//

type WorkflowDefinitions map[string]*WorkflowDefinition

func (self WorkflowDefinitions) Normalize(normalServiceTemplate *normal.ServiceTemplate) {
	for _, workflowDefinition := range self {
		normalServiceTemplate.Workflows[workflowDefinition.Name] = workflowDefinition.Normalize(normalServiceTemplate)
	}
}
