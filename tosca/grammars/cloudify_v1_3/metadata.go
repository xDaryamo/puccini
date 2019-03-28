package cloudify_v1_3

import (
	"strings"

	"github.com/tliron/puccini/ard"
	"github.com/tliron/puccini/tosca"
)

//
// Metadata
//
// Note: not in spec
//

type Metadata ard.Map

// tosca.Reader signature
func ReadMetadata(context *tosca.Context) interface{} {
	var self map[string]interface{}

	if context.ValidateType("map") {
		self = context.Data.(ard.Map)
	}

	if self != nil {
		for k, v := range self {
			if strings.HasPrefix(k, "puccini-js.import.") {
				name := k[18:]
				context.ImportScript(name, v.(string))
				delete(self, k)
			} else if strings.HasPrefix(k, "puccini-js.source.") {
				name := k[18:]
				context.SourceScript(name, v.(string))
				delete(self, k)
			}
		}
	}

	return self
}
