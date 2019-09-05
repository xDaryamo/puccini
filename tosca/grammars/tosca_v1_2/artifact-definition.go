package tosca_v1_2

import (
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/grammars/tosca_v1_3"
)

//
// ArtifactDefinition
//
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.7
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.6
//

// tosca.Reader signature
func ReadArtifactDefinition(context *tosca.Context) interface{} {
	if context.ReadOverrides == nil {
		context.ReadOverrides = make(map[string]string)
	}
	context.ReadOverrides["ArtifactVersion"] = ""
	context.ReadOverrides["ChecksumAlgorithm"] = ""
	context.ReadOverrides["Checksum"] = ""

	return tosca_v1_3.ReadArtifactDefinition(context)
}
