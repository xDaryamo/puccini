package tosca_v2_0

import (
	"github.com/tliron/puccini/tosca"
)

//
// ArtifactType
//
// [TOSCA-v2.0] @ ?
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.7.4
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.7.4
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.6.4
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.6.3
//

type ArtifactType struct {
	*Type `name:"artifact type"`

	PropertyDefinitions PropertyDefinitions `read:"properties,PropertyDefinition" inherit:"properties,Parent"`
	MIMEType            *string             `read:"mime_type" inherit:"mime_type,Parent"`
	FileExtension       *[]string           `read:"file_ext" inherit:"file_ext,Parent"`

	Parent *ArtifactType `lookup:"derived_from,ParentName" json:"-" yaml:"-"`
}

func NewArtifactType(context *tosca.Context) *ArtifactType {
	return &ArtifactType{
		Type:                NewType(context),
		PropertyDefinitions: make(PropertyDefinitions),
	}
}

// tosca.Reader signature
func ReadArtifactType(context *tosca.Context) tosca.EntityPtr {
	self := NewArtifactType(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))
	return self
}

// tosca.Hierarchical interface
func (self *ArtifactType) GetParent() tosca.EntityPtr {
	return self.Parent
}

// tosca.Inherits interface
func (self *ArtifactType) Inherit() {
	logInherit.Debugf("artifact type: %s", self.Name)

	if self.Parent == nil {
		return
	}

	self.PropertyDefinitions.Inherit(self.Parent.PropertyDefinitions)
}

//
// ArtifactTypes
//

type ArtifactTypes []*ArtifactType
