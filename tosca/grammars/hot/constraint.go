package hot

import (
	"github.com/tliron/kutil/ard"
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
	profile "github.com/tliron/puccini/tosca/profiles/hot/v1_0"
	"github.com/tliron/yamlkeys"
)

const constraintPathPrefix = "/hot/1.0/js/constraints/"

// Built-in constraint functions
var ConstraintScriptlets = map[string]string{
	tosca.ConstraintScriptletPrefix + "length":            profile.Profile[constraintPathPrefix+"length.js"],
	tosca.ConstraintScriptletPrefix + "range":             profile.Profile[constraintPathPrefix+"range.js"],
	tosca.ConstraintScriptletPrefix + "modulo":            profile.Profile[constraintPathPrefix+"modulo.js"],
	tosca.ConstraintScriptletPrefix + "allowed_values":    profile.Profile[constraintPathPrefix+"allowed_values.js"],
	tosca.ConstraintScriptletPrefix + "allowed_pattern":   profile.Profile[constraintPathPrefix+"allowed_pattern.js"],
	tosca.ConstraintScriptletPrefix + "custom_constraint": profile.Profile[constraintPathPrefix+"custom_constraint.js"],
}

var ConstraintNativeArgumentIndexes = map[string][]int{}

//
// Constraint
//
// [https://docs.openstack.org/heat/wallaby/template_guide/hot_spec.html#parameter-constraints]
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
			context.ReportValueMalformed("constraint", "map size not 1 or 2")
			return self
		}

		for key, value := range map_ {
			operator := yamlkeys.KeyString(key)

			if operator == "description" {
				self.Description = context.FieldChild(operator, value).ReadString()
				continue
			}

			scriptletName := tosca.ConstraintScriptletPrefix + operator
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
	return context.NewFunctionCall(tosca.ConstraintScriptletPrefix+self.Operator, self.Arguments)
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
