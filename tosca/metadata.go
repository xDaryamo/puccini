package tosca

//
// HasMetadata
//

type HasMetadata interface {
	GetDescription() (string, bool)
	GetMetadata() (map[string]string, bool)
	SetMetadata(name string, value string) bool
}

// From HasMetadata interface
func GetDescription(entityPtr interface{}) (string, bool) {
	if hasMetadata, ok := entityPtr.(HasMetadata); ok {
		return hasMetadata.GetDescription()
	}
	return "", false
}

// From HasMetadata interface
func GetMetadata(entityPtr interface{}) (map[string]string, bool) {
	if hasMetadata, ok := entityPtr.(HasMetadata); ok {
		return hasMetadata.GetMetadata()
	}
	return nil, false
}

// From HasMetadata interface
func SetMetadata(entityPtr interface{}, name string, value string) bool {
	if hasMetadata, ok := entityPtr.(HasMetadata); ok {
		hasMetadata.SetMetadata(name, value)
		return true
	}
	return false
}
