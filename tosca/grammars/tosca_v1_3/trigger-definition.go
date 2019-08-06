package tosca_v1_3

import (
	"github.com/tliron/puccini/ard"
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
)

//
// TriggerDefinition
//
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.18
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.16
//

// Note: The TOSCA 1.1 spec is mangled, we will jump right to 1.2 here

type TriggerDefinition struct {
	*Entity `name:"trigger definition" json:"-" yaml:"-"`
	Name    string

	Description     *string                     `read:"description"`
	EventType       *string                     `read:"event_type"`
	Schedule        *Value                      `read:"schedule,Value"` // tosca.datatypes.TimeInterval
	TargetFilter    *EventFilter                `read:"target_filter,EventFilter"`
	Condition       *TriggerDefinitionCondition `read:"condition,TriggerDefinitionCondition"`
	Period          *ScalarUnitTime             `read:"period,scalar-unit.time"`
	Evaluations     *int                        `read:"evaluations"`
	Method          *string                     `read:"method"`
	OperationAction *OperationDefinition
	WorkflowAction  *string

	WorkflowDefinition *WorkflowDefinition `lookup:"action,WorkflowAction"`
}

func NewTriggerDefinition(context *tosca.Context) *TriggerDefinition {
	return &TriggerDefinition{
		Entity: NewEntity(context),
		Name:   context.Name,
	}
}

// tosca.Reader signature
func ReadTriggerDefinition(context *tosca.Context) interface{} {
	self := NewTriggerDefinition(context)
	context.ValidateUnsupportedFields(append(context.ReadFields(self), "action"))

	if context.Is("map") {
		map_ := context.Data.(ard.Map)
		childContext := context.FieldChild("action", nil)
		var ok bool
		if childContext.Data, ok = map_["action"]; ok {
			if childContext.ValidateType("map", "string") {
				if childContext.Is("map") {
					// Note that OperationDefinition can also be a string, but there is no way
					// for us to differentiate between that and a workflow ID, so we support only
					// the long notation
					self.OperationAction = ReadOperationDefinition(childContext).(*OperationDefinition)
				} else {
					self.WorkflowAction = childContext.ReadString()
				}
			}
		} else {
			childContext.ReportFieldMissing()
		}
	}

	return self
}

// tosca.Mappable interface
func (self *TriggerDefinition) GetKey() string {
	return self.Name
}

// tosca.Renderable interface
func (self *TriggerDefinition) Render() {
	log.Infof("{render} trigger definition: %s", self.Name)
	if self.Schedule != nil {
		self.Schedule.RenderDataType("tosca.datatypes.TimeInterval")
	}
}

func (self *TriggerDefinition) Normalize(p *normal.Policy, s *normal.ServiceTemplate) *normal.PolicyTrigger {
	t := p.NewTrigger()

	if self.OperationAction != nil {
		self.OperationAction.Normalize(t.NewOperation())
	} else if self.WorkflowDefinition != nil {
		t.Workflow = s.Workflows[self.WorkflowDefinition.Name]
	}

	// TODO: missing fields

	return t
}

//
// TriggerDefinitions
//

type TriggerDefinitions map[string]*TriggerDefinition

func (self TriggerDefinitions) Normalize(p *normal.Policy, s *normal.ServiceTemplate) {
	for _, triggerDefinition := range self {
		triggerDefinition.Normalize(p, s)
	}
}
