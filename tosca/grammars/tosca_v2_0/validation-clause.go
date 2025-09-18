package tosca_v2_0

import (
	"fmt"
	"strings"

	"github.com/tliron/go-ard"
	"github.com/tliron/puccini/assets/tosca/profiles"
	"github.com/tliron/puccini/normal"
	"github.com/tliron/puccini/tosca/parsing"
	"github.com/tliron/yamlkeys"
)

const validationPathPrefix = "implicit/2.0/js/constraints/"

// Built‑in validation functions — reuse existing JavaScript implementations when possible
var ValidationClauseScriptlets = map[string]string{
	parsing.MetadataValidationPrefix + "equal":            profiles.GetString(validationPathPrefix + "equal.js"),
	parsing.MetadataValidationPrefix + "greater_than":     profiles.GetString(validationPathPrefix + "greater_than.js"),
	parsing.MetadataValidationPrefix + "greater_or_equal": profiles.GetString(validationPathPrefix + "greater_or_equal.js"),
	parsing.MetadataValidationPrefix + "less_than":        profiles.GetString(validationPathPrefix + "less_than.js"),
	parsing.MetadataValidationPrefix + "less_or_equal":    profiles.GetString(validationPathPrefix + "less_or_equal.js"),
	parsing.MetadataValidationPrefix + "in_range":         profiles.GetString(validationPathPrefix + "in_range.js"),
	parsing.MetadataValidationPrefix + "valid_values":     profiles.GetString(validationPathPrefix + "valid_values.js"),
	parsing.MetadataValidationPrefix + "min_length":       profiles.GetString(validationPathPrefix + "min_length.js"),
	parsing.MetadataValidationPrefix + "max_length":       profiles.GetString(validationPathPrefix + "max_length.js"),
	parsing.MetadataValidationPrefix + "pattern":          profiles.GetString(validationPathPrefix + "pattern.js"),
	parsing.MetadataValidationPrefix + "schema":           profiles.GetString(validationPathPrefix + "schema.js"),
	parsing.MetadataValidationPrefix + "matches":          profiles.GetString(validationPathPrefix + "matches.js"),
	parsing.MetadataValidationPrefix + "has_suffix":       profiles.GetString(validationPathPrefix + "has_suffix.js"),
	parsing.MetadataValidationPrefix + "has_prefix":       profiles.GetString(validationPathPrefix + "has_prefix.js"),
	parsing.MetadataValidationPrefix + "has_entry":        profiles.GetString(validationPathPrefix + "has_entry.js"),
	parsing.MetadataValidationPrefix + "has_key":          profiles.GetString(validationPathPrefix + "has_key.js"),
	parsing.MetadataValidationPrefix + "has_all_entries":  profiles.GetString(validationPathPrefix + "has_all_entries.js"),
	parsing.MetadataValidationPrefix + "has_all_keys":     profiles.GetString(validationPathPrefix + "has_all_keys.js"),
	parsing.MetadataValidationPrefix + "has_any_entry":    profiles.GetString(validationPathPrefix + "has_any_entry.js"),
	parsing.MetadataValidationPrefix + "has_any_key":      profiles.GetString(validationPathPrefix + "has_any_key.js"),
	parsing.MetadataValidationPrefix + "contains":         profiles.GetString(validationPathPrefix + "contains.js"),
	parsing.MetadataValidationPrefix + "and":              profiles.GetString(validationPathPrefix + "and.js"),
	parsing.MetadataValidationPrefix + "or":               profiles.GetString(validationPathPrefix + "or.js"),
	parsing.MetadataValidationPrefix + "not":              profiles.GetString(validationPathPrefix + "not.js"),
	parsing.MetadataValidationPrefix + "xor":              profiles.GetString(validationPathPrefix + "xor.js"),
	parsing.MetadataValidationPrefix + "_format":          profiles.GetString(validationPathPrefix + "$format.js"),
}

