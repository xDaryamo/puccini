package hot

import (
	"github.com/tliron/puccini/ard"
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
	profile "github.com/tliron/puccini/tosca/profiles/hot/v1_0"
)

// Built-in constraint functions
var ConstraintSourceCode = map[string]string{
	"length":            profile.Profile["/hot/1.0/js/length.js"],
	"range":             profile.Profile["/hot/1.0/js/range.js"],
	"modulo":            profile.Profile["/hot/1.0/js/modulo.js"],
	"allowed_values":    profile.Profile["/hot/1.0/js/allowed_values.js"],
	"allowed_pattern":   profile.Profile["/hot/1.0/js/allowed_pattern.js"],
	"custom_constraint": profile.Profile["/hot/1.0/js/custom_constraint.js"],
}

//
// Constraint
//
// [https://docs.openstack.org/heat/stein/template_guide/hot_spec.html#parameter-constraints]
//

type Constraint struct {
	*Entity `name:"constraint"`

	Description *string
	Operator    string
	Arguments   ard.List
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

			if _, ok := context.ScriptNamespace[operator]; !ok {
				context.WithData(operator).ReportValueMalformed("constraint", "unsupported operator")
				return self
			}

			self.Operator = operator

			if list, ok := value.(ard.List); ok {
				self.Arguments = list
			} else {
				self.Arguments = ard.List{value}
			}
		}
	}

	return self
}

func (self *Constraint) NewFunctionCall(context *tosca.Context) *tosca.FunctionCall {
	return context.NewFunctionCall(self.Operator, self.Arguments)
}

//
// Constraints
//

type Constraints []*Constraint

func (self Constraints) Normalize(context *tosca.Context, constrainable normal.Constrainable) {
	for _, constraint := range self {
		functionCall := constraint.NewFunctionCall(context)
		NormalizeFunctionCallArguments(functionCall, context)
		constrainable.AddConstraint(functionCall)
		// TODO: normalize constraint description somewhere
	}
}
