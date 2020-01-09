package tosca_v1_2

import (
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/grammars/tosca_v1_3"
)

//
// Unit
//
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.10
//

// tosca.Reader signature
func ReadUnit(context *tosca.Context) interface{} {
	self := tosca_v1_3.NewUnit(context)
	context.ScriptletNamespace.Merge(DefaultScriptletNamespace)
	context.ValidateUnsupportedFields(append(context.ReadFields(self), "dsl_definitions"))
	return self
}
