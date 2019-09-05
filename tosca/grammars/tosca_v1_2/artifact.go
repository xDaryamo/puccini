package tosca_v1_2

import (
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/grammars/tosca_v1_3"
)

//
// Artifact
//

// tosca.Reader signature
func ReadArtifact(context *tosca.Context) interface{} {
	if context.ReadOverrides == nil {
		context.ReadOverrides = make(map[string]string)
	}
	context.ReadOverrides["ArtifactVersion"] = ""
	context.ReadOverrides["ChecksumAlgorithm"] = ""
	context.ReadOverrides["Checksum"] = ""

	return tosca_v1_3.ReadArtifact(context)
}
