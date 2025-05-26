package tosca_v2_0

import (
	"fmt"

	"github.com/tliron/go-ard"
	"github.com/tliron/puccini/assets/tosca/profiles"
	"github.com/tliron/puccini/normal"
	"github.com/tliron/puccini/tosca/parsing"
	"github.com/tliron/yamlkeys"
)

const validationPathPrefix = "implicit/2.0/js/constraints/"

// Built-in validation functions - reuse existing constraint JavaScript implementations where possible
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
}

var ValidationClauseNativeArgumentIndexes = map[string][]int{
	parsing.MetadataValidationPrefix + "equal":            {0},
	parsing.MetadataValidationPrefix + "greater_than":     {0},
	parsing.MetadataValidationPrefix + "greater_or_equal": {0},
	parsing.MetadataValidationPrefix + "less_than":        {0},
	parsing.MetadataValidationPrefix + "less_or_equal":    {0},
	parsing.MetadataValidationPrefix + "in_range":         {0, 1},
	parsing.MetadataValidationPrefix + "valid_values":     {-1}, // -1 means all
	parsing.MetadataValidationPrefix + "matches":          {1},
	parsing.MetadataValidationPrefix + "has_entry":        {0},
	parsing.MetadataValidationPrefix + "has_key":          {0},
	parsing.MetadataValidationPrefix + "has_all_entries":  {0},
	parsing.MetadataValidationPrefix + "has_all_keys":     {0},
	parsing.MetadataValidationPrefix + "has_any_entry":    {0},
	parsing.MetadataValidationPrefix + "has_any_key":      {0},
	parsing.MetadataValidationPrefix + "contains":         {0},
}

// ValidationClause represents a single validation constraint
type ValidationClause struct {
	*Entity `name:"validation clause"`

	Operator              string
	Arguments             ard.List
	NativeArgumentIndexes []int
	DataType              *DataType      `traverse:"ignore" json:"-" yaml:"-"`
	Definition            DataDefinition `traverse:"ignore" json:"-" yaml:"-"`
}

func NewValidationClause(context *parsing.Context) *ValidationClause {
	return &ValidationClause{Entity: NewEntity(context)}
}

// Reads a validation clause from context
func ReadValidationClause(context *parsing.Context) parsing.EntityPtr {
	self := NewValidationClause(context)

	if context.ValidateType(ard.TypeMap) {
		map_ := context.Data.(ard.Map)
		if len(map_) != 1 {
			context.ReportValueMalformed("validation clause", "map size not 1")
			return self
		}

		for key, value := range map_ {
			originalOperator := yamlkeys.KeyString(key)
			operator := originalOperator

			// Remove '$' prefix for internal processing (TOSCA 2.0)
			if len(operator) > 0 && operator[0] == '$' {
				operator = operator[1:]
			}

			scriptletName := parsing.MetadataValidationPrefix + operator
			scriptlet, ok := context.ScriptletNamespace.Lookup(scriptletName)
			if !ok {
				context.Clone(originalOperator).ReportValueMalformed("validation clause", "unsupported operator")
				return self
			}

			self.Operator = operator

			if list, ok := value.(ard.List); ok {
				self.Arguments = list
			} else {
				self.Arguments = ard.List{value}
			}

			self.NativeArgumentIndexes = scriptlet.NativeArgumentIndexes
			break // Only one key is allowed
		}
	}

	return self
}

// Recursively process an argument for a validation clause
func (self *ValidationClause) processValidationArgument(arg any, validationContext *parsing.Context) any {
	// Case 1: Argument is the string "$value"
	if strArg, ok := arg.(string); ok && strArg == "$value" {
		return validationContext.Data
	}

	// Case 2: Argument is a map with a single key
	if mapArg, ok := arg.(ard.Map); ok && len(mapArg) == 1 {
		var key, value any
		for k, v := range mapArg {
			key = k
			value = v
			break
		}
		keyStr := yamlkeys.KeyString(key)

		// Case 2a: Map is { "$value": [...] }
		if keyStr == "$value" {
			return dereferenceValuePath(validationContext.Data, value)
		}

		// Case 2b: Map could be a function { "$length": [...] }, { "$concat": [...] }, etc.
		if len(keyStr) > 0 && keyStr[0] == '$' {
			functionName := keyStr[1:] // Remove '$' prefix
			scriptletName := parsing.MetadataFunctionPrefix + functionName

			if _, isFunction := validationContext.ScriptletNamespace.Lookup(scriptletName); isFunction {
				// Recursively process function arguments
				var funcArgs ard.List
				if listVal, isList := value.(ard.List); isList {
					funcArgs = listVal
				} else {
					funcArgs = ard.List{value}
				}

				processedArgs := make([]any, len(funcArgs))
				for i, funcArg := range funcArgs {
					processedArgs[i] = self.processValidationArgument(funcArg, validationContext)
				}

				// Create a function call
				return validationContext.NewFunctionCall(scriptletName, processedArgs)
			}
		}
	}

	// Case 3: Literal value
	return arg
}

