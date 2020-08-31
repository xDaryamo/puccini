package tosca_v2_0

import (
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

var FunctionScriptlets = map[string]string{
	"tosca.function.concat":               profile.Profile["/tosca/implicit/2.0/js/functions/concat.js"],
	"tosca.function.join":                 profile.Profile["/tosca/implicit/2.0/js/functions/join.js"], // introduced in TOSCA 1.2
	"tosca.function.token":                profile.Profile["/tosca/implicit/2.0/js/functions/token.js"],
	"tosca.function.get_input":            profile.Profile["/tosca/implicit/2.0/js/functions/get_input.js"],
	"tosca.function.get_property":         profile.Profile["/tosca/implicit/2.0/js/functions/get_property.js"],
	"tosca.function.get_attribute":        profile.Profile["/tosca/implicit/2.0/js/functions/get_attribute.js"],
	"tosca.function.get_operation_output": profile.Profile["/tosca/implicit/2.0/js/functions/get_operation_output.js"],
	"tosca.function.get_nodes_of_type":    profile.Profile["/tosca/implicit/2.0/js/functions/get_nodes_of_type.js"],
	"tosca.function.get_artifact":         profile.Profile["/tosca/implicit/2.0/js/functions/get_artifact.js"],
}

func ToFunctionCall(context *tosca.Context) bool {
	if _, ok := context.Data.(*tosca.FunctionCall); ok {
		// It's already a function call
		return true
	}

	map_, ok := context.Data.(ard.Map)
	if !ok || len(map_) != 1 {
		return false
	}

	for key, data := range map_ {
		name := yamlkeys.KeyString(key)

		scriptletName := "tosca.function." + name
		_, ok := context.ScriptletNamespace.Lookup(scriptletName)
		if !ok {
			// Not a function call, despite having the right data structure
			return false
		}

		// Some functions accept a list of arguments, some just one argument
		originalArguments, ok := data.(ard.List)
		if !ok {
			originalArguments = ard.List{data}
		}

		// Arguments may be function calls
		arguments := make(ard.List, len(originalArguments))
		for index, argument := range originalArguments {
			argumentContext := context.Clone(argument)
			ToFunctionCall(argumentContext)
			arguments[index] = argumentContext.Data
		}

		context.Data = context.NewFunctionCall(scriptletName, arguments)

		// We have only one key
		return true
	}

	return false
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
