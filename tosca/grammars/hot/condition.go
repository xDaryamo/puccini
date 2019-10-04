package hot

import (
	"github.com/tliron/puccini/ard"
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/yamlkeys"
)

//
// Condition
//

type Condition struct {
	*Entity `name:"condition"`

	Operator              string
	Arguments             ard.List
	NativeArgumentIndexes []uint

	DefinitionName *string

	Value *bool
}

func NewCondition(context *tosca.Context) *Condition {
	return &Condition{Entity: NewEntity(context)}
}

// tosca.Reader signature
func ReadCondition(context *tosca.Context) interface{} {
	self := NewCondition(context)

	if context.Is("map") {
		map_ := context.Data.(ard.Map)
		if len(map_) != 1 {
			context.ReportValueMalformed("condition", "map length not 1")
			return self
		}

		for key, value := range map_ {
			operator := yamlkeys.KeyString(key)

			script, ok := context.ScriptletNamespace[operator]
			if !ok {
				context.Clone(operator).ReportValueMalformed("condition", "unsupported operator")
				return self
			}

			self.Operator = operator

			if list, ok := value.(ard.List); ok {
				self.Arguments = list
			} else {
				self.Arguments = ard.List{value}
			}

			self.NativeArgumentIndexes = script.NativeArgumentIndexes

			// We have only one key
			break
		}
	} else if context.Is("string") {
		self.DefinitionName = context.ReadString()
	} else if context.ValidateType("map", "string", "bool") {
		self.Value = context.ReadBoolean()
	}

	return self
}