// Converts the validation clause into a function call for validation
func (self *ValidationClause) ToFunctionCall(context *parsing.Context, strict bool) *parsing.FunctionCall {
	processedArguments := make([]any, len(self.Arguments))
	for i, argument := range self.Arguments {
		processedArguments[i] = self.processValidationArgument(argument, context)

		// VALUTA LE FUNCTION CALL PRIMA DI PASSARLE ALLE VALIDAZIONI
		processedArguments[i] = self.evaluateNestedFunctions(processedArguments[i], context)

		// Normalize only if not already a FunctionCall
		if self.IsNativeArgument(i) {
			if _, isValue := processedArguments[i].(*Value); !isValue {
				if _, isFuncCall := processedArguments[i].(*parsing.FunctionCall); !isFuncCall {
					if self.DataType != nil {
						value := ReadValue(self.Context.ListChild(i, processedArguments[i])).(*Value)
						value.Render(self.DataType, self.Definition, true, false)
						processedArguments[i] = value
					} else if strict {
						panic(fmt.Sprintf("no data type for native argument at index %d", i))
					}
				}
			}
		}
	}

	return context.NewFunctionCall(parsing.MetadataValidationPrefix+self.Operator, processedArguments)
}

// Nuova funzione helper per valutare le funzioni annidate
func (self *ValidationClause) evaluateNestedFunctions(arg any, context *parsing.Context) any {
	// Se è già una FunctionCall, lascia che il sistema la gestisca automaticamente
	// Il trucco è che dobbiamo assicurarci che sia valutata prima di arrivare alle validazioni
	if functionCall, isFuncCall := arg.(*parsing.FunctionCall); isFuncCall {
		// Valuta ricorsivamente gli argomenti della funzione
		evaluatedArgs := make([]any, len(functionCall.Arguments))
		for i, arg := range functionCall.Arguments {
			evaluatedArgs[i] = self.evaluateNestedFunctions(arg, context)
		}

		// Crea una nuova function call con argomenti valutati
		return context.NewFunctionCall(functionCall.Name, evaluatedArgs)
	}

	// GESTIONE DIRETTA: Se l'argomento è ancora una mappa con funzione non risolta
	if mapArg, ok := arg.(ard.Map); ok && len(mapArg) == 1 {
		for key, value := range mapArg {
			keyStr := yamlkeys.KeyString(key)
			if len(keyStr) > 0 && keyStr[0] == '$' {
				functionName := keyStr[1:]
				scriptletName := parsing.MetadataFunctionPrefix + functionName

				if _, exists := context.ScriptletNamespace.Lookup(scriptletName); exists {
					// Valuta gli argomenti della funzione ricorsivamente
					var funcArgs ard.List
					if listVal, isList := value.(ard.List); isList {
						funcArgs = listVal
					} else {
						funcArgs = ard.List{value}
					}

					evaluatedArgs := make([]any, len(funcArgs))
					for i, funcArg := range funcArgs {
						evaluatedArgs[i] = self.evaluateNestedFunctions(funcArg, context)
					}

					// Crea una function call che dovrebbe essere valutata dal sistema
					return context.NewFunctionCall(scriptletName, evaluatedArgs)
				}
			}
		}
	}

	return arg
}

// Helper to dereference a path like { $value: [ "foo", 0, "bar" ] }
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

// Checks if the argument at the given index is native
func (self *ValidationClause) IsNativeArgument(index int) bool {
	for _, i := range self.NativeArgumentIndexes {
		if (i == -1) || (i == index) {
			return true
		}
	}
	return false
}

// ValidationClauses is a slice of ValidationClause pointers
type ValidationClauses []*ValidationClause

func (self ValidationClauses) Append(validations ValidationClauses) ValidationClauses {
	r := append(self[:0:0], self...)
	return append(r, validations...)
}

func (self ValidationClauses) Normalize(context *parsing.Context) normal.FunctionCalls {
	var normalFunctionCalls normal.FunctionCalls
	for _, validationClause := range self {
		functionCall := validationClause.ToFunctionCall(context, false)
		NormalizeFunctionCallArguments(functionCall, context)
		normalFunctionCalls = append(normalFunctionCalls, normal.NewFunctionCall(functionCall))
	}
	return normalFunctionCalls
}

func (self ValidationClauses) AddToMeta(context *parsing.Context, normalValueMeta *normal.ValueMeta) {
	for _, validationClause := range self {
		functionCall := validationClause.ToFunctionCall(context, true)
		NormalizeFunctionCallArguments(functionCall, context)
		normalValueMeta.AddValidator(functionCall)
	}
}
