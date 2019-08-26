package hot

import (
	"github.com/tliron/puccini/ard"
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
	profile "github.com/tliron/puccini/tosca/profiles/hot/v1_0"
)

//
// Built-in functions
//
// [https://docs.openstack.org/heat/stein/template_guide/hot_spec.html#intrinsic-functions]
//

var FunctionSourceCode = map[string]string{
	"and":                 profile.Profile["/hot/1.0/js/and.js"],
	"contains":            profile.Profile["/hot/1.0/js/contains.js"],
	"digest":              profile.Profile["/hot/1.0/js/digest.js"],
	"equals":              profile.Profile["/hot/1.0/js/equals.js"],
	"filter":              profile.Profile["/hot/1.0/js/filter.js"],
	"get_attr":            profile.Profile["/hot/1.0/js/get_attr.js"],
	"get_file":            profile.Profile["/hot/1.0/js/get_file.js"],
	"get_param":           profile.Profile["/hot/1.0/js/get_param.js"],
	"get_resource":        profile.Profile["/hot/1.0/js/get_resource.js"],
	"if":                  profile.Profile["/hot/1.0/js/if.js"],
	"list_concat_unique":  profile.Profile["/hot/1.0/js/list_concat_unique.js"],
	"list_concat":         profile.Profile["/hot/1.0/js/list_concat.js"],
	"list_join":           profile.Profile["/hot/1.0/js/list_join.js"],
	"make_url":            profile.Profile["/hot/1.0/js/make_url.js"],
	"map_merge":           profile.Profile["/hot/1.0/js/map_merge.js"],
	"map_replace":         profile.Profile["/hot/1.0/js/map_replace.js"],
	"not":                 profile.Profile["/hot/1.0/js/not.js"],
	"or":                  profile.Profile["/hot/1.0/js/or.js"],
	"repeat":              profile.Profile["/hot/1.0/js/repeat.js"],
	"resolve":             profile.Profile["/hot/1.0/js/resolve.js"],
	"resource_facade":     profile.Profile["/hot/1.0/js/resource_facade.js"],
	"str_replace_strict":  profile.Profile["/hot/1.0/js/str_replace_strict.js"],
	"str_replace_vstrict": profile.Profile["/hot/1.0/js/str_replace_vstrict.js"],
	"str_replace":         profile.Profile["/hot/1.0/js/str_replace.js"],
	"str_split":           profile.Profile["/hot/1.0/js/str_split.js"],
	"yaql":                profile.Profile["/hot/1.0/js/yaql.js"],
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

		// The "list_join" function has a nested argument structure that we need to flatten
		// https://docs.openstack.org/heat/stein/template_guide/hot_spec.html#list-join
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
