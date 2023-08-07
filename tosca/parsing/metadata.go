package parsing

import (
	"strings"
)

//
// HasMetadata
//

const (
	METADATA_TYPE                    = "puccini.type"
	METADATA_CONVERTER               = "puccini.converter"
	METADATA_COMPARER                = "puccini.comparer"
	METADATA_QUIRKS                  = "puccini.quirks"
	METADATA_DATA_TYPE_PREFIX        = "puccini.data-type:"
	METADATA_SCRIPTLET_PREFIX        = "puccini.scriptlet:"
	METADATA_SCRIPTLET_IMPORT_PREFIX = "puccini.scriptlet.import:"

	METADATA_CANONICAL_NAME    = "tosca.canonical-name"
	METADATA_NORMATIVE         = "tosca.normative"
	METADATA_FUNCTION_PREFIX   = "tosca.function."
	METADATA_CONSTRAINT_PREFIX = "tosca.constraint."
)

type HasMetadata interface {
	GetDescription() (string, bool)
	GetMetadata() (map[string]string, bool) // should return a copy
	SetMetadata(name string, value string) bool
}

// From HasMetadata interface
func GetDescription(entityPtr EntityPtr) (string, bool) {
	if hasMetadata, ok := entityPtr.(HasMetadata); ok {
		return hasMetadata.GetDescription()
	} else {
		return "", false
	}
}

// From HasMetadata interface
func GetMetadata(entityPtr EntityPtr) (map[string]string, bool) {
	if hasMetadata, ok := entityPtr.(HasMetadata); ok {
		return hasMetadata.GetMetadata()
	} else {
		return nil, false
	}
}

// From HasMetadata interface
func SetMetadata(entityPtr EntityPtr, name string, value string) bool {
	if hasMetadata, ok := entityPtr.(HasMetadata); ok {
		hasMetadata.SetMetadata(name, value)
		return true
	} else {
		return false
	}
}

func GetDataTypeMetadata(metadata map[string]string) map[string]string {
	dataTypeMetadata := make(map[string]string)
	if metadata != nil {
		for key, value := range metadata {
			if strings.HasPrefix(key, METADATA_DATA_TYPE_PREFIX) {
				dataTypeMetadata[key[len(METADATA_DATA_TYPE_PREFIX):]] = value
			}
		}
	}
	return dataTypeMetadata
}
