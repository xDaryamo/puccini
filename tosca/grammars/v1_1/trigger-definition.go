package v1_1

import (
	"github.com/tliron/puccini/tosca"
)

//
// TriggerDefinition
//

type TriggerDefinition struct {
	*Entity `name:"trigger definition" json:"-" yaml:"-"`
	Name    string

	Description  *string                  `read:"description"`
	EventType    *string                  `read:"event_type" required:"event_type"`
	Schedule     *Value                   `read:"schedule,Value"` // tosca.datatypes.TimeInterval
	TargetFilter *EventFilter             `read:"target_filter,EventFilter"`
	Condition    *string                  `read:"condition"`  // TODO: ??
	Constraint   *string                  `read:"constraint"` // TODO: ??
	Period       *ScalarUnitTime          `read:"period,scalar-unit.time"`
	Evaluations  *int                     `read:"evaluations"`
	Method       *string                  `read:"method"`
	Action       *OperationImplementation `read:"action,OperationImplementation" required:"action"`
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
	context.ReadFields(self, Readers)
	return self
}

// tosca.Renderable interface
func (self *TriggerDefinition) Render() {
	log.Infof("{render} trigger definition: %s", self.Name)
	if self.Schedule != nil {
		self.Schedule.RenderDataType("tosca.datatypes.TimeInterval")
	}
}

// tosca.Mappable interface
func (self *TriggerDefinition) GetKey() string {
	return self.Name
}

//
// TriggerDefinitions
//

type TriggerDefinitions map[string]*TriggerDefinition
