package tosca_v1_3

import (
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
	*ArtifactDefinition
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
	if self.ArtifactType == nil {
		return
	}

	if definition != nil {
		// Our artifact type must be compatible with definition's
		if (definition.ArtifactType != nil) && !self.Context.Hierarchy.IsCompatible(definition.ArtifactType, self.ArtifactType) {
			self.Context.ReportIncompatible(self.ArtifactType.Name, "artifact", "type")
		}

		// Copy values from definition
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
	}

	self.Properties.RenderProperties(self.ArtifactType.PropertyDefinitions, "property", self.Context.FieldChild("properties", nil))

	// Validate extension (if "file_ext" was not set in type, then anything goes)
	if self.ArtifactType.FileExt != nil {
		extension := self.GetExtension()
		found := false
		for _, e := range *self.ArtifactType.FileExt {
			if e == extension {
				found = true
				break
			}
		}
		if !found {
			self.Context.FieldChild("file", nil).ReportIncompatibleExtension(extension, *self.ArtifactType.FileExt)
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
		a.SourcePath = *self.File
	}
	if self.DeployPath != nil {
		a.TargetPath = *self.DeployPath
	}

	return a
}

//
// Artifacts
//

type Artifacts map[string]*Artifact

func (self Artifacts) Render(definitions ArtifactDefinitions, context *tosca.Context) {
	for key, definition := range definitions {
		artifact, ok := self[key]
		if ok {
			artifact.Render(definition)
		} else {
			artifact := NewArtifact(context.MapChild(key, nil))
			artifact.ArtifactDefinition = definition
			self[key] = artifact
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
