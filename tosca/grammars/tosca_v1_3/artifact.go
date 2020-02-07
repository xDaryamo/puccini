package tosca_v1_3

import (
	"path/filepath"

	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
)

//
// Artifact
//
// Attaches to NodeTemplate
//
// (See ArtifactDefinition for a variation that attaches to NodeType)
//

type Artifact struct {
	*ArtifactDefinition `name:"artifact"`
}

func NewArtifact(context *tosca.Context) *Artifact {
	return &Artifact{ArtifactDefinition: NewArtifactDefinition(context)}
}

// tosca.Reader signature
func ReadArtifact(context *tosca.Context) interface{} {
	self := NewArtifact(context)
	self.ArtifactDefinition = ReadArtifactDefinition(context).(*ArtifactDefinition)
	return self
}

func (self *Artifact) Render(definition *ArtifactDefinition) {
	if definition != nil {
		if self.ArtifactType == nil {
			self.ArtifactType = definition.ArtifactType
		} else {
			// Our artifact type must be compatible with definition's
			if (definition.ArtifactType != nil) && !self.Context.Hierarchy.IsCompatible(definition.ArtifactType, self.ArtifactType) {
				self.Context.ReportIncompatible(self.ArtifactType.Name, "artifact", "type")
			}
		}

		// Copy values from definition
		if self.ArtifactVersion == nil {
			self.ArtifactVersion = definition.ArtifactVersion
		}
		if self.Description == nil {
			self.Description = definition.Description
		}
		for name, value := range definition.Properties {
			if _, ok := self.Properties[name]; !ok {
				self.Properties[name] = value
			}
		}
		if self.RepositoryName == nil {
			self.RepositoryName = definition.RepositoryName
		}
		if self.File == nil {
			self.File = definition.File
		}
		if self.DeployPath == nil {
			self.DeployPath = definition.DeployPath
		}
		if self.Repository == nil {
			self.Repository = definition.Repository
		}
		if self.ChecksumAlgorithm == nil {
			self.ChecksumAlgorithm = definition.ChecksumAlgorithm
		}
		if self.Checksum == nil {
			self.Checksum = definition.Checksum
		}
	}

	if self.File == nil {
		// Avoid reporting more than once
		if !self.fileMissingProblemReported {
			self.Context.FieldChild("file", nil).ReportFieldMissing()
			self.fileMissingProblemReported = true
		}
	}

	if self.ArtifactType == nil {
		return
	}

	self.Properties.RenderProperties(self.ArtifactType.PropertyDefinitions, "property", self.Context.FieldChild("properties", nil))

	// Validate extension (if "file_ext" was not set in type, then anything goes)
	if self.ArtifactType.FileExtension != nil {
		extension := self.GetExtension()
		found := false
		for _, e := range *self.ArtifactType.FileExtension {
			if e == extension {
				found = true
				break
			}
		}
		if !found {
			self.Context.FieldChild("file", nil).ReportIncompatibleExtension(extension, *self.ArtifactType.FileExtension)
		}
	}
}

func (self *Artifact) Normalize(n *normal.NodeTemplate) *normal.Artifact {
	log.Debugf("{normalize} artifact: %s", self.Name)

	a := n.NewArtifact(self.Name)

	if self.Description != nil {
		a.Description = *self.Description
	}

	if types, ok := normal.GetTypes(self.Context.Hierarchy, self.ArtifactType); ok {
		a.Types = types
	}

	self.Properties.Normalize(a.Properties)

	if self.File != nil {
		a.Filename = filepath.Base(*self.File)
	}
	url_ := self.GetURL()
	if url_ != nil {
		a.SourcePath = url_.String()
	}
	if self.DeployPath != nil {
		a.TargetPath = *self.DeployPath
	}
	if self.ArtifactVersion != nil {
		a.Version = *self.ArtifactVersion
	}
	if self.ChecksumAlgorithm != nil {
		a.ChecksumAlgorithm = *self.ChecksumAlgorithm
	}
	if self.Checksum != nil {
		a.Checksum = *self.Checksum
	}
	if (self.Repository != nil) && (self.Repository.Credential != nil) {
		a.Credential = self.Repository.Credential.Normalize()
	}

	return a
}

//
// Artifacts
//

type Artifacts map[string]*Artifact

func (self Artifacts) Render(definitions ArtifactDefinitions, context *tosca.Context) {
	for key, definition := range definitions {
		if artifact, ok := self[key]; ok {
			artifact.Render(definition)
		}
	}

	for key, artifact := range self {
		if _, ok := definitions[key]; !ok {
			artifact.Render(nil)
		}
	}
}

func (self Artifacts) Normalize(n *normal.NodeTemplate) {
	for key, artifact := range self {
		n.Artifacts[key] = artifact.Normalize(n)
	}
}
