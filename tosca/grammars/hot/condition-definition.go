package hot

import (
	"github.com/tliron/puccini/ard"
	"github.com/tliron/puccini/tosca"
)

//
// ConditionDefinition
//
// [https://docs.openstack.org/heat/rocky/template_guide/hot_spec.html#conditions-section]
//

type ConditionDefinition struct {
	*Entity `name:"condition definition"`
	Name    string `namespace:""`

	Operator              string
	Arguments             ard.List
	NativeArgumentIndexes []uint
}

func NewConditionDefinition(context *tosca.Context) *ConditionDefinition {
	return &ConditionDefinition{
		Entity: NewEntity(context),
		Name:   context.Name,
	}
}

// tosca.Reader signature
func ReadConditionDefinition(context *tosca.Context) interface{} {
	self := NewConditionDefinition(context)

	if context.ValidateType("map") {
		map_ := context.Data.(ard.Map)
		if len(map_) != 1 {
			context.ReportValueMalformed("condition definition", "map length not 1")
			return self
		}

		for operator, value := range map_ {
			script, ok := context.ScriptNamespace[operator]
			if !ok {
				context.WithData(operator).ReportValueMalformed("condition definition", "unsupported operator")
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
	}

	return self
}
