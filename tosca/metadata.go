package tosca

//
// HasMetadata
//

type HasMetadata interface {
	GetMetadata() (map[string]string, bool)
	SetMetadata(name string, value string)
}

// From HasMetadata interface
func GetMetadata(entityPtr interface{}) (map[string]string, bool) {
	if hasMetadata, ok := entityPtr.(HasMetadata); ok {
		if metadata, ok := hasMetadata.GetMetadata(); ok {
			return metadata, true
		}
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
