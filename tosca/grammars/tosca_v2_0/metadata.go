package tosca_v2_0

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/tliron/go-ard"
	"github.com/tliron/puccini/tosca/parsing"
	"github.com/tliron/yamlkeys"
)

//
// Metadata
//

type Metadata map[string]string

// ([parsing.Reader] signature)
func ReadMetadata(context *parsing.Context) parsing.EntityPtr {
	var self Metadata

	if context.Is(ard.TypeMap) {
		// Support all YAML 1.2 data types as per TOSCA 2.0 specification
		// !!map, !!seq, !!str, !!null, !!bool, !!int, !!float
		if data, ok := context.Data.(ard.Map); ok {
			self = make(Metadata)
			for key, value := range data {
				if key != nil {
					keyString := yamlkeys.KeyString(key)

					// Convert value to string representation
					var valueString string
					switch v := value.(type) {
					case string:
						valueString = v
					case nil:
						valueString = ""
					case bool:
						if v {
							valueString = "true"
						} else {
							valueString = "false"
						}
					case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
						valueString = fmt.Sprintf("%d", v)
					case float32, float64:
						valueString = fmt.Sprintf("%g", v)
					default:
						// For complex types (maps, lists), convert to JSON
						if jsonBytes, err := json.Marshal(v); err == nil {
							valueString = string(jsonBytes)
						} else {
							valueString = fmt.Sprintf("%v", v)
						}
					}

					self[keyString] = valueString
				}
			}
		}
	}

	if self != nil {
		for key, value := range self {
			// Handle scriptlet imports and embeds (these must be strings)
			if strings.HasPrefix(key, parsing.MetadataScriptletImportPrefix) {
				context.ImportScriptlet(key[len(parsing.MetadataScriptletImportPrefix):], value)
				delete(self, key)
			} else if strings.HasPrefix(key, parsing.MetadataScriptletPrefix) {
				context.EmbedScriptlet(key[len(parsing.MetadataScriptletPrefix):], value)
				delete(self, key)
			}
		}
	}

	return self
}
