package tosca_v1_3

import (
	"github.com/tliron/puccini/tosca"
)

//
// ConditionClause
//
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.21
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.19
//

type ConditionClause struct {
	*Entity `name:"condition clause"`
}

func NewConditionClause(context *tosca.Context) *ConditionClause {
	return &ConditionClause{Entity: NewEntity(context)}
}

// tosca.Reader signature
func ReadConditionClause(context *tosca.Context) interface{} {
	self := NewConditionClause(context)

	if context.ValidateType("map") {
		for _, childContext := range context.FieldChildren() {
			if !self.readField(childContext) {
				childContext.ReportFieldUnsupported()
			}
		}
	}

	return self
}

func (self *ConditionClause) readField(context *tosca.Context) bool {
	switch context.Name {
	case "and":
	case "or":
	case "assert":
	default:
		return false
	}
	return true
}

//
// ConditionClauses
//

type ConditionClauses []*ConditionClause
