package tosca_v2_0

import (
	"strings"

	"github.com/tliron/kutil/ard"
	"github.com/tliron/puccini/tosca"
)

//
// Metadata
//

type Metadata map[string]string

// tosca.Reader signature
func ReadMetadata(context *tosca.Context) tosca.EntityPtr {
	var self Metadata

	if context.Is(ard.TypeMap) {
		metadata := context.ReadStringStringMap()
		if metadata != nil {
			self = *metadata
		}
	}

	if self != nil {
		for key, value := range self {
			if strings.HasPrefix(key, "puccini.scriptlet.import:") {
				context.ImportScriptlet(key[25:], value)
				delete(self, key)
			} else if strings.HasPrefix(key, "puccini.scriptlet:") {
				context.EmbedScriptlet(key[18:], value)
				delete(self, key)
			}
		}
	}

	return self
}
