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
	var self map[string]string

	if context.ValidateType("map") {
		metadata := context.ReadStringMap()
		if metadata != nil {
			self = *metadata
		}
	}

	if self != nil {
		for k, v := range self {
			if strings.HasPrefix(k, "puccini-js.import.") {
				name := k[18:]
				context.ImportScript(name, v)
				delete(self, k)
			} else if strings.HasPrefix(k, "puccini-js.source.") {
				name := k[18:]
				context.SourceScript(name, v)
				delete(self, k)
			}
		}
	}

	return self
}
