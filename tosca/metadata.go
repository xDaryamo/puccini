package tosca

import (
	"strings"
)

//
// HasMetadata
//

// Must be thread-safe!
type HasMetadata interface {
	GetDescription() (string, bool)
	GetMetadata() (map[string]string, bool) // should return a copy
	SetMetadata(name string, value string) bool
}

// From HasMetadata interface
func GetDescription(entityPtr EntityPtr) (string, bool) {
	if hasMetadata, ok := entityPtr.(HasMetadata); ok {
		return hasMetadata.GetDescription()
	}
	return "", false
}

// From HasMetadata interface
func GetMetadata(entityPtr EntityPtr) (map[string]string, bool) {
	if hasMetadata, ok := entityPtr.(HasMetadata); ok {
		return hasMetadata.GetMetadata()
	}
	return nil, false
}

// From HasMetadata interface
func SetMetadata(entityPtr EntityPtr, name string, value string) bool {
	if hasMetadata, ok := entityPtr.(HasMetadata); ok {
		hasMetadata.SetMetadata(name, value)
		return true
	}
	return false
}

func GetValueMetadata(metadata map[string]string) map[string]string {
	valueMetadata := make(map[string]string)
	if metadata != nil {
		for key, value := range metadata {
			if strings.HasPrefix(key, "puccini.information:") {
				valueMetadata[key[20:]] = value
			}
		}
	}
	return valueMetadata
}
