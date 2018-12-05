package v1_1

import (
	"github.com/tliron/puccini/ard"
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
	profile "github.com/tliron/puccini/tosca/profiles/simple/v1_1"
)

//
// Built-in functions
//
// [TOSCA-Simple-Profile-YAML-v1.1] @ 4
//

var FunctionSourceCode = map[string]string{
	"concat":               profile.Profile["/tosca/simple/1.1/js/concat.js"],
	"token":                profile.Profile["/tosca/simple/1.1/js/token.js"],
	"get_input":            profile.Profile["/tosca/simple/1.1/js/get_input.js"],
	"get_property":         profile.Profile["/tosca/simple/1.1/js/get_property.js"],
	"get_attribute":        profile.Profile["/tosca/simple/1.1/js/get_attribute.js"],
	"get_operation_output": profile.Profile["/tosca/simple/1.1/js/get_operation_output.js"],
	"get_nodes_of_type":    profile.Profile["/tosca/simple/1.1/js/get_nodes_of_type.js"],
	"get_artifact":         profile.Profile["/tosca/simple/1.1/js/get_artifact.js"],
}

func GetFunction(context *tosca.Context) (*tosca.Function, bool) {
	if _, ok := context.Data.(*tosca.Function); ok {
		// It's already a function
		return nil, false
	}

	map_, ok := context.Data.(ard.Map)
	if !ok || len(map_) != 1 {
		return nil, false
	}

	for key, data := range map_ {
		_, ok := context.ScriptNamespace[key]
		if !ok {
			// Not a function, despite having the right data structure
			return nil, false
		}

		// Some functions accept a list of arguments, some just one argument
		originalArguments, ok := data.(ard.List)
		if !ok {
			originalArguments = ard.List{data}
		}

		// Arguments may be functions
		arguments := make(ard.List, len(originalArguments))
		for index, argument := range originalArguments {
			if f, ok := GetFunction(context.WithData(argument)); ok {
				argument = f
			}
			arguments[index] = argument
		}

		// We have only one key
		return tosca.NewFunction(context.Path, key, arguments), true
	}

	return nil, false
}

func NormalizeFunctionArguments(function *tosca.Function, context *tosca.Context) {
	for index, argument := range function.Arguments {
		if _, ok := argument.(normal.Constrainable); ok {
			// Because the same constraint clause may be shared among many values, this func
			// might be called more than once on the same arguments, so we must make sure not
			// to normalize more than once
			return
		}
		value := NewValue(context.ListChild(index, argument))
		function.Arguments[index] = value.Normalize()
	}
}
