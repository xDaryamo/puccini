package parser

import (
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/grammars/v1_1"
)

type Grammars map[string]tosca.Reader

var UnitGrammars = Grammars{
	"tosca_simple_yaml_1_1":            v1_1.ReadUnit,
	"tosca_simple_yaml_1_0":            v1_1.ReadUnit, // TODO: hmmm
	"tosca_simple_profile_for_nfv_1_0": v1_1.ReadUnit,
}

var ServiceTemplateGrammars = Grammars{
	"tosca_simple_yaml_1_1":            v1_1.ReadServiceTemplate,
	"tosca_simple_yaml_1_0":            v1_1.ReadServiceTemplate, // TODO: hmmm
	"tosca_simple_profile_for_nfv_1_0": v1_1.ReadServiceTemplate,
}

func GetGrammar(toscaContext *tosca.Context, grammars Grammars) (tosca.Reader, bool) {
	toscaContext, ok := toscaContext.RequiredFieldChild("tosca_definitions_version")
	if !ok {
		return nil, false
	}

	toscaDefinitionsVersion := toscaContext.ReadString()
	if toscaDefinitionsVersion == nil {
		return nil, false
	}

	reader, ok := grammars[*toscaDefinitionsVersion]
	if !ok {
		toscaContext.ReportFieldUnsupportedValue()
		return nil, false
	}

	return reader, true
}
