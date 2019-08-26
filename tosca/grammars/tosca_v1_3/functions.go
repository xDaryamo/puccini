package tosca_v1_3

import (
	"github.com/tliron/puccini/ard"
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
	profile "github.com/tliron/puccini/tosca/profiles/simple/v1_2"
)

//
// Built-in functions
//
// [TOSCA-Simple-Profile-YAML-v1.2] @ 4
// [TOSCA-Simple-Profile-YAML-v1.1] @ 4
//

var FunctionSourceCode = map[string]string{
	"concat":               profile.Profile["/tosca/simple/1.2/js/concat.js"],
	"join":                 profile.Profile["/tosca/simple/1.2/js/join.js"], // introduced in 1.2
	"token":                profile.Profile["/tosca/simple/1.2/js/token.js"],
	"get_input":            profile.Profile["/tosca/simple/1.2/js/get_input.js"],
	"get_property":         profile.Profile["/tosca/simple/1.2/js/get_property.js"],
	"get_attribute":        profile.Profile["/tosca/simple/1.2/js/get_attribute.js"],
	"get_operation_output": profile.Profile["/tosca/simple/1.2/js/get_operation_output.js"],
	"get_nodes_of_type":    profile.Profile["/tosca/simple/1.2/js/get_nodes_of_type.js"],
	"get_artifact":         profile.Profile["/tosca/simple/1.2/js/get_artifact.js"],
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
		_, ok := context.ScriptNamespace[key]
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
			argumentContext := context.WithData(argument)
			ToFunctionCall(argumentContext)
			arguments[index] = argumentContext.Data
		}

		context.Data = context.NewFunctionCall(key, arguments)

		// We have only one key
		return true
	}

	return false
}

func NormalizeFunctionCallArguments(functionCall *tosca.FunctionCall, context *tosca.Context) {
	for index, argument := range functionCall.Arguments {
		if _, ok := argument.(normal.Constrainable); ok {
			// Because the same constraint instance may be shared among more than one value, this
			// func might be called more than once on the same arguments, so we must make sure not
			// to normalize more than once
			return
		}
		value := NewValue(context.ListChild(index, argument))
		functionCall.Arguments[index] = value.Normalize()
	}
}
