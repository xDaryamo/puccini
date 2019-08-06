package tosca_v1_3

import (
	"path/filepath"

	"github.com/tliron/puccini/tosca"
)

//
// ArtifactDefinition
//
// Attaches to NodeType
//
// (See Artifact for a variation that attaches to NodeTemplate)
//
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.7
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.6
//

type ArtifactDefinition struct {
	*Entity `name:"artifact"`
	Name    string

	ArtifactTypeName *string `read:"type"` // required only if cannot be inherited
	Description      *string `read:"description" inherit:"description,ArtifactType"`
	Properties       Values  `read:"properties,Value"`
	RepositoryName   *string `read:"repository"`
	File             *string `read:"file"` // required only if cannot be inherited
	DeployPath       *string `read:"deploy_path"`

	ArtifactType *ArtifactType `lookup:"type,ArtifactTypeName" json:"-" yaml:"-"`
	Repository   *Repository   `lookup:"repository,RepositoryName" json:"-" yaml:"-"`

	fileMissingProblemReported bool
}

func NewArtifactDefinition(context *tosca.Context) *ArtifactDefinition {
	return &ArtifactDefinition{
		Entity:     NewEntity(context),
		Name:       context.Name,
		Properties: make(Values),
	}
}

// tosca.Reader signature
func ReadArtifactDefinition(context *tosca.Context) interface{} {
	self := NewArtifactDefinition(context)

	if context.Is("map") {
		// Long notation
		context.ValidateUnsupportedFields(context.ReadFields(self))
	} else if context.ValidateType("map", "string") {
		// Short notation
		self.File = context.FieldChild("file", context.Data).ReadString()
		// TODO: infer ArtifactTypeName from content's URI
	}

	return self
}

func (self *ArtifactDefinition) GetExtension() string {
	if self.File == nil {
		return ""
	}
	extension := filepath.Ext(*self.File)
	if extension == "" {
		return ""
	}
	return extension[1:]
}

// tosca.Mappable interface
func (self *ArtifactDefinition) GetKey() string {
	return self.Name
}

func (self *ArtifactDefinition) Inherit(parentDefinition *ArtifactDefinition) {
	if parentDefinition != nil {
		if (self.ArtifactTypeName == nil) && (parentDefinition.ArtifactTypeName != nil) {
			self.ArtifactTypeName = parentDefinition.ArtifactTypeName
		}
		if (self.Description == nil) && (parentDefinition.Description != nil) {
			self.Description = parentDefinition.Description
		}
		if (self.Properties == nil) && (parentDefinition.Properties != nil) {
			self.Properties = parentDefinition.Properties
		}
		if (self.RepositoryName == nil) && (parentDefinition.RepositoryName != nil) {
			self.RepositoryName = parentDefinition.RepositoryName
		}
		if (self.File == nil) && (parentDefinition.File != nil) {
			self.File = parentDefinition.File
		}
		if (self.DeployPath == nil) && (parentDefinition.DeployPath != nil) {
			self.DeployPath = parentDefinition.DeployPath
		}
		if (self.ArtifactType == nil) && (parentDefinition.ArtifactType != nil) {
			self.ArtifactType = parentDefinition.ArtifactType
		}
		if (self.Repository == nil) && (parentDefinition.Repository != nil) {
			self.Repository = parentDefinition.Repository
		}

		// Validate type compatibility
		if (self.ArtifactType != nil) && (parentDefinition.ArtifactType != nil) && !self.Context.Hierarchy.IsCompatible(parentDefinition.ArtifactType, self.ArtifactType) {
			self.Context.ReportIncompatibleType(self.ArtifactType.Name, parentDefinition.ArtifactType.Name)
		}
	}

	if self.File == nil {
		// Avoid reporting more than once
		if !self.fileMissingProblemReported {
			self.Context.FieldChild("file", nil).ReportFieldMissing()
			self.fileMissingProblemReported = true
		}
	}
}

//
// ArtifactDefinitions
//

type ArtifactDefinitions map[string]*ArtifactDefinition

func (self ArtifactDefinitions) Inherit(parentDefinitions ArtifactDefinitions) {
	for name, definition := range parentDefinitions {
		if _, ok := self[name]; !ok {
			self[name] = definition
		}
	}

	for name, definition := range self {
		if parentDefinition, ok := parentDefinitions[name]; ok {
			if definition != parentDefinition {
				definition.Inherit(parentDefinition)
			}
		} else {
			definition.Inherit(nil)
		}
	}
}