// Indexes of arguments that must be converted to native Go types before evaluation
var ValidationClauseNativeArgumentIndexes = map[string][]int{
	parsing.MetadataValidationPrefix + "equal":            {0},
	parsing.MetadataValidationPrefix + "greater_than":     {0},
	parsing.MetadataValidationPrefix + "greater_or_equal": {0},
	parsing.MetadataValidationPrefix + "less_than":        {0},
	parsing.MetadataValidationPrefix + "less_or_equal":    {0},
	parsing.MetadataValidationPrefix + "in_range":         {0},
	parsing.MetadataValidationPrefix + "valid_values":     {0}, // Only first argument (value to test)
	parsing.MetadataValidationPrefix + "matches":          {1},
	parsing.MetadataValidationPrefix + "has_entry":        {0},
	parsing.MetadataValidationPrefix + "has_key":          {0},
	parsing.MetadataValidationPrefix + "has_all_entries":  {0},
	parsing.MetadataValidationPrefix + "has_all_keys":     {0},
	parsing.MetadataValidationPrefix + "has_any_entry":    {0},
	parsing.MetadataValidationPrefix + "has_any_key":      {0},
	parsing.MetadataValidationPrefix + "contains":         {0},
	parsing.MetadataValidationPrefix + "_format":          {0},
}

// ValidationClause represents a single validation constraint
type ValidationClause struct {
	*Entity `name:"validation clause"`

	Operator              string
	ScriptletName         string         // The actual scriptlet name found during parsing
	Arguments             ard.List
	NativeArgumentIndexes []int
	DataType              *DataType      `traverse:"ignore" json:"-" yaml:"-"`
	Definition            DataDefinition `traverse:"ignore" json:"-" yaml:"-"`
}

func NewValidationClause(context *parsing.Context) *ValidationClause {
	return &ValidationClause{Entity: NewEntity(context)}
}

// Read a validation clause from the parsing context
func ReadValidationClause(context *parsing.Context) parsing.EntityPtr {
	self := NewValidationClause(context)

	if context.ValidateType(ard.TypeMap) {
		m := context.Data.(ard.Map)
		if len(m) != 1 {
			context.ReportValueMalformed("validation clause", "map size not 1")
			return self
		}

		for key, value := range m {
			originalOp := yamlkeys.KeyString(key)
			op := strings.TrimPrefix(originalOp, "$") // remove leading '$' if present

			scriptletName := parsing.MetadataValidationPrefix + op
			scriptlet, ok := context.ScriptletNamespace.Lookup(scriptletName)
			if !ok {
				// Also try looking in the constraint namespace (for backward compatibility)
				constraintScriptletName := parsing.MetadataContraintPrefix + op
				scriptlet, ok = context.ScriptletNamespace.Lookup(constraintScriptletName)
				if !ok {
					context.Clone(originalOp).ReportValueMalformed("validation clause", "unsupported operator")
					return self
				}
				scriptletName = constraintScriptletName
			}

			self.Operator = op
			self.ScriptletName = scriptletName

			switch v := value.(type) {
			case ard.List:
				self.Arguments = v
			default:
				self.Arguments = ard.List{v}
			}

			self.NativeArgumentIndexes = scriptlet.NativeArgumentIndexes
			break
		}
	}

	return self
}

// Recursively process an argument, expanding "$value" and nested functions
func (self *ValidationClause) processValidationArgument(arg any, ctx *parsing.Context) any {
	// "$value" placeholder -> replace with the current value being validated
	if s, ok := arg.(string); ok && s == "$value" {
		// Special handling: if DataType is a scalar type, we're likely validating individual elements
		// In this case, keep $value as string to let JavaScript handle it properly
		if self.DataType != nil && self.DataType.IsScalarType() {
			return "$value"
		}
		return ctx.Data
	}

	// Auto-convert numeric values to strings with canonical units for scalar types
	if self.DataType != nil && self.DataType.IsScalarType() {
		if canonicalUnit := self.getCanonicalUnit(); canonicalUnit != "" {
			switch v := arg.(type) {
			case int:
				return fmt.Sprintf("%d %s", v, canonicalUnit)
			case int32:
				return fmt.Sprintf("%d %s", v, canonicalUnit)
			case int64:
				return fmt.Sprintf("%d %s", v, canonicalUnit)
			case float32:
				return fmt.Sprintf("%g %s", v, canonicalUnit)
			case float64:
				return fmt.Sprintf("%g %s", v, canonicalUnit)
			}
		}
	}

	// Map with exactly one key may represent nested function or dereference path
	if m, ok := arg.(ard.Map); ok && len(m) == 1 {
		var k, v any
		for key, val := range m {
			k, v = key, val
			break
		}
		keyStr := yamlkeys.KeyString(k)

		// Path dereference syntax: { "$value": ["foo", 0, "bar"] }
		if keyStr == "$value" {
			return dereferenceValuePath(ctx.Data, v)
		}

		// Nested function call: { "$length": [...] }, { "$concat": [...] }, etc.
		if strings.HasPrefix(keyStr, "$") {
			fnName := keyStr[1:]

			// Check if it's a validation clause first
			validationScriptletName := parsing.MetadataValidationPrefix + fnName
			_, exists := ctx.ScriptletNamespace.Lookup(validationScriptletName)
			if !exists {
				// Also try looking in the constraint namespace
				constraintScriptletName := parsing.MetadataContraintPrefix + fnName
				_, exists = ctx.ScriptletNamespace.Lookup(constraintScriptletName)
			}
			if exists {
				// It's a validation clause, leave as-is for proper processing
				return arg
			}

			// Then check if it's a function
			functionScriptletName := parsing.MetadataFunctionPrefix + fnName
			if _, isFn := ctx.ScriptletNamespace.Lookup(functionScriptletName); isFn {
				var fnArgs ard.List
				if listVal, isList := v.(ard.List); isList {
					fnArgs = listVal
				} else {
					fnArgs = ard.List{v}
				}

				processed := make(ard.List, len(fnArgs))
				for i, a := range fnArgs {
					processed[i] = self.processValidationArgument(a, ctx)
				}
				return ctx.NewFunctionCall(functionScriptletName, processed)
			}
		}
	}

	// Literal value, leave as‑is
	return arg
}

