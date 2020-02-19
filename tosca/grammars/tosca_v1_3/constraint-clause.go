package tosca_v1_3

import (
	"github.com/tliron/puccini/ard"
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
	profile "github.com/tliron/puccini/tosca/profiles/simple/v1_3"
	"github.com/tliron/yamlkeys"
)

// Built-in constraint functions
var ConstraintClauseScriptlets = map[string]string{
	"tosca.constraint.equal":            profile.Profile["/tosca/simple/1.3/js/constraints/equal.js"],
	"tosca.constraint.greater_than":     profile.Profile["/tosca/simple/1.3/js/constraints/greater_than.js"],
	"tosca.constraint.greater_or_equal": profile.Profile["/tosca/simple/1.3/js/constraints/greater_or_equal.js"],
	"tosca.constraint.less_than":        profile.Profile["/tosca/simple/1.3/js/constraints/less_than.js"],
	"tosca.constraint.less_or_equal":    profile.Profile["/tosca/simple/1.3/js/constraints/less_or_equal.js"],
	"tosca.constraint.in_range":         profile.Profile["/tosca/simple/1.3/js/constraints/in_range.js"],
	"tosca.constraint.valid_values":     profile.Profile["/tosca/simple/1.3/js/constraints/valid_values.js"],
	"tosca.constraint.length":           profile.Profile["/tosca/simple/1.3/js/constraints/length.js"],
	"tosca.constraint.min_length":       profile.Profile["/tosca/simple/1.3/js/constraints/min_length.js"],
	"tosca.constraint.max_length":       profile.Profile["/tosca/simple/1.3/js/constraints/max_length.js"],
	"tosca.constraint.pattern":          profile.Profile["/tosca/simple/1.3/js/constraints/pattern.js"],
	"tosca.constraint.schema":           profile.Profile["/tosca/simple/1.3/js/constraints/schema.js"], // introduced in TOSCA 1.3
}

var ConstraintClauseNativeArgumentIndexes = map[string][]uint{
	"tosca.constraint.equal":            {0},
	"tosca.constraint.greater_than":     {0},
	"tosca.constraint.greater_or_equal": {0},
	"tosca.constraint.less_than":        {0},
	"tosca.constraint.less_or_equal":    {0},
	"tosca.constraint.in_range":         {0, 1},
}

//
// ConstraintClause
//
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.6.3
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.3
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.2
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.5.2
//

type ConstraintClause struct {
	*Entity `name:"constraint clause"`

	Operator              string
	Arguments             ard.List
	NativeArgumentIndexes []uint
	DataType              *DataType `traverse:"ignore" json:"-" yaml:"-"`
}

func NewConstraintClause(context *tosca.Context) *ConstraintClause {
	return &ConstraintClause{Entity: NewEntity(context)}
}

// tosca.Reader signature
func ReadConstraintClause(context *tosca.Context) interface{} {
	self := NewConstraintClause(context)

	if context.ValidateType("map") {
		map_ := context.Data.(ard.Map)
		if len(map_) != 1 {
			context.ReportValueMalformed("constraint clause", "map length not 1")
			return self
		}

		for key, value := range map_ {
			operator := yamlkeys.KeyString(key)

			scriptletName := "tosca.constraint." + operator
			scriptlet, ok := context.ScriptletNamespace[scriptletName]
			if !ok {
				context.Clone(operator).ReportValueMalformed("constraint clause", "unsupported operator")
				return self
			}

			self.Operator = operator

			if list, ok := value.(ard.List); ok {
				self.Arguments = list
			} else {
				self.Arguments = ard.List{value}
			}

			self.NativeArgumentIndexes = scriptlet.NativeArgumentIndexes

			// We have only one key
			break
		}
	}

	return self
}

func (self *ConstraintClause) ToFunctionCall(context *tosca.Context, strict bool) *tosca.FunctionCall {
	arguments := make([]interface{}, len(self.Arguments))
	for index, argument := range self.Arguments {
		if self.IsNativeArgument(uint(index)) {
			if _, ok := argument.(*Value); !ok {
				if self.DataType != nil {
					value := ReadValue(context.ListChild(index, argument)).(*Value)
					value.RenderAttribute(self.DataType, nil, true, false)
					argument = value
				} else if strict {
					panic("no data type for native argument")
				}
			}
		}
		arguments[index] = argument
	}
	return context.NewFunctionCall("tosca.constraint."+self.Operator, arguments)
}

func (self *ConstraintClause) IsNativeArgument(index uint) bool {
	for _, i := range self.NativeArgumentIndexes {
		if i == index {
			return true
		}
	}
	return false
}

//
// ConstraintClauses
//

type ConstraintClauses []*ConstraintClause

func (self ConstraintClauses) RenderAndAppend(constraints *ConstraintClauses, dataType *DataType) {
	for _, constraintClause := range self {
		if (constraintClause.DataType != nil) && (constraintClause.DataType != dataType) {
			panic("constraint clause cannot be used with different data type")
		}
		constraintClause.DataType = dataType
		*constraints = append(*constraints, constraintClause)
	}
}

func (self ConstraintClauses) Normalize(context *tosca.Context) normal.FunctionCalls {
	var functionCalls normal.FunctionCalls
	for _, constraintClause := range self {
		functionCall := constraintClause.ToFunctionCall(context, false)
		NormalizeFunctionCallArguments(functionCall, context)
		functionCalls = append(functionCalls, normal.NewFunctionCall(functionCall))
	}
	return functionCalls
}

func (self ConstraintClauses) NormalizeConstrainable(context *tosca.Context, constrainable normal.Constrainable) {
	for _, constraintClause := range self {
		functionCall := constraintClause.ToFunctionCall(context, true)
		NormalizeFunctionCallArguments(functionCall, context)
		constrainable.AddConstraint(functionCall)
	}
}

func (self ConstraintClauses) NormalizeListEntries(context *tosca.Context, l *normal.List) {
	for _, constraintClause := range self {
		functionCall := constraintClause.ToFunctionCall(context, true)
		NormalizeFunctionCallArguments(functionCall, context)
		l.AddEntryConstraint(functionCall)
	}
}

func (self ConstraintClauses) NormalizeMapKeys(context *tosca.Context, m *normal.Map) {
	for _, constraintClause := range self {
		functionCall := constraintClause.ToFunctionCall(context, true)
		NormalizeFunctionCallArguments(functionCall, context)
		m.AddKeyConstraint(functionCall)
	}
}

func (self ConstraintClauses) NormalizeMapValues(context *tosca.Context, m *normal.Map) {
	for _, constraintClause := range self {
		functionCall := constraintClause.ToFunctionCall(context, true)
		NormalizeFunctionCallArguments(functionCall, context)
		m.AddValueConstraint(functionCall)
	}
}
