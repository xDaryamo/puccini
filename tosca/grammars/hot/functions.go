package hot

import (
	"github.com/tliron/kutil/ard"
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
	profile "github.com/tliron/puccini/tosca/profiles/hot/v1_0"
	"github.com/tliron/yamlkeys"
)

//
// Built-in functions
//
// [https://docs.openstack.org/heat/stein/template_guide/hot_spec.html#intrinsic-functions]
//

var FunctionScriptlets = map[string]string{
	"tosca.function.and":                 profile.Profile["/hot/1.0/js/functions/and.js"],
	"tosca.function.contains":            profile.Profile["/hot/1.0/js/functions/contains.js"],
	"tosca.function.digest":              profile.Profile["/hot/1.0/js/functions/digest.js"],
	"tosca.function.equals":              profile.Profile["/hot/1.0/js/functions/equals.js"],
	"tosca.function.filter":              profile.Profile["/hot/1.0/js/functions/filter.js"],
	"tosca.function.get_attr":            profile.Profile["/hot/1.0/js/functions/get_attr.js"],
	"tosca.function.get_file":            profile.Profile["/hot/1.0/js/functions/get_file.js"],
	"tosca.function.get_param":           profile.Profile["/hot/1.0/js/functions/get_param.js"],
	"tosca.function.get_resource":        profile.Profile["/hot/1.0/js/functions/get_resource.js"],
	"tosca.function.if":                  profile.Profile["/hot/1.0/js/functions/if.js"],
	"tosca.function.list_concat_unique":  profile.Profile["/hot/1.0/js/functions/list_concat_unique.js"],
	"tosca.function.list_concat":         profile.Profile["/hot/1.0/js/functions/list_concat.js"],
	"tosca.function.list_join":           profile.Profile["/hot/1.0/js/functions/list_join.js"],
	"tosca.function.make_url":            profile.Profile["/hot/1.0/js/functions/make_url.js"],
	"tosca.function.map_merge":           profile.Profile["/hot/1.0/js/functions/map_merge.js"],
	"tosca.function.map_replace":         profile.Profile["/hot/1.0/js/functions/map_replace.js"],
	"tosca.function.not":                 profile.Profile["/hot/1.0/js/functions/not.js"],
	"tosca.function.or":                  profile.Profile["/hot/1.0/js/functions/or.js"],
	"tosca.function.repeat":              profile.Profile["/hot/1.0/js/functions/repeat.js"],
	"tosca.function.resolve":             profile.Profile["/hot/1.0/js/functions/resolve.js"],
	"tosca.function.resource_facade":     profile.Profile["/hot/1.0/js/functions/resource_facade.js"],
	"tosca.function.str_replace_strict":  profile.Profile["/hot/1.0/js/functions/str_replace_strict.js"],
	"tosca.function.str_replace_vstrict": profile.Profile["/hot/1.0/js/functions/str_replace_vstrict.js"],
	"tosca.function.str_replace":         profile.Profile["/hot/1.0/js/functions/str_replace.js"],
	"tosca.function.str_split":           profile.Profile["/hot/1.0/js/functions/str_split.js"],
	"tosca.function.yaql":                profile.Profile["/hot/1.0/js/functions/yaql.js"],
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

		// The "list_join" function has a nested argument structure that we need to flatten
		// https://docs.openstack.org/heat/stein/template_guide/hot_spec.html#list-join
		if name == "list_join" {
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
				yamlkeys.MapPut(map_, key, childContext.Data) // support complex keys
			}
		}
	}
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
