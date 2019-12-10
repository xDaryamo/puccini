package cloudify_v1_3

import (
	"strings"

	"github.com/tliron/puccini/tosca"
)

//
// Metadata
//
// Note: not in spec
//

type Metadata map[string]interface{}

// tosca.Reader signature
func ReadMetadata(context *tosca.Context) interface{} {
	var self map[string]interface{}

	if context.ValidateType("map") {
		metadata := context.ReadStringMap()
		if metadata != nil {
			self = *metadata
		}
	}

	if self != nil {
		for key, value := range self {
			if strings.HasPrefix(key, "puccini.scriptlet.import.") {
				if v, ok := value.(string); ok {
					context.ImportScriptlet(key[25:], v)
					delete(self, key)
				} else {
					context.MapChild(key, value).ReportValueWrongType("string")
				}
			} else if strings.HasPrefix(key, "puccini.scriptlet.") {
				if v, ok := value.(string); ok {
					context.EmbedScriptlet(key[18:], v)
					delete(self, key)
				} else {
					context.MapChild(key, value).ReportValueWrongType("string")
				}
			}
		}
	}

	return self
}
