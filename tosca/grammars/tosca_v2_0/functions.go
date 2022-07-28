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
	tosca.FunctionScriptletPrefix + "concat":               profile.Profile[functionPathPrefix+"concat.js"],
	tosca.FunctionScriptletPrefix + "join":                 profile.Profile[functionPathPrefix+"join.js"], // introduced in TOSCA 1.2
	tosca.FunctionScriptletPrefix + "token":                profile.Profile[functionPathPrefix+"token.js"],
	tosca.FunctionScriptletPrefix + "get_input":            profile.Profile[functionPathPrefix+"get_input.js"],
	tosca.FunctionScriptletPrefix + "get_property":         profile.Profile[functionPathPrefix+"get_property.js"],
	tosca.FunctionScriptletPrefix + "get_attribute":        profile.Profile[functionPathPrefix+"get_attribute.js"],
	tosca.FunctionScriptletPrefix + "get_operation_output": profile.Profile[functionPathPrefix+"get_operation_output.js"],
	tosca.FunctionScriptletPrefix + "get_nodes_of_type":    profile.Profile[functionPathPrefix+"get_nodes_of_type.js"],
	tosca.FunctionScriptletPrefix + "get_artifact":         profile.Profile[functionPathPrefix+"get_artifact.js"],
}

func ParseFunctionCalls(context *tosca.Context) bool {
	// TODO: traverse ARD

	if _, ok := context.Data.(*tosca.FunctionCall); ok {
		// It's already a function call
		return true
	}

	map_, ok := context.Data.(ard.Map)
	if !ok {
		return false
	}
	count := len(map_)

	if prefixLength := len(context.FunctionPrefix); prefixLength > 0 {
		if count == 0 {
			return false
		}

		for key, data := range map_ {
			key_ := yamlkeys.KeyString(key)

			if strings.HasPrefix(key_, context.FunctionPrefix) {
				scriptletName := key_[prefixLength:]

				if strings.HasPrefix(scriptletName, context.FunctionPrefix) {
					// Double prefix means escape
					delete(map_, key)
					map_[scriptletName] = data
					continue
				}

				if count != 1 {
					context.ReportValueMalformed("function", "more than one entry in map")
					return false
				}

				scriptletName = tosca.FunctionScriptletPrefix + scriptletName
				if _, ok := context.ScriptletNamespace.Lookup(scriptletName); !ok {
					// Not a function call, despite having the right data structure
					context.Clone(scriptletName).ReportValueInvalid("function", "unsupported")
					return false
				}

				setFunctionCall(context, scriptletName, data)
				return true
			}
		}
	} else {
		if count != 1 {
			return false
		}

		// Only one iteration
		for key, data := range map_ {
			scriptletName := tosca.FunctionScriptletPrefix + yamlkeys.KeyString(key)

			if _, ok := context.ScriptletNamespace.Lookup(scriptletName); !ok {
				// Not a function call, despite having the right data structure
				return false
			}

			setFunctionCall(context, scriptletName, data)
			return true
		}
	}

	return false
}

func setFunctionCall(context *tosca.Context, scriptletName string, data ard.Value) {
	// TODO: function calls can be nested *anywhere*, and it shouldn't matter if it's a list or not

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

func NormalizeFunctionCallArguments(functionCall *tosca.FunctionCall, context *tosca.Context) {
	for index, argument := range functionCall.Arguments {
		// Because the same constraint instance may be shared among more than one value, this
		// func might be called more than once on the same arguments, so we must make sure not
		// to normalize more than once
		if _, ok := argument.(normal.Constrainable); !ok {
			if value, ok := argument.(*Value); ok {
				functionCall.Arguments[index] = value.Normalize()
			} else {
				// Note: this literal value will not have a $type field
				functionCall.Arguments[index] = NewValue(context.ListChild(index, argument)).Normalize()
			}
		}
	}
}
