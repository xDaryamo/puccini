package v1_1

import (
	"strings"

	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"

	simpleForNFV_v1_0 "github.com/tliron/puccini/tosca/profiles/simple-for-nfv/v1_0"
	simple_v1_1 "github.com/tliron/puccini/tosca/profiles/simple/v1_1"
)

//
// Unit
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
	context.ValidateUnsupportedFields(append(context.ReadFields(self, Readers), "dsl_definitions"))
	return self
}

// parser.Importer interface
func (self *Unit) GetImportSpecs() []*tosca.ImportSpec {
	var importSpecs []*tosca.ImportSpec

	// TODO: importing should also import repositories

	// TODO: better way to decide when not to import the profile
	ok := strings.HasPrefix(self.Context.URL.String(), "internal:/tosca/simple/1.1/")
	if !ok {
		importSpecs = append(importSpecs, &tosca.ImportSpec{simple_v1_1.GetURL(), nil})
	}

	if (self.ToscaDefinitionsVersion != nil) && (*self.ToscaDefinitionsVersion == "tosca_simple_profile_for_nfv_1_0") {
		importSpecs = append(importSpecs, &tosca.ImportSpec{simpleForNFV_v1_0.GetURL(), nil})
	}

	// Our imports
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
