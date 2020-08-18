package tosca_v1_2

import (
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/grammars/tosca_v2_0"
)

//
// ArtifactDefinition
//
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.7
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.6
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.5.6
//

// tosca.Reader signature
func ReadArtifactDefinition(context *tosca.Context) tosca.EntityPtr {
	context.SetReadTag("ArtifactVersion", "")
	context.SetReadTag("ChecksumAlgorithm", "")
	context.SetReadTag("Checksum", "")

	return tosca_v2_0.ReadArtifactDefinition(context)
}
