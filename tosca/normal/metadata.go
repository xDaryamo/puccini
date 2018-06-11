package normal

//
// HasMetadata
//

type HasMetadata interface {
	GetMetadata() (map[string]string, bool)
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
