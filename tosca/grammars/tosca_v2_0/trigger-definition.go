package tosca_v2_0

import (
	"github.com/tliron/puccini/normal"
	"github.com/tliron/puccini/tosca/parsing"
)

//
// TriggerDefinition
//
// [TOSCA-v2.0] @ 16.5
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.6.22
//

type TriggerDefinition struct {
	*Entity `name:"trigger definition" json:"-" yaml:"-"`
	Name    string

	Description *string                     `read:"description"`
	Event       *string                     `read:"event" mandatory:""`
	Condition   *ValidationClause           `read:"condition,ValidationClause"`
	Action      WorkflowActivityDefinitions `read:"action,[]WorkflowActivityDefinition" mandatory:""`
}

func NewTriggerDefinition(context *parsing.Context) *TriggerDefinition {
	return &TriggerDefinition{
		Entity: NewEntity(context),
		Name:   context.Name,
	}
}

// ([parsing.Reader] signature)
func ReadTriggerDefinition(context *parsing.Context) parsing.EntityPtr {
	self := NewTriggerDefinition(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))
	return self
}

// ([parsing.Mappable] interface)
func (self *TriggerDefinition) GetKey() string {
	return self.Name
}

// ([parsing.Renderable] interface)
func (self *TriggerDefinition) Render() {
	// Avoid rendering more than once
	self.renderOnce.Do(self.render)
}

func (self *TriggerDefinition) render() {
	logRender.Debugf("trigger definition: %s", self.Name)
}

func (self *TriggerDefinition) Normalize(normalPolicy *normal.Policy) *normal.PolicyTrigger {
	normalPolicyTrigger := normalPolicy.NewTrigger()

	if self.Description != nil {
		normalPolicyTrigger.Description = *self.Description
	}

	if self.Event != nil {
		normalPolicyTrigger.Event = *self.Event
		// Maintain backward compatibility
		normalPolicyTrigger.EventType = *self.Event
	}

	// Normalize condition as boolean FunctionCall
	if self.Condition != nil {
		fc := self.Condition.ToFunctionCall(self.Context, false)
		NormalizeFunctionCallArguments(fc, self.Context)
		normalPolicyTrigger.Condition = normal.NewFunctionCall(fc)
	}

	// TODO: normalize actions

	return normalPolicyTrigger
}

//
// TriggerDefinitions
//

type TriggerDefinitions map[string]*TriggerDefinition

func (self TriggerDefinitions) Normalize(normalPolicy *normal.Policy) {
	for _, triggerDefinition := range self {
		triggerDefinition.Normalize(normalPolicy)
	}
}
