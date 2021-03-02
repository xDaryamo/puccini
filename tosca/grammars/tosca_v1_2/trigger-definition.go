package tosca_v1_2

import (
	"github.com/tliron/kutil/ard"
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/grammars/tosca_v2_0"
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
	*tosca_v2_0.Entity `name:"trigger definition" json:"-" yaml:"-"`
	Name               string

	Description     *string                                `read:"description"`
	Event           *string                                `read:"event_type"`
	Schedule        *tosca_v2_0.Value                      `read:"schedule,Value"` // tosca:TimeInterval
	TargetFilter    *tosca_v2_0.EventFilter                `read:"target_filter,EventFilter"`
	Condition       *tosca_v2_0.TriggerDefinitionCondition `read:"condition,TriggerDefinitionCondition"`
	OperationAction *tosca_v2_0.OperationDefinition
	WorkflowAction  *string

	WorkflowDefinition *tosca_v2_0.WorkflowDefinition `lookup:"action,WorkflowAction"`
}

func NewTriggerDefinition(context *tosca.Context) *TriggerDefinition {
	return &TriggerDefinition{
		Entity: tosca_v2_0.NewEntity(context),
		Name:   context.Name,
	}
}

// tosca.Reader signature
func ReadTriggerDefinition(context *tosca.Context) tosca.EntityPtr {
	self := NewTriggerDefinition(context)
	context.ValidateUnsupportedFields(append(context.ReadFields(self), "action"))

	if context.Is(ard.TypeMap) {
		map_ := context.Data.(ard.Map)
		childContext := context.FieldChild("action", nil)
		var ok bool
		if childContext.Data, ok = map_["action"]; ok {
			if childContext.ValidateType(ard.TypeMap, ard.TypeString) {
				if childContext.Is(ard.TypeMap) {
					// Note that OperationDefinition can also be a string, but there is no way
					// for us to differentiate between that and a workflow ID, so we support only
					// the long notation
					self.OperationAction = tosca_v2_0.ReadOperationDefinition(childContext).(*tosca_v2_0.OperationDefinition)
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

// parser.Renderable interface
func (self *TriggerDefinition) Render() {
	logRender.Debugf("trigger definition: %s", self.Name)
	if self.Schedule != nil {
		self.Schedule.RenderDataType("tosca:TimeInterval")
	}
}

func (self *TriggerDefinition) Normalize(normalPolicy *normal.Policy) *normal.PolicyTrigger {
	normalPolicyTrigger := normalPolicy.NewTrigger()

	if self.OperationAction != nil {
		self.OperationAction.Normalize(normalPolicyTrigger.NewOperation())
	} else if self.WorkflowDefinition != nil {
		normalPolicyTrigger.Workflow = normalPolicy.ServiceTemplate.Workflows[self.WorkflowDefinition.Name]
	}

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
