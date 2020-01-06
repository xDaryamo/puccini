package parser

import (
	"fmt"

	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/grammars/cloudify_v1_3"
	"github.com/tliron/puccini/tosca/grammars/hot"
	"github.com/tliron/puccini/tosca/grammars/tosca_v1_0"
	"github.com/tliron/puccini/tosca/grammars/tosca_v1_1"
	"github.com/tliron/puccini/tosca/grammars/tosca_v1_2"
	"github.com/tliron/puccini/tosca/grammars/tosca_v1_3"
	"github.com/tliron/puccini/url"
)

// Map of keyword -> version -> grammar
var Grammars = make(map[string]map[string]*tosca.Grammar)

// Map of keyword -> version -> internal URL path
var ImplicitProfilePaths = make(map[string]map[string]string)

func init() {
	initGrammar(&tosca_v1_3.Grammar)
	initGrammar(&tosca_v1_2.Grammar)
	initGrammar(&tosca_v1_1.Grammar)
	initGrammar(&tosca_v1_0.Grammar)
	initGrammar(&cloudify_v1_3.Grammar)
	initGrammar(&hot.Grammar)
}

func initGrammar(grammar *tosca.Grammar) {
	for keyword, versions := range grammar.Versions {
		var grammars map[string]*tosca.Grammar
		var ok bool

		if grammars, ok = Grammars[keyword]; !ok {
			grammars = make(map[string]*tosca.Grammar)
			Grammars[keyword] = grammars
		}

		for _, version := range versions {
			if _, ok := grammars[version.Version]; !ok {
				grammars[version.Version] = grammar
			} else {
				panic(fmt.Sprintf("initGrammar version conflict: %s = %s", keyword, version.Version))
			}

			if version.ImplicitProfilePath != "" {
				var paths map[string]string
				if paths, ok = ImplicitProfilePaths[keyword]; !ok {
					paths = make(map[string]string)
					ImplicitProfilePaths[keyword] = paths
				}

				if _, ok := paths[version.Version]; !ok {
					paths[version.Version] = version.ImplicitProfilePath
				} else {
					panic(fmt.Sprintf("initGrammar implicit profile path conflict: %s = %s -> %s", keyword, version.Version, version.ImplicitProfilePath))
				}
			}
		}
	}
}

func DetectGrammar(context *tosca.Context) bool {
	if versionContext, version := GetVersion(context); version != nil {
		if grammars, ok := Grammars[versionContext.Name]; ok {
			if context.Grammar, ok = grammars[*version]; ok {
				return true
			}
		}
		versionContext.ReportFieldUnsupportedValue()
	}

	return false
}

func GetImplicitImportSpec(context *tosca.Context) (*tosca.ImportSpec, bool) {
	if versionContext, version := GetVersion(context); version != nil {
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
