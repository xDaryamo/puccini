package tosca_v2_0

import (
	"strings"

	"github.com/tliron/go-ard"
	"github.com/tliron/puccini/tosca/parsing"
)

//
// Metadata
//

type Metadata map[string]string

// parsing.Reader signature
func ReadMetadata(context *parsing.Context) parsing.EntityPtr {
	var self Metadata

	if context.Is(ard.TypeMap) {
		metadata := context.ReadStringStringMap()
		if metadata != nil {
			self = *metadata
		}
	}

	if self != nil {
		for key, value := range self {
			if strings.HasPrefix(key, parsing.METADATA_SCRIPTLET_IMPORT_PREFIX) {
				context.ImportScriptlet(key[len(parsing.METADATA_SCRIPTLET_IMPORT_PREFIX):], value)
				delete(self, key)
			} else if strings.HasPrefix(key, parsing.METADATA_SCRIPTLET_PREFIX) {
				context.EmbedScriptlet(key[len(parsing.METADATA_SCRIPTLET_PREFIX):], value)
				delete(self, key)
			}
		}
	}

	return self
}
