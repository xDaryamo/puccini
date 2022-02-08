package tosca_v2_0

import (
	"github.com/tliron/kutil/ard"
	"github.com/tliron/puccini/tosca"
)

//
// TriggerDefinitionCondition
//

type TriggerDefinitionCondition struct {
	*Entity `name:"trigger definition condition" json:"-" yaml:"-"`

	Constraint  *ConditionClause `read:"constraint,ConditionClause"` // why is this called "constraint"?
	Period      *Value           `read:"period,Value"`               // scalar-unit.time
	Evaluations *int             `read:"evaluations"`
	Method      *string          `read:"method"`
}

func NewTriggerDefinitionCondition(context *tosca.Context) *TriggerDefinitionCondition {
	return &TriggerDefinitionCondition{Entity: NewEntity(context)}
}

// tosca.Reader signature
func ReadTriggerDefinitionCondition(context *tosca.Context) tosca.EntityPtr {
	self := NewTriggerDefinitionCondition(context)

	if context.Is(ard.TypeMap) {
		map_ := context.Data.(ard.Map)
		if _, ok := map_["constraint"]; ok {
			// Long map notation
			context.ValidateUnsupportedFields(context.ReadFields(self))
		} else {
			// Short map notation
			self.Constraint = ReadConditionClause(context).(*ConditionClause)
		}
	} else if context.ValidateType(ard.TypeMap, ard.TypeList) {
		// Short list notation
		// See note: [TOSCA-Simple-Profile-YAML-v1.3] @ 3.6.25.1
		context = context.Clone(ard.Map{"and": context.Data})
		self.Constraint = ReadConditionClause(context).(*ConditionClause)
	}

	return self
}

// tosca.Renderable interface
func (self *TriggerDefinitionCondition) Render() {
	self.renderOnce.Do(self.render)
}

func (self *TriggerDefinitionCondition) render() {
	logRender.Debug("trigger definition condition")

	if self.Period != nil {
		self.Period.RenderDataType("scalar-unit.time")
	}
}
