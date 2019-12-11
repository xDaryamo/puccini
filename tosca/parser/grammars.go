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
)

var Grammars = make(map[string]map[string]*tosca.Grammar)

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
			if _, ok := grammars[version.Version]; ok {
				panic(fmt.Sprintf("grammar version conflict: %s = %s", keyword, version.Version))
			}
			grammars[version.Version] = grammar

			if version.ProfileInternalPath != "" {
				var paths map[string]string
				if paths, ok = ProfileInternalPaths[keyword]; !ok {
					paths = make(map[string]string)
					ProfileInternalPaths[keyword] = paths
				}
				paths[version.Version] = version.ProfileInternalPath
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
