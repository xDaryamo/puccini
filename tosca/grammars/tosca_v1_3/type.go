package tosca_v1_3

import (
	"github.com/tliron/puccini/tosca"
)

//
// Type
//
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.7.1
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.7.1
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.6.1
//

type Type struct {
	*Entity `json:"-" yaml:"-"`
	Name    string `namespace:""`

	ParentName  *string  `read:"derived_from"`
	Version     *Version `read:"version,version"`
	Metadata    Metadata `read:"metadata,!Metadata"`
	Description *string  `read:"description"`
}

func NewType(context *tosca.Context) *Type {
	return &Type{
		Entity: NewEntity(context),
		Name:   context.Name,
	}
}

// tosca.HasMetadata interface
func (self *Type) GetDescription() (string, bool) {
	if self.Description != nil {
		return *self.Description, true
	}
	return "", false
}

// tosca.HasMetadata interface
func (self *Type) GetMetadata() (map[string]string, bool) {
	metadata := make(map[string]string)
	if self.Metadata != nil {
		for key, value := range self.Metadata {
			metadata[key] = value
		}
	}
	return metadata, true
}

// tosca.HasMetadata interface
func (self *Type) SetMetadata(name string, value string) bool {
	self.Metadata[name] = value
	return true
}

func (self *Type) GetMetadataValue(key string) (string, bool) {
	if self.Metadata != nil {
		value, ok := self.Metadata[key]
		return value, ok
	}
	return "", false
}
