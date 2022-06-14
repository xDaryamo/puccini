package tosca_v1_1

import (
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/grammars/tosca_v2_0"
)

//
// ServiceFile
//
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.9
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.9
//

// tosca.Reader signature
func ReadServiceFile(context *tosca.Context) tosca.EntityPtr {
	context.SetReadTag("ServiceTemplate", "topology_template,ServiceTemplate")
	context.SetReadTag("Profile", "")

	self := tosca_v2_0.NewServiceFile(context)
	context.ScriptletNamespace.Merge(DefaultScriptletNamespace)
	ignore := []string{"dsl_definitions"}
	if context.HasQuirk(tosca.QuirkAnnotationsIgnore) {
		ignore = append(ignore, "annotation_types")
	}
	context.ValidateUnsupportedFields(append(context.ReadFields(self), ignore...))
	return self
}
