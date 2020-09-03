package tosca_v1_0

import (
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/grammars/tosca_v1_2"
)

//
// Group
//
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.7.5
//

// tosca.Reader signature
func ReadGroup(context *tosca.Context) tosca.EntityPtr {
	context.SetReadTag("Metadata", "")

	return tosca_v1_2.ReadGroup(context)
}
