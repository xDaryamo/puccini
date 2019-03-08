package parser

import (
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/grammars/hot"
	"github.com/tliron/puccini/tosca/grammars/tosca_v1_1"
	"github.com/tliron/puccini/tosca/grammars/tosca_v1_2"
)

var Grammars = map[string]tosca.Grammar{
	"tosca_simple_yaml_1_2":            tosca_v1_2.Grammar,
	"tosca_simple_yaml_1_1":            tosca_v1_1.Grammar,
	"tosca_simple_yaml_1_0":            tosca_v1_1.Grammar, // TODO: properly support 1.0
	"tosca_simple_profile_for_nfv_1_0": tosca_v1_2.Grammar,
	"2018-08-31":                       hot.Grammar, // rocky
	"2018-03-02":                       hot.Grammar, // queens
	"2017-09-01":                       hot.Grammar, // pike
	"2017-02-24":                       hot.Grammar, // ocata
	"2016-10-14":                       hot.Grammar, // newton
	"2016-04-08":                       hot.Grammar, // mitaka
	"2015-10-15":                       hot.Grammar, // liberty
	"2015-04-30":                       hot.Grammar, // kilo
	"2014-10-16":                       hot.Grammar, // juno
	"2013-05-23":                       hot.Grammar, // icehouse
}

func DetectGrammar(context *tosca.Context) bool {
	var versionContext *tosca.Context
	var ok bool
	if versionContext, ok = context.GetFieldChild("tosca_definitions_version"); !ok {
		if versionContext, ok = context.GetFieldChild("heat_template_version"); !ok {
			return false
		}
	}

	if version := versionContext.ReadString(); version != nil {
		if context.Grammar, ok = Grammars[*version]; ok {
			return true
		} else {
			versionContext.ReportFieldUnsupportedValue()
		}
	}

	return false
}
