package cloudify_v1_3

import (
	"strings"

	"github.com/tliron/puccini/tosca"
	"github.com/tliron/yamlkeys"
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
			name := yamlkeys.KeyString(key)
			if strings.HasPrefix(name, "puccini-js.import.") {
				name := name[18:]
				context.ImportScriptlet(name, value.(string))
				delete(self, key)
			} else if strings.HasPrefix(key, "puccini-js.embed.") {
				name := name[17:]
				context.EmbedScriptlet(name, value.(string))
				delete(self, key)
			}
		}
	}

	return self
}
