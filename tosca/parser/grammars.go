package parser

import (
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/grammars/cloudify_v1_3"
	"github.com/tliron/puccini/tosca/grammars/hot"
	"github.com/tliron/puccini/tosca/grammars/tosca_v1_0"
	"github.com/tliron/puccini/tosca/grammars/tosca_v1_1"
	"github.com/tliron/puccini/tosca/grammars/tosca_v1_2"
	"github.com/tliron/puccini/tosca/grammars/tosca_v1_3"
)

var Grammars = map[string]map[string]tosca.Grammar{
	"tosca_definitions_version": {
		"tosca_simple_yaml_1_3":            tosca_v1_3.Grammar,
		"tosca_simple_yaml_1_2":            tosca_v1_2.Grammar,
		"tosca_simple_yaml_1_1":            tosca_v1_1.Grammar,
		"tosca_simple_yaml_1_0":            tosca_v1_0.Grammar,
		"tosca_simple_profile_for_nfv_1_0": tosca_v1_3.Grammar,
		"cloudify_dsl_1_3":                 cloudify_v1_3.Grammar,
	},
	"heat_template_version": {
		"stein":      hot.Grammar, // stein (not mentioned in spec, but probably supported)
		"rocky":      hot.Grammar, // rocky
		"queens":     hot.Grammar, // queens
		"pike":       hot.Grammar, // pike
		"newton":     hot.Grammar, // newton
		"ocata":      hot.Grammar, // ocata
		"2018-08-31": hot.Grammar, // stein, rocky
		"2018-03-02": hot.Grammar, // queens
		"2017-09-01": hot.Grammar, // pike
		"2017-02-24": hot.Grammar, // ocata
		"2016-10-14": hot.Grammar, // newton
		"2016-04-08": hot.Grammar, // mitaka
		"2015-10-15": hot.Grammar, // liberty
		"2015-04-30": hot.Grammar, // kilo
		"2014-10-16": hot.Grammar, // juno
		"2013-05-23": hot.Grammar, // icehouse
	},
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
