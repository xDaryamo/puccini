package tosca_v2_0

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
func ReadArtifact(context *tosca.Context) tosca.EntityPtr {
	self := NewArtifact(context)
	self.ArtifactDefinition = ReadArtifactDefinition(context).(*ArtifactDefinition)
	return self
}

func (self *Artifact) Copy(definition *ArtifactDefinition) {
	// Validate type compatibility
	if (self.ArtifactType != nil) && (definition.ArtifactType != nil) && !self.Context.Hierarchy.IsCompatible(definition.ArtifactType, self.ArtifactType) {
		self.Context.ReportIncompatible(tosca.GetCanonicalName(self.ArtifactType), "artifact", "type")
	}

	if self.ArtifactType == nil {
		self.ArtifactType = definition.ArtifactType
	}
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

func (self *Artifact) DoRender() {
	logRender.Debugf("artifact: %s", self.Name)

	if self.File == nil {
		self.Context.FieldChild("file", nil).ReportFieldMissing()
	}

	if self.ArtifactType == nil {
		return
	}

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

	self.Properties.RenderProperties(self.ArtifactType.PropertyDefinitions, "property", self.Context.FieldChild("properties", nil))
}

func (self *Artifact) Normalize(normalNodeTemplate *normal.NodeTemplate) *normal.Artifact {
	logNormalize.Debugf("artifact: %s", self.Name)

	normalArtifact := normalNodeTemplate.NewArtifact(self.Name)

	if self.Description != nil {
		normalArtifact.Description = *self.Description
	}

	if types, ok := normal.GetTypes(self.Context.Hierarchy, self.ArtifactType); ok {
		normalArtifact.Types = types
	}

	self.Properties.Normalize(normalArtifact.Properties)

	if self.File != nil {
		normalArtifact.Filename = filepath.Base(*self.File)
	}
	url := self.GetURL()
	if url != nil {
		normalArtifact.SourcePath = url.String()
	}
	if self.DeployPath != nil {
		normalArtifact.TargetPath = *self.DeployPath
	}
	if self.ArtifactVersion != nil {
		normalArtifact.Version = *self.ArtifactVersion
	}
	if self.ChecksumAlgorithm != nil {
		normalArtifact.ChecksumAlgorithm = *self.ChecksumAlgorithm
	}
	if self.Checksum != nil {
		normalArtifact.Checksum = *self.Checksum
	}
	if (self.Repository != nil) && (self.Repository.Credential != nil) {
		normalArtifact.Credential = self.Repository.Credential.Normalize()
	}

	return normalArtifact
}

//
// Artifacts
//

type Artifacts map[string]*Artifact

func (self Artifacts) Render(definitions ArtifactDefinitions, context *tosca.Context) {
	for key, definition := range definitions {
		if artifact, ok := self[key]; ok {
			artifact.Copy(definition)
		}
	}

	for key, artifact := range self {
		if _, ok := definitions[key]; !ok {
			artifact.DoRender()
		}
	}
}

func (self Artifacts) Normalize(normalNodeTemplate *normal.NodeTemplate) {
	for key, artifact := range self {
		normalNodeTemplate.Artifacts[key] = artifact.Normalize(normalNodeTemplate)
	}
}
