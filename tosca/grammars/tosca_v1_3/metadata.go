package tosca_v1_3

import (
	"strings"

	"github.com/tliron/puccini/tosca"
)

//
// Metadata
//

type Metadata map[string]string

// tosca.Reader signature
func ReadMetadata(context *tosca.Context) interface{} {
	var self Metadata

	if context.Is("map") {
		metadata := context.ReadStringStringMap()
		if metadata != nil {
			self = *metadata
		}
	}

	if self != nil {
		for key, value := range self {
			if strings.HasPrefix(key, "puccini-js.import.") {
				name := key[18:]
				context.ImportScript(name, value)
				delete(self, key)
			} else if strings.HasPrefix(key, "puccini-js.source.") {
				name := key[18:]
				context.SourceScript(name, value)
				delete(self, key)
			}
		}
	}

	return self
}
