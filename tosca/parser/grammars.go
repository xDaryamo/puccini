package parser

import (
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/grammars/hot"
	"github.com/tliron/puccini/tosca/grammars/tosca_v1_1"
	"github.com/tliron/puccini/tosca/grammars/tosca_v1_2"
)

type Grammars map[string]tosca.Reader

var ImportedEntityGrammars = Grammars{
	"tosca_simple_yaml_1_2":            tosca_v1_1.ReadUnit,
	"tosca_simple_yaml_1_1":            tosca_v1_1.ReadUnit,
	"tosca_simple_yaml_1_0":            tosca_v1_1.ReadUnit, // TODO: properly support 1.0
	"tosca_simple_profile_for_nfv_1_0": tosca_v1_1.ReadUnit,
}

var RootEntityGrammars = Grammars{
	"tosca_simple_yaml_1_2":            tosca_v1_2.ReadServiceTemplate,
	"tosca_simple_yaml_1_1":            tosca_v1_1.ReadServiceTemplate,
	"tosca_simple_yaml_1_0":            tosca_v1_1.ReadServiceTemplate, // TODO: properly support 1.0
	"tosca_simple_profile_for_nfv_1_0": tosca_v1_1.ReadServiceTemplate,
	"2018-08-31":                       hot.ReadTemplate, // rocky
	"2018-03-02":                       hot.ReadTemplate, // queens
	"2017-09-01":                       hot.ReadTemplate, // pike
	"2017-02-24":                       hot.ReadTemplate, // ocata
	"2016-10-14":                       hot.ReadTemplate, // newton
	"2016-04-08":                       hot.ReadTemplate, // mitaka
	"2015-10-15":                       hot.ReadTemplate, // liberty
	"2015-04-30":                       hot.ReadTemplate, // kilo
	"2014-10-16":                       hot.ReadTemplate, // juno
	"2013-05-23":                       hot.ReadTemplate, // icehouse
}

func GetGrammar(toscaContext *tosca.Context, grammars Grammars) (tosca.Reader, bool) {
	if versionContext, ok := toscaContext.GetFieldChild("tosca_definitions_version"); ok {
		if version := versionContext.ReadString(); version != nil {
			if reader, ok := grammars[*version]; ok {
				return reader, true
			}
		}

		versionContext.ReportFieldUnsupportedValue()
	} else if versionContext, ok := toscaContext.GetFieldChild("heat_template_version"); ok {
		if version := versionContext.ReadString(); version != nil {
			if reader, ok := grammars[*version]; ok {
				return reader, true
			}
		}

		versionContext.ReportFieldUnsupportedValue()
	}

	return nil, false
}
