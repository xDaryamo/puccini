package parser

import (
	"github.com/tliron/puccini/hot/grammars/v2018_08_31"
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/grammars/v1_1"
	"github.com/tliron/puccini/tosca/grammars/v1_2"
)

type Grammars map[string]tosca.Reader

var ImportedEntityGrammars = Grammars{
	"tosca_simple_yaml_1_2":            v1_1.ReadUnit,
	"tosca_simple_yaml_1_1":            v1_1.ReadUnit,
	"tosca_simple_yaml_1_0":            v1_1.ReadUnit, // TODO: properly support 1.0
	"tosca_simple_profile_for_nfv_1_0": v1_1.ReadUnit,
}

var RootEntityGrammars = Grammars{
	"tosca_simple_yaml_1_2":            v1_2.ReadServiceTemplate,
	"tosca_simple_yaml_1_1":            v1_1.ReadServiceTemplate,
	"tosca_simple_yaml_1_0":            v1_1.ReadServiceTemplate, // TODO: properly support 1.0
	"tosca_simple_profile_for_nfv_1_0": v1_1.ReadServiceTemplate,
	"2018-08-31":                       v2018_08_31.ReadTemplate, // rocky
	"2018-03-02":                       v2018_08_31.ReadTemplate, // queens
	"2017-09-01":                       v2018_08_31.ReadTemplate, // pike
	"2017-02-24":                       v2018_08_31.ReadTemplate, // ocata
	"2016-10-14":                       v2018_08_31.ReadTemplate, // newton
	"2016-04-08":                       v2018_08_31.ReadTemplate, // mitaka
	"2015-10-15":                       v2018_08_31.ReadTemplate, // liberty
	"2015-04-30":                       v2018_08_31.ReadTemplate, // kilo
	"2014-10-16":                       v2018_08_31.ReadTemplate, // juno
	"2013-05-23":                       v2018_08_31.ReadTemplate, // icehouse
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
