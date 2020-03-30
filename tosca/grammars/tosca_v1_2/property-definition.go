package tosca_v1_2

import (
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/grammars/tosca_v1_3"
)

//
// PropertyDefinition
//
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.9
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.8
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.5.8
//

// tosca.Reader signature
func ReadPropertyDefinition(context *tosca.Context) interface{} {
	context.SetReadTag("Metadata", "")
	context.SetReadTag("KeySchema", "")

	return tosca_v1_3.ReadPropertyDefinition(context)
}
