package hot

import (
	"github.com/tliron/go-ard"
	"github.com/tliron/puccini/assets/tosca/profiles"
	"github.com/tliron/puccini/normal"
	"github.com/tliron/puccini/tosca/parsing"
	"github.com/tliron/yamlkeys"
)

//
// Built-in functions
//
// [https://docs.openstack.org/heat/wallaby/template_guide/hot_spec.html#intrinsic-functions]
//

const functionPathPrefix = "hot/1.0/js/functions/"

var FunctionScriptlets = map[string]string{
	parsing.MetadataFunctionPrefix + "and":                 profiles.GetString(functionPathPrefix + "and.js"),
	parsing.MetadataFunctionPrefix + "contains":            profiles.GetString(functionPathPrefix + "contains.js"),
	parsing.MetadataFunctionPrefix + "digest":              profiles.GetString(functionPathPrefix + "digest.js"),
	parsing.MetadataFunctionPrefix + "equals":              profiles.GetString(functionPathPrefix + "equals.js"),
	parsing.MetadataFunctionPrefix + "filter":              profiles.GetString(functionPathPrefix + "filter.js"),
	parsing.MetadataFunctionPrefix + "get_attr":            profiles.GetString(functionPathPrefix + "get_attr.js"),
	parsing.MetadataFunctionPrefix + "get_file":            profiles.GetString(functionPathPrefix + "get_file.js"),
	parsing.MetadataFunctionPrefix + "get_param":           profiles.GetString(functionPathPrefix + "get_param.js"),
	parsing.MetadataFunctionPrefix + "get_resource":        profiles.GetString(functionPathPrefix + "get_resource.js"),
	parsing.MetadataFunctionPrefix + "if":                  profiles.GetString(functionPathPrefix + "if.js"),
	parsing.MetadataFunctionPrefix + "list_concat_unique":  profiles.GetString(functionPathPrefix + "list_concat_unique.js"),
	parsing.MetadataFunctionPrefix + "list_concat":         profiles.GetString(functionPathPrefix + "list_concat.js"),
	parsing.MetadataFunctionPrefix + "list_join":           profiles.GetString(functionPathPrefix + "list_join.js"),
	parsing.MetadataFunctionPrefix + "make_url":            profiles.GetString(functionPathPrefix + "make_url.js"),
	parsing.MetadataFunctionPrefix + "map_merge":           profiles.GetString(functionPathPrefix + "map_merge.js"),
	parsing.MetadataFunctionPrefix + "map_replace":         profiles.GetString(functionPathPrefix + "map_replace.js"),
	parsing.MetadataFunctionPrefix + "not":                 profiles.GetString(functionPathPrefix + "not.js"),
	parsing.MetadataFunctionPrefix + "or":                  profiles.GetString(functionPathPrefix + "or.js"),
	parsing.MetadataFunctionPrefix + "repeat":              profiles.GetString(functionPathPrefix + "repeat.js"),
	parsing.MetadataFunctionPrefix + "resource_facade":     profiles.GetString(functionPathPrefix + "resource_facade.js"),
	parsing.MetadataFunctionPrefix + "str_replace_strict":  profiles.GetString(functionPathPrefix + "str_replace_strict.js"),
	parsing.MetadataFunctionPrefix + "str_replace_vstrict": profiles.GetString(functionPathPrefix + "str_replace_vstrict.js"),
	parsing.MetadataFunctionPrefix + "str_replace":         profiles.GetString(functionPathPrefix + "str_replace.js"),
	parsing.MetadataFunctionPrefix + "str_split":           profiles.GetString(functionPathPrefix + "str_split.js"),
	parsing.MetadataFunctionPrefix + "yaql":                profiles.GetString(functionPathPrefix + "yaql.js"),
}

func ParseFunctionCall(context *parsing.Context) bool {
	if _, ok := context.Data.(*parsing.FunctionCall); ok {
		// It's already a function call
		return true
	}

	map_, ok := context.Data.(ard.Map)
	if !ok || len(map_) != 1 {
		return false
	}

	for key, data := range map_ {
		name := yamlkeys.KeyString(key)

		scriptletName := parsing.MetadataFunctionPrefix + name
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

func ParseFunctionCalls(context *parsing.Context) bool {
	changed := false
	if ParseFunctionCall(context) {
		changed = true
	} else if list, ok := context.Data.(ard.List); ok {
		for index, value := range list {
			childContext := context.ListChild(index, value)
			if ParseFunctionCalls(childContext) {
				changed = true
			}
			list[index] = childContext.Data
		}
	} else if map_, ok := context.Data.(ard.Map); ok {
		for key, value := range map_ {
			childContext := context.MapChild(key, value)
			if ParseFunctionCalls(childContext) {
				changed = true
			}
			yamlkeys.MapPut(map_, key, childContext.Data) // support complex keys
		}
	}
	return changed
}

func NormalizeFunctionCallArguments(functionCall *parsing.FunctionCall, context *parsing.Context) {
	for index, argument := range functionCall.Arguments {
		// Because the same constraint instance may be shared among more than one value, this
		// func might be called more than once on the same arguments, so we must make sure not
		// to normalize more than once
		if _, ok := argument.(normal.Value); !ok {
			if value, ok := argument.(*Value); ok {
				functionCall.Arguments[index] = value.Normalize()
			} else {
				// Note: this literal value will not have a $type field
				functionCall.Arguments[index] = NewValue(context.ListChild(index, argument)).Normalize()
			}
		}
	}
}
