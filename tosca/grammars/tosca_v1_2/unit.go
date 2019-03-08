package tosca_v1_2

import (
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
)

//
// Unit
//
// See ServiceTemplate
//
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.10
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.9
//

type Unit struct {
	*Entity `name:"unit"`

	ToscaDefinitionsVersion *string             `read:"tosca_definitions_version" require:"tosca_definitions_version"`
	Metadata                Metadata            `read:"metadata,!Metadata"`
	Repositories            []*Repository       `read:"repositories,Repository"`
	Imports                 []*Import           `read:"imports,[]Import"`
	ArtifactTypes           []*ArtifactType     `read:"artifact_types,ArtifactType" hierarchy:""`
	CapabilityTypes         []*CapabilityType   `read:"capability_types,CapabilityType" hierarchy:""`
	DataTypes               []*DataType         `read:"data_types,DataType" hierarchy:""`
	GroupTypes              []*GroupType        `read:"group_types,GroupType" hierarchy:""`
	InterfaceTypes          []*InterfaceType    `read:"interface_types,InterfaceType" hierarchy:""`
	NodeTypes               []*NodeType         `read:"node_types,NodeType" hierarchy:""`
	PolicyTypes             []*PolicyType       `read:"policy_types,PolicyType" hierarchy:""`
	RelationshipTypes       []*RelationshipType `read:"relationship_types,RelationshipType" hierarchy:""`
}

func NewUnit(context *tosca.Context) *Unit {
	return &Unit{Entity: NewEntity(context)}
}

// tosca.Reader signature
func ReadUnit(context *tosca.Context) interface{} {
	self := NewUnit(context)
	context.ScriptNamespace.Merge(DefaultScriptNamespace)
	context.ValidateUnsupportedFields(append(context.ReadFields(self), "dsl_definitions"))
	return self
}

// parser.Importer interface
func (self *Unit) GetImportSpecs() []*tosca.ImportSpec {
	// TODO: importing should also import repositories

	var importSpecs = make([]*tosca.ImportSpec, 0, len(self.Imports))
	for _, import_ := range self.Imports {
		if importSpec, ok := import_.NewImportSpec(self); ok {
			importSpecs = append(importSpecs, importSpec)
		}
	}
	return importSpecs
}

func (self *Unit) Normalize(s *normal.ServiceTemplate) {
	log.Info("{normalize} unit")

	if self.Metadata != nil {
		for k, v := range self.Metadata {
			s.Metadata[k] = v
		}
	}
}
