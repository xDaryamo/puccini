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

		// The "list_join" function has a nested argument structure that we need to flatten
		// https://docs.openstack.org/heat/rocky/template_guide/hot_spec.html#list-join
		if key == "list_join" {
			newArguments := ard.List{originalArguments[0]}
			for _, argument := range originalArguments[1:] {
				if nestedArguments, ok := argument.(ard.List); ok {
					newArguments = append(newArguments, nestedArguments...)
				} else {
					newArguments = append(newArguments, argument)
				}
			}
			originalArguments = newArguments
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

func ToFunctions(context *tosca.Context) interface{} {
	data := context.Data
	if function, ok := GetFunction(context); ok {
		data = function
	} else if list, ok := data.(ard.List); ok {
		for index, value := range list {
			childContext := context.ListChild(index, value)
			list[index] = ToFunctions(childContext)
		}
	} else if map_, ok := data.(ard.Map); ok {
		for key, value := range map_ {
			childContext := context.MapChild(key, value)
			map_[key] = ToFunctions(childContext)
		}
	}
	return data
}

func NormalizeFunctionArguments(function *tosca.Function, context *tosca.Context) {
	for index, argument := range function.Arguments {
		if _, ok := argument.(normal.Constrainable); ok {
			// Because the same constraint instance may be shared among more than one value, this
			// func might be called more than once on the same arguments, so we must make sure not
			// to normalize more than once
			return
		}
		value := NewValue(context.ListChild(index, argument))
		function.Arguments[index] = value.Normalize()
	}
}