// Helper function to get canonical unit from DataType
func (self *ValidationClause) getCanonicalUnit() string {
	if self.DataType == nil || !self.DataType.IsScalarType() {
		return ""
	}

	// First check if CanonicalUnit is explicitly set
	if self.DataType.CanonicalUnit != nil {
		return *self.DataType.CanonicalUnit
	}

	// If no canonical unit specified, find the unit with multiplier 1
	if self.DataType.Units != nil {
		for unitInterface, multiplierInterface := range self.DataType.Units {
			if unitName, ok := unitInterface.(string); ok {
				// Try to convert multiplier to float64
				var multiplier float64
				switch m := multiplierInterface.(type) {
				case float64:
					multiplier = m
				case int:
					multiplier = float64(m)
				case int32:
					multiplier = float64(m)
				case int64:
					multiplier = float64(m)
				default:
					continue
				}

				if multiplier == 1.0 {
					return unitName
				}
			}
		}
	}

	return ""
}

// Convert the clause into a FunctionCall that Puccini can evaluate
func (self *ValidationClause) ToFunctionCall(ctx *parsing.Context, strict bool) *parsing.FunctionCall {
	processed := make(ard.List, len(self.Arguments))

	for i, arg := range self.Arguments {
		processed[i] = self.processValidationArgument(arg, ctx)
		processed[i] = self.evaluateNestedFunctions(processed[i], ctx)

		// Convert native arguments to *Value when a datatype is available
		if self.IsNativeArgument(i) {
			if _, isVal := processed[i].(*Value); !isVal {
				if _, isCall := processed[i].(*parsing.FunctionCall); !isCall {
					// Skip if it's already a ValidationClause (this can happen with nested validations)
					if _, isValidationClause := processed[i].(*ValidationClause); !isValidationClause {
						if self.DataType != nil {
							val := ReadValue(self.Context.ListChild(i, processed[i])).(*Value)
							val.Render(self.DataType, self.Definition, true, false)
							processed[i] = val
						} else if strict {
							panic(fmt.Sprintf("no data type for native argument at index %d", i))
						}
					}
				}
			}
		}
	}

	return ctx.NewFunctionCall(self.ScriptletName, processed)
}

// Evaluate nested functions contained inside an argument before validation
func (self *ValidationClause) evaluateNestedFunctions(arg any, ctx *parsing.Context) any {
	if fc, ok := arg.(*parsing.FunctionCall); ok {
		evaluated := make(ard.List, len(fc.Arguments))
		for i, a := range fc.Arguments {
			evaluated[i] = self.evaluateNestedFunctions(a, ctx)
		}
		return ctx.NewFunctionCall(fc.Name, evaluated)
	}

	// Handle map‑style function notation that hasn't yet been converted
	if m, ok := arg.(ard.Map); ok && len(m) == 1 {
		for k, v := range m {
			keyStr := yamlkeys.KeyString(k)
			if strings.HasPrefix(keyStr, "$") {
				fnName := keyStr[1:]

				// Check if it's a validation clause first
				validationScriptletName := parsing.MetadataValidationPrefix + fnName
				_, exists := ctx.ScriptletNamespace.Lookup(validationScriptletName)
				if !exists {
					// Also try looking in the constraint namespace
					constraintScriptletName := parsing.MetadataContraintPrefix + fnName
					_, exists = ctx.ScriptletNamespace.Lookup(constraintScriptletName)
				}
				if exists {
					// It's a validation clause, leave as-is for proper processing
					return arg
				}

				// Then check if it's a function
				functionScriptletName := parsing.MetadataFunctionPrefix + fnName
				if _, exists := ctx.ScriptletNamespace.Lookup(functionScriptletName); exists {
					var fnArgs ard.List
					if listVal, isList := v.(ard.List); isList {
						fnArgs = listVal
					} else {
						fnArgs = ard.List{v}
					}

					evaluated := make(ard.List, len(fnArgs))
					for i, a := range fnArgs {
						evaluated[i] = self.evaluateNestedFunctions(a, ctx)
					}
					return ctx.NewFunctionCall(functionScriptletName, evaluated)
				}
			}
		}
	}

	return arg
}

