package tosca_v1_0

import (
	"github.com/op/go-logging"
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/grammars/tosca_v1_1"
	"github.com/tliron/puccini/tosca/grammars/tosca_v1_2"
	"github.com/tliron/puccini/tosca/grammars/tosca_v1_3"
)

var log = logging.MustGetLogger("grammars.tosca_v1_0")

var Grammar = tosca.NewGrammar()

func init() {
	Grammar.RegisterVersion("tosca_definitions_version", "tosca_simple_yaml_1_0", "/tosca/simple/1.0/profile.yaml")

	Grammar.RegisterReader("$Root", tosca_v1_1.ReadServiceTemplate) // 1.1
	Grammar.RegisterReader("$Unit", tosca_v1_1.ReadUnit)            // 1.1

	Grammar.RegisterReader("Artifact", tosca_v1_2.ReadArtifact)                     // 1.2
	Grammar.RegisterReader("ArtifactDefinition", tosca_v1_2.ReadArtifactDefinition) // 1.2
	Grammar.RegisterReader("ArtifactType", tosca_v1_3.ReadArtifactType)
	Grammar.RegisterReader("AttributeDefinition", tosca_v1_2.ReadAttributeDefinition) // 1.2
	Grammar.RegisterReader("AttributeValue", tosca_v1_3.ReadAttributeValue)
	Grammar.RegisterReader("CapabilityAssignment", tosca_v1_3.ReadCapabilityAssignment)
	Grammar.RegisterReader("CapabilityDefinition", tosca_v1_3.ReadCapabilityDefinition)
	Grammar.RegisterReader("CapabilityFilter", tosca_v1_3.ReadCapabilityFilter)
	Grammar.RegisterReader("CapabilityMapping", tosca_v1_3.ReadCapabilityMapping)
	Grammar.RegisterReader("CapabilityType", tosca_v1_3.ReadCapabilityType)
	Grammar.RegisterReader("ConstraintClause", tosca_v1_3.ReadConstraintClause)
	Grammar.RegisterReader("DataType", tosca_v1_3.ReadDataType)
	Grammar.RegisterReader("Group", tosca_v1_2.ReadGroup)         // 1.2
	Grammar.RegisterReader("GroupType", tosca_v1_2.ReadGroupType) // 1.2
	Grammar.RegisterReader("Import", tosca_v1_3.ReadImport)
	Grammar.RegisterReader("InterfaceAssignment", tosca_v1_2.ReadInterfaceAssignment) // 1.2
	Grammar.RegisterReader("InterfaceDefinition", tosca_v1_2.ReadInterfaceDefinition) // 1.2
	Grammar.RegisterReader("InterfaceType", tosca_v1_2.ReadInterfaceType)             // 1.2
	Grammar.RegisterReader("Metadata", tosca_v1_3.ReadMetadata)
	Grammar.RegisterReader("NodeFilter", tosca_v1_3.ReadNodeFilter)
	Grammar.RegisterReader("NodeTemplate", tosca_v1_3.ReadNodeTemplate)
	Grammar.RegisterReader("NodeType", tosca_v1_3.ReadNodeType)
	Grammar.RegisterReader("OperationAssignment", tosca_v1_3.ReadOperationAssignment)
	Grammar.RegisterReader("OperationDefinition", tosca_v1_3.ReadOperationDefinition)
	Grammar.RegisterReader("InterfaceImplementation", tosca_v1_1.ReadInterfaceImplementation) // 1.1
	Grammar.RegisterReader("ParameterDefinition", tosca_v1_3.ReadParameterDefinition)
	Grammar.RegisterReader("Policy", ReadPolicy)                                    // override
	Grammar.RegisterReader("PolicyType", ReadPolicyType)                            // override
	Grammar.RegisterReader("PropertyDefinition", tosca_v1_2.ReadPropertyDefinition) // 1.2
	Grammar.RegisterReader("PropertyFilter", tosca_v1_3.ReadPropertyFilter)
	Grammar.RegisterReader("range", tosca_v1_3.ReadRange)
	Grammar.RegisterReader("RangeEntity", tosca_v1_3.ReadRangeEntity)
	Grammar.RegisterReader("RelationshipAssignment", tosca_v1_3.ReadRelationshipAssignment)
	Grammar.RegisterReader("RelationshipDefinition", tosca_v1_3.ReadRelationshipDefinition)
	Grammar.RegisterReader("RelationshipTemplate", tosca_v1_3.ReadRelationshipTemplate)
	Grammar.RegisterReader("RelationshipType", tosca_v1_3.ReadRelationshipType)
	Grammar.RegisterReader("Repository", tosca_v1_3.ReadRepository)
	Grammar.RegisterReader("RequirementAssignment", tosca_v1_3.ReadRequirementAssignment)
	Grammar.RegisterReader("RequirementDefinition", tosca_v1_3.ReadRequirementDefinition)
	Grammar.RegisterReader("RequirementMapping", tosca_v1_3.ReadRequirementMapping)
	Grammar.RegisterReader("scalar-unit.frequency", tosca_v1_3.ReadScalarUnitFrequency)
	Grammar.RegisterReader("scalar-unit.size", tosca_v1_3.ReadScalarUnitSize)
	Grammar.RegisterReader("scalar-unit.time", tosca_v1_3.ReadScalarUnitTime)
	Grammar.RegisterReader("Schema", tosca_v1_3.ReadSchema)
	Grammar.RegisterReader("SubstitutionMappings", tosca_v1_1.ReadSubstitutionMappings) // 1.1
	Grammar.RegisterReader("timestamp", tosca_v1_3.ReadTimestamp)
	Grammar.RegisterReader("TopologyTemplate", ReadTopologyTemplate) // override
	Grammar.RegisterReader("Value", tosca_v1_3.ReadValue)
	Grammar.RegisterReader("version", tosca_v1_3.ReadVersion)
}
