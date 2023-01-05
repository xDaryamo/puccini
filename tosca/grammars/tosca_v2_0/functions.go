package tosca_v2_0

import (
	"strings"

	"github.com/tliron/kutil/ard"
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
	profile "github.com/tliron/puccini/tosca/profiles/implicit/v2_0"
	"github.com/tliron/yamlkeys"
)

//
// Built-in functions and constraints
//
// [TOSCA-v2.0] @ ?
// [TOSCA-Simple-Profile-YAML-v1.3] @ 4
// [TOSCA-Simple-Profile-YAML-v1.2] @ 4
// [TOSCA-Simple-Profile-YAML-v1.1] @ 4
// [TOSCA-Simple-Profile-YAML-v1.0] @ 4
//
// [TOSCA-v2.0] @ ?
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.6.3
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.3
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.2
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.5.2
//

const functionPathPrefix = "/tosca/implicit/2.0/js/functions/"

var FunctionScriptlets = map[string]string{
	tosca.METADATA_FUNCTION_PREFIX + "concat":               profile.Profile[functionPathPrefix+"concat.js"],
	tosca.METADATA_FUNCTION_PREFIX + "join":                 profile.Profile[functionPathPrefix+"join.js"], // introduced in TOSCA 1.2
	tosca.METADATA_FUNCTION_PREFIX + "token":                profile.Profile[functionPathPrefix+"token.js"],
	tosca.METADATA_FUNCTION_PREFIX + "get_input":            profile.Profile[functionPathPrefix+"get_input.js"],
	tosca.METADATA_FUNCTION_PREFIX + "get_property":         profile.Profile[functionPathPrefix+"get_property.js"],
	tosca.METADATA_FUNCTION_PREFIX + "get_attribute":        profile.Profile[functionPathPrefix+"get_attribute.js"],
	tosca.METADATA_FUNCTION_PREFIX + "get_operation_output": profile.Profile[functionPathPrefix+"get_operation_output.js"],
	tosca.METADATA_FUNCTION_PREFIX + "get_nodes_of_type":    profile.Profile[functionPathPrefix+"get_nodes_of_type.js"],
	tosca.METADATA_FUNCTION_PREFIX + "get_artifact":         profile.Profile[functionPathPrefix+"get_artifact.js"],
	tosca.METADATA_FUNCTION_PREFIX + "_get_target_name":     profile.Profile[functionPathPrefix+"_get_target_name.js"],
}

func ParseFunctionCall(context *tosca.Context) bool {
	if _, ok := context.Data.(*tosca.FunctionCall); ok {
		// It's already a function call
		return true
	}

	map_, ok := context.Data.(ard.Map)
	if !ok {
		return false
	}
	count := len(map_)

	changed := false
	if prefixLength := len(context.FunctionPrefix); prefixLength > 0 {
		if count == 0 {
			return false
		}

		for key, value := range map_ {
			key_ := yamlkeys.KeyString(key)

			if strings.HasPrefix(key_, context.FunctionPrefix) {
				scriptletName := key_[prefixLength:]

				// Double prefix means escape
				if strings.HasPrefix(scriptletName, context.FunctionPrefix) {
					delete(map_, key)
					map_[scriptletName] = value
					changed = true
					continue
				}

				if count != 1 {
					context.ReportValueMalformed("function", "more than one entry in map")
					return false
				}

				scriptletName = tosca.METADATA_FUNCTION_PREFIX + scriptletName
				if _, ok := context.ScriptletNamespace.Lookup(scriptletName); !ok {
					// Not a function call, despite having the right data structure
					context.Clone(scriptletName).ReportValueInvalid("function", "unsupported")
					return false
				}

				setFunctionCall(context, scriptletName, value)
				return true
			}
		}
	} else {
		if count != 1 {
			return false
		}

		// Only one iteration
		for key, data := range map_ {
			scriptletName := tosca.METADATA_FUNCTION_PREFIX + yamlkeys.KeyString(key)

			if _, ok := context.ScriptletNamespace.Lookup(scriptletName); !ok {
				// Not a function call, despite having the right data structure
				return false
			}

			setFunctionCall(context, scriptletName, data)
			return true
		}
	}

	return changed
}

func ParseFunctionCalls(context *tosca.Context) bool {
	changed := false
	if ParseFunctionCall(context) {
		changed = true
	} else if list, ok := context.Data.(ard.List); ok {
		for index, value := range list {
			childContext := context.ListChild(index, value)
			if ParseFunctionCalls(childContext) {
				changed = true
				list[index] = childContext.Data
			}
		}
	} else if map_, ok := context.Data.(ard.Map); ok {
		for key, value := range map_ {
			childContext := context.MapChild(key, value)
			if ParseFunctionCalls(childContext) {
				changed = true
				yamlkeys.MapPut(map_, key, childContext.Data) // support complex keys
			}
		}
	}
	return changed
}

func NormalizeFunctionCallArguments(functionCall *tosca.FunctionCall, context *tosca.Context) {
	for index, argument := range functionCall.Arguments {
		// Because the same constraint instance may be shared among more than one value, this
		// func might be called more than once on the same arguments, so we must make sure not
		// to normalize more than once
		if _, ok := argument.(normal.Value); !ok {
			if value, ok := argument.(*Value); ok {
				functionCall.Arguments[index] = value.Normalize()
			} else {
				// Note: this literal value will not have a $meta field
				functionCall.Arguments[index] = NewValue(context.ListChild(index, argument)).Normalize()
			}
		}
	}
}

// Utils

func setFunctionCall(context *tosca.Context, scriptletName string, data ard.Value) {
	// Some functions accept a list of arguments, some just one argument
	originalArguments, ok := data.(ard.List)
	if !ok {
		originalArguments = ard.List{data}
	}

	// Arguments may be function calls, recursively
	arguments := make(ard.List, len(originalArguments))
	for index, argument := range originalArguments {
		argumentContext := context.Clone(argument)
		ParseFunctionCalls(argumentContext)
		arguments[index] = argumentContext.Data
	}

	context.Data = context.NewFunctionCall(scriptletName, arguments)
}
