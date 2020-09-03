package tosca_v1_0

import (
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/grammars/tosca_v2_0"
)

//
// RelationshipTemplate
//
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.7.4
//

// tosca.Reader signature
func ReadRelationshipTemplate(context *tosca.Context) tosca.EntityPtr {
	context.SetReadTag("Metadata", "")

	return tosca_v2_0.ReadRelationshipTemplate(context)
}
