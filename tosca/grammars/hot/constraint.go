package hot

import (
	"github.com/tliron/puccini/ard"
	"github.com/tliron/puccini/tosca"
	profile "github.com/tliron/puccini/tosca/profiles/hot/v2018_08_31"
)

// Built-in constraint functions
var ConstraintSourceCode = map[string]string{
	"length":            profile.Profile["/hot/2018-08-31/js/length.js"],
	"range":             profile.Profile["/hot/2018-08-31/js/range.js"],
	"modulo":            profile.Profile["/hot/2018-08-31/js/modulo.js"],
	"allowed_values":    profile.Profile["/hot/2018-08-31/js/allowed_values.js"],
	"allowed_pattern":   profile.Profile["/hot/2018-08-31/js/allowed_pattern.js"],
	"custom_constraint": profile.Profile["/hot/2018-08-31/js/custom_constraint.js"],
}

var ConstraintNativeArgumentIndexes = map[string][]uint{
	"range": {0, 1}, // TODO
}

//
// Constraint
//
// [https://docs.openstack.org/heat/rocky/template_guide/hot_spec.html#parameter-constraints]
//

type Constraint struct {
	*Entity `name:"constraint"`

	Description           *string
	Operator              string
	Arguments             ard.List
	NativeArgumentIndexes []uint
}

func NewConstraint(context *tosca.Context) *Constraint {
	return &Constraint{Entity: NewEntity(context)}
}

// tosca.Reader signature
func ReadConstraint(context *tosca.Context) interface{} {
	self := NewConstraint(context)

	if context.ValidateType("map") {
		map_ := context.Data.(ard.Map)
		var length = len(map_)
		if (length != 1) && (length != 2) {
			context.ReportValueMalformed("constraint", "map length not 1 or 2")
			return self
		}

		for operator, value := range map_ {
			if operator == "description" {
				self.Description = context.FieldChild(operator, value).ReadString()
				continue
			}

			script, ok := context.ScriptNamespace[operator]
			if !ok {
				context.WithData(operator).ReportValueMalformed("constraint", "unsupported operator")
				return self
			}

			self.Operator = operator

			if list, ok := value.(ard.List); ok {
				self.Arguments = list
			} else {
				self.Arguments = ard.List{value}
			}

			self.NativeArgumentIndexes = script.NativeArgumentIndexes
		}
	}

	return self
}

//
// Constraints
//

type Constraints []*Constraint
