package cloudify_v1_3

import (
	"github.com/tliron/puccini/ard"
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
	profile "github.com/tliron/puccini/tosca/profiles/cloudify/v4_5"
)

//
// Built-in functions
//
// [https://docs.cloudify.co/4.5.5/developer/blueprints/spec-intrinsic-functions/]
//

var FunctionSourceCode = map[string]string{
	"concat":         profile.Profile["/cloudify/4.5/js/get_secret.js"],
	"get_attribute":  profile.Profile["/cloudify/4.5/js/get_attribute.js"],
	"get_capability": profile.Profile["/cloudify/4.5/js/get_capability.js"],
	"get_input":      profile.Profile["/cloudify/4.5/js/get_input.js"],
	"get_property":   profile.Profile["/cloudify/4.5/js/get_property.js"],
	"get_secret":     profile.Profile["/cloudify/4.5/js/get_secret.js"],
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

func ToFunctionCalls(context *tosca.Context) {
	if !ToFunctionCall(context) {
		if list, ok := context.Data.(ard.List); ok {
			for index, value := range list {
				childContext := context.ListChild(index, value)
				ToFunctionCalls(childContext)
				list[index] = childContext.Data
			}
		} else if map_, ok := context.Data.(ard.Map); ok {
			for key, value := range map_ {
				childContext := context.MapChild(key, value)
				ToFunctionCalls(childContext)
				map_[key] = childContext.Data
			}
		}
	}
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