// Dereference a path such as { $value: ["foo", 0, "bar"] }
func dereferenceValuePath(data any, path any) any {
	list, ok := path.(ard.List)
	if !ok {
		return data
	}
	current := data
	for _, p := range list {
		switch c := current.(type) {
		case ard.Map:
			current = c[p]
		case map[string]any:
			if key, ok := p.(string); ok {
				current = c[key]
			}
		case []any:
			if idx, ok := p.(int); ok && idx < len(c) {
				current = c[idx]
			}
		default:
			return nil
		}
	}
	return current
}

// IsNativeArgument reports whether the argument at the given index must be treated as a native value
func (self *ValidationClause) IsNativeArgument(index int) bool {
	// First check if it's a $value token - these shouldn't be treated as native values
	if index < len(self.Arguments) {
		if s, ok := self.Arguments[index].(string); ok && s == "$value" {
			return false
		}

		// Check if it's a validation clause map - these shouldn't be treated as native values
		if m, ok := self.Arguments[index].(ard.Map); ok && len(m) == 1 {
			for key := range m {
				keyStr := yamlkeys.KeyString(key)
				if strings.HasPrefix(keyStr, "$") {
					fnName := keyStr[1:]
					validationScriptletName := parsing.MetadataValidationPrefix + fnName
					if _, exists := self.Context.ScriptletNamespace.Lookup(validationScriptletName); exists {
						// It's a validation clause, don't treat as native
						return false
					}
				}
			}
		}
	}

	// Then check if it's in the native argument indexes list
	for _, i := range self.NativeArgumentIndexes {
		if i == -1 || i == index {
			return true
		}
	}

	// Special handling for scalar types: only the value argument (usually index 0) should be native
	// This prevents bounds arrays and other structural arguments from being converted incorrectly
	if self.DataType != nil && self.DataType.IsScalarType() && index == 0 {
		return true
	}

	return false
}

// --- Collections ----------------------------------------------------------

type ValidationClauses []*ValidationClause

func (self ValidationClauses) Append(v ValidationClauses) ValidationClauses {
	out := append(self[:0:0], self...)
	return append(out, v...)
}

func (self ValidationClauses) Normalize(ctx *parsing.Context) normal.FunctionCalls {
	var calls normal.FunctionCalls
	for _, c := range self {
		fc := c.ToFunctionCall(ctx, false)
		NormalizeFunctionCallArguments(fc, ctx)
		calls = append(calls, normal.NewFunctionCall(fc))
	}
	return calls
}

func (self ValidationClauses) AddToMeta(ctx *parsing.Context, meta *normal.ValueMeta) {
	// Skip direct validation for lists when validation should be applied to elements
	if meta.Type == "list" && meta.Element != nil {
		// Instead of just returning, ensure validations are added to the element schema
		if len(self) > 0 && meta.Element != nil {
			for _, c := range self {
				fc := c.ToFunctionCall(ctx, true)
				NormalizeFunctionCallArguments(fc, ctx)
				meta.Element.AddValidator(fc)
			}
		}
		return
	}

	// Skip direct validation for maps when validation should be applied to values
	if meta.Type == "map" && meta.Value != nil {
		// Instead of just returning, ensure validations are added to the value schema
		if len(self) > 0 && meta.Value != nil {
			for _, c := range self {
				fc := c.ToFunctionCall(ctx, true)
				NormalizeFunctionCallArguments(fc, ctx)
				meta.Value.AddValidator(fc)
			}
		}
		return
	}

	for _, c := range self {
		fc := c.ToFunctionCall(ctx, true)
		NormalizeFunctionCallArguments(fc, ctx)
		meta.AddValidator(fc)
	}
}
