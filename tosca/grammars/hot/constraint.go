package hot

import (
	"github.com/tliron/kutil/ard"
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
	profile "github.com/tliron/puccini/tosca/profiles/hot/v1_0"
	"github.com/tliron/yamlkeys"
)

// Built-in constraint functions
var ConstraintScriptlets = map[string]string{
	"tosca.constraint.length":            profile.Profile["/hot/1.0/js/constraints/length.js"],
	"tosca.constraint.range":             profile.Profile["/hot/1.0/js/constraints/range.js"],
	"tosca.constraint.modulo":            profile.Profile["/hot/1.0/js/constraints/modulo.js"],
	"tosca.constraint.allowed_values":    profile.Profile["/hot/1.0/js/constraints/allowed_values.js"],
	"tosca.constraint.allowed_pattern":   profile.Profile["/hot/1.0/js/constraints/allowed_pattern.js"],
	"tosca.constraint.custom_constraint": profile.Profile["/hot/1.0/js/constraints/custom_constraint.js"],
}

var ConstraintNativeArgumentIndexes = map[string][]uint{}

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
func ReadConstraint(context *tosca.Context) tosca.EntityPtr {
	self := NewConstraint(context)

	if context.ValidateType(ard.TypeMap) {
		map_ := context.Data.(ard.Map)
		var length = len(map_)
		if (length != 1) && (length != 2) {
			context.ReportValueMalformed("constraint", "map length not 1 or 2")
			return self
		}

		for key, value := range map_ {
			operator := yamlkeys.KeyString(key)

			if operator == "description" {
				self.Description = context.FieldChild(operator, value).ReadString()
				continue
			}

			scriptletName := "tosca.constraint." + operator
			if _, ok := context.ScriptletNamespace.Lookup(scriptletName); !ok {
				context.Clone(operator).ReportValueMalformed("constraint", "unsupported operator")
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
	return context.NewFunctionCall("tosca.constraint."+self.Operator, self.Arguments)
}

//
// Constraints
//

type Constraints []*Constraint

func (self Constraints) Normalize(context *tosca.Context, normalConstrainable normal.Constrainable) {
	for _, constraint := range self {
		functionCall := constraint.NewFunctionCall(context)
		NormalizeFunctionCallArguments(functionCall, context)
		normalConstrainable.AddConstraint(functionCall)
		// TODO: normalize constraint description somewhere
	}
}
