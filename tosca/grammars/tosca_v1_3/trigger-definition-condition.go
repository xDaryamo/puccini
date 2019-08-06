package tosca_v1_3

import (
	"github.com/tliron/puccini/tosca"
)

//
// TriggerDefinitionCondition
//

type TriggerDefinitionCondition struct {
	*Entity `name:"trigger definition condition" json:"-" yaml:"-"`

	ConstraintClauses ConstraintClauses `read:"constraint,[]ConstraintClause"` // this should be "constraints"...
	Period            *ScalarUnitTime   `read:"period,scalar-unit.time"`
	Evaluations       *int              `read:"evaluations"`
	Method            *string           `read:"method"`
}

func NewTriggerDefinitionCondition(context *tosca.Context) *TriggerDefinitionCondition {
	return &TriggerDefinitionCondition{Entity: NewEntity(context)}
}

// tosca.Reader signature
func ReadTriggerDefinitionCondition(context *tosca.Context) interface{} {
	self := NewTriggerDefinitionCondition(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))
	return self
}
