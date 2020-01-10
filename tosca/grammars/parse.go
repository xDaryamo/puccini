package grammars

import (
	"time"

	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/url"
)

func Detect(context *tosca.Context) bool {
	if context.Grammar == nil {
		if context.Grammar = GetGrammar(context); context.Grammar == nil {
			context.ReportFieldUnsupportedValue()
		}
	}
	return context.Grammar != nil
}

func GetGrammar(context *tosca.Context) *tosca.Grammar {
	if versionContext, version := DetectVersion(context); version != nil {
		if grammars, ok := Grammars[versionContext.Name]; ok {
			if grammar, ok := grammars[*version]; ok {
				return grammar
			}
		}
	}
	return nil
}

func CompatibleGrammars(context1 *tosca.Context, context2 *tosca.Context) bool {
	return GetGrammar(context1) == GetGrammar(context2)
}

func DetectVersion(context *tosca.Context) (*tosca.Context, *string) {
	var versionContext *tosca.Context
	var ok bool

	for keyword := range Grammars {
		if versionContext, ok = context.GetFieldChild(keyword); ok {
			if keyword == "heat_template_version" {
				// Hack to allow HOT to use YAML !!timestamp values

				if versionContext.Is("string") {
					return versionContext, versionContext.ReadString()
				}

				switch data := versionContext.Data.(type) {
				case time.Time:
					versionContext.Data = data.Format("2006-01-02")
					return versionContext, versionContext.ReadString()
				}

				versionContext.ReportValueWrongType("string", "timestamp")
			} else {
				if versionContext.ValidateType("string") {
					return versionContext, versionContext.ReadString()
				}
			}
		}
	}

	return nil, nil
}

func GetImplicitImportSpec(context *tosca.Context) (*tosca.ImportSpec, bool) {
	if versionContext, version := DetectVersion(context); version != nil {
		if paths, ok := ImplicitProfilePaths[versionContext.Name]; ok {
			if path, ok := paths[*version]; ok {
				if url_, err := url.NewValidInternalURL(path); err == nil {
					return &tosca.ImportSpec{url_, nil, true}, true
				} else {
					context.ReportError(err)
				}
			}
		}
	}

	return nil, false
}
