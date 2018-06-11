package v1_1

import (
	"github.com/tliron/puccini/tosca"
)

//
// WorkflowDefinition
//

type WorkflowDefinition struct {
	*Entity `name:"workflow definition"`
	Name    string `namespace:""`

	Metadata                Metadata                          `read:"metadata,Metadata"`
	Description             *string                           `read:"description"`
	InputDefinitions        PropertyDefinitions               `read:"inputs,PropertyDefinition"`
	PreconditionDefinitions []*WorkflowPreconditionDefinition `read:"preconditions,WorkflowPreconditionDefinition"`
	StepDefinitions         WorkflowStepDefinitions           `read:"steps,WorkflowStepDefinition"`
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
func ReadWorkflowDefinition(context *tosca.Context) interface{} {
	self := NewWorkflowDefinition(context)
	context.ValidateUnsupportedFields(context.ReadFields(self, Readers))
	return self
}

func init() {
	Readers["WorkflowDefinition"] = ReadWorkflowDefinition
}

// tosca.Mappable interface
func (self *WorkflowDefinition) GetKey() string {
	return self.Name
}

// tosca.Renderable interface
func (self *WorkflowDefinition) Render() {
	log.Info("{render} workflow definition")

	self.StepDefinitions.Render(self.Context.FieldChild("steps", nil))
}

//
// WorkflowDefinitions
//

type WorkflowDefinitions map[string]*WorkflowDefinition
