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
// [https://docs.openstack.org/heat/wallaby/template_guide/hot_spec.html#intrinsic-functions]
//

const functionPathPrefix = "/hot/1.0/js/functions/"

var FunctionScriptlets = map[string]string{
	tosca.FunctionScriptletPrefix + "and":                 profile.Profile[functionPathPrefix+"and.js"],
	tosca.FunctionScriptletPrefix + "contains":            profile.Profile[functionPathPrefix+"contains.js"],
	tosca.FunctionScriptletPrefix + "digest":              profile.Profile[functionPathPrefix+"digest.js"],
	tosca.FunctionScriptletPrefix + "equals":              profile.Profile[functionPathPrefix+"equals.js"],
	tosca.FunctionScriptletPrefix + "filter":              profile.Profile[functionPathPrefix+"filter.js"],
	tosca.FunctionScriptletPrefix + "get_attr":            profile.Profile[functionPathPrefix+"get_attr.js"],
	tosca.FunctionScriptletPrefix + "get_file":            profile.Profile[functionPathPrefix+"get_file.js"],
	tosca.FunctionScriptletPrefix + "get_param":           profile.Profile[functionPathPrefix+"get_param.js"],
	tosca.FunctionScriptletPrefix + "get_resource":        profile.Profile[functionPathPrefix+"get_resource.js"],
	tosca.FunctionScriptletPrefix + "if":                  profile.Profile[functionPathPrefix+"if.js"],
	tosca.FunctionScriptletPrefix + "list_concat_unique":  profile.Profile[functionPathPrefix+"list_concat_unique.js"],
	tosca.FunctionScriptletPrefix + "list_concat":         profile.Profile[functionPathPrefix+"list_concat.js"],
	tosca.FunctionScriptletPrefix + "list_join":           profile.Profile[functionPathPrefix+"list_join.js"],
	tosca.FunctionScriptletPrefix + "make_url":            profile.Profile[functionPathPrefix+"make_url.js"],
	tosca.FunctionScriptletPrefix + "map_merge":           profile.Profile[functionPathPrefix+"map_merge.js"],
	tosca.FunctionScriptletPrefix + "map_replace":         profile.Profile[functionPathPrefix+"map_replace.js"],
	tosca.FunctionScriptletPrefix + "not":                 profile.Profile[functionPathPrefix+"not.js"],
	tosca.FunctionScriptletPrefix + "or":                  profile.Profile[functionPathPrefix+"or.js"],
	tosca.FunctionScriptletPrefix + "repeat":              profile.Profile[functionPathPrefix+"repeat.js"],
	tosca.FunctionScriptletPrefix + "resolve":             profile.Profile[functionPathPrefix+"resolve.js"],
	tosca.FunctionScriptletPrefix + "resource_facade":     profile.Profile[functionPathPrefix+"resource_facade.js"],
	tosca.FunctionScriptletPrefix + "str_replace_strict":  profile.Profile[functionPathPrefix+"str_replace_strict.js"],
	tosca.FunctionScriptletPrefix + "str_replace_vstrict": profile.Profile[functionPathPrefix+"str_replace_vstrict.js"],
	tosca.FunctionScriptletPrefix + "str_replace":         profile.Profile[functionPathPrefix+"str_replace.js"],
	tosca.FunctionScriptletPrefix + "str_split":           profile.Profile[functionPathPrefix+"str_split.js"],
	tosca.FunctionScriptletPrefix + "yaql":                profile.Profile[functionPathPrefix+"yaql.js"],
}

func ParseFunctionCalls(context *tosca.Context) bool {
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

		scriptletName := tosca.FunctionScriptletPrefix + name
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
		// https://docs.openstack.org/heat/wallaby/template_guide/hot_spec.html#list-join
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
			ParseFunctionCalls(argumentContext)
			arguments[index] = argumentContext.Data
		}

		context.Data = context.NewFunctionCall(scriptletName, arguments)

		// We have only one key
		return true
	}

	return false
}

func ToFunctionCalls(context *tosca.Context) {
	if !ParseFunctionCalls(context) {
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
