package tosca_v1_2

import (
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/grammars/tosca_v2_0"
)

//
// Artifact
//

// tosca.Reader signature
func ReadArtifact(context *tosca.Context) tosca.EntityPtr {
	context.SetReadTag("ArtifactVersion", "")
	context.SetReadTag("ChecksumAlgorithm", "")
	context.SetReadTag("Checksum", "")

	return tosca_v2_0.ReadArtifact(context)
}
