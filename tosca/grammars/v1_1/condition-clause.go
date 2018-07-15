package v1_1

import (
	"github.com/tliron/puccini/tosca"
)

//
// ConditionClause
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
