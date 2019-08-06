package tosca_v1_3

import (
	"github.com/tliron/puccini/tosca"
)

//
// Type
//
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.7.1
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.6.1
//

type Type struct {
	*Entity `json:"-" yaml:"-"`
	Name    string `namespace:""`

	ParentName  *string  `read:"derived_from"`
	Version     *Version `read:"version,version"`
	Metadata    Metadata `read:"metadata,Metadata"`
	Description *string  `read:"description" inherit:"description,Parent"`
}

func NewType(context *tosca.Context) *Type {
	return &Type{
		Entity: NewEntity(context),
		Name:   context.Name,
	}
}

// normal.HasMetadata interface
func (self *Type) GetMetadata() (map[string]string, bool) {
	metadata := make(map[string]string)
	if self.Metadata != nil {
		for k, v := range self.Metadata {
			metadata[k] = v
		}
	}
	return metadata, true
}

func (self *Type) GetMetadataValue(key string) (string, bool) {
	if self.Metadata != nil {
		value, ok := self.Metadata[key]
		return value, ok
	}
	return "", false
}
