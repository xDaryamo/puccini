package parser

import (
	"github.com/tliron/puccini/hot/grammars/v2018_08_31"
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/grammars/v1_1"
	"github.com/tliron/puccini/tosca/grammars/v1_2"
)

type Grammars map[string]tosca.Reader

var UnitGrammars = Grammars{
	"tosca_simple_yaml_1_2":            v1_1.ReadUnit,
	"tosca_simple_yaml_1_1":            v1_1.ReadUnit,
	"tosca_simple_yaml_1_0":            v1_1.ReadUnit, // TODO: hmmm
	"tosca_simple_profile_for_nfv_1_0": v1_1.ReadUnit,
}

var ServiceTemplateGrammars = Grammars{
	"tosca_simple_yaml_1_2":            v1_2.ReadServiceTemplate,
	"tosca_simple_yaml_1_1":            v1_1.ReadServiceTemplate,
	"tosca_simple_yaml_1_0":            v1_1.ReadServiceTemplate, // TODO: hmmm
	"tosca_simple_profile_for_nfv_1_0": v1_1.ReadServiceTemplate,
	"2018-08-31":                       v2018_08_31.ReadTemplate,
}

func GetGrammar(toscaContext *tosca.Context, grammars Grammars) (tosca.Reader, bool) {
	versionContext, ok := toscaContext.GetFieldChild("tosca_definitions_version")
	if !ok {
		versionContext, ok = toscaContext.GetFieldChild("heat_template_version")
		if !ok {
			return nil, false
		}
	}

	version := versionContext.ReadString()
	if version == nil {
		return nil, false
	}

	reader, ok := grammars[*version]
	if !ok {
		versionContext.ReportFieldUnsupportedValue()
		return nil, false
	}

	return reader, true
}
