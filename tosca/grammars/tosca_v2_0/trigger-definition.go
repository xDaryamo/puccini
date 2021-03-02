package tosca_v2_0

import (
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
)

//
// TriggerDefinition
//
// [TOSCA-v2.0] @ ?
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.6.22
//

type TriggerDefinition struct {
	*Entity `name:"trigger definition" json:"-" yaml:"-"`
	Name    string

	Description  *string                     `read:"description"`
	Event        *string                     `read:"event" require:""`
	Schedule     *Value                      `read:"schedule,Value"` // tosca:TimeInterval
	TargetFilter *EventFilter                `read:"target_filter,EventFilter"`
	Condition    *TriggerDefinitionCondition `read:"condition,TriggerDefinitionCondition"`
	Action       WorkflowActivityDefinitions `read:"action,[]WorkflowActivityDefinition" require:""`
}

func NewTriggerDefinition(context *tosca.Context) *TriggerDefinition {
	return &TriggerDefinition{
		Entity: NewEntity(context),
		Name:   context.Name,
	}
}

// tosca.Reader signature
func ReadTriggerDefinition(context *tosca.Context) tosca.EntityPtr {
	self := NewTriggerDefinition(context)
	context.ValidateUnsupportedFields(append(context.ReadFields(self)))
	return self
}

// tosca.Mappable interface
func (self *TriggerDefinition) GetKey() string {
	return self.Name
}

// parser.Renderable interface
func (self *TriggerDefinition) Render() {
	logRender.Debugf("trigger definition: %s", self.Name)
	if self.Schedule != nil {
		self.Schedule.RenderDataType("tosca:TimeInterval")
	}
}

func (self *TriggerDefinition) Normalize(normalPolicy *normal.Policy) *normal.PolicyTrigger {
	normalPolicyTrigger := normalPolicy.NewTrigger()

	// TODO: missing fields

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
