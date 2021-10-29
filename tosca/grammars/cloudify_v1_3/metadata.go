package cloudify_v1_3

import (
	"strings"

	"github.com/tliron/kutil/ard"
	"github.com/tliron/puccini/tosca"
)

//
// Metadata
//
// Note: not in spec
//

type Metadata map[string]ard.Value

// tosca.Reader signature
func ReadMetadata(context *tosca.Context) tosca.EntityPtr {
	var self map[string]ard.Value

	if context.ValidateType(ard.TypeMap) {
		metadata := context.ReadStringMap()
		if metadata != nil {
			self = *metadata
		}
	}

	if self != nil {
		for key, value := range self {
			if strings.HasPrefix(key, tosca.METADATA_SCRIPTLET_IMPORT_PREFIX) {
				if v, ok := value.(string); ok {
					context.ImportScriptlet(key[len(tosca.METADATA_SCRIPTLET_IMPORT_PREFIX):], v)
					delete(self, key)
				} else {
					context.MapChild(key, value).ReportValueWrongType(ard.TypeString)
				}
			} else if strings.HasPrefix(key, tosca.METADATA_SCRIPTLET_PREFIX) {
				if v, ok := value.(string); ok {
					context.EmbedScriptlet(key[len(tosca.METADATA_SCRIPTLET_PREFIX):], v)
					delete(self, key)
				} else {
					context.MapChild(key, value).ReportValueWrongType(ard.TypeString)
				}
			}
		}
	}

	return self
}
