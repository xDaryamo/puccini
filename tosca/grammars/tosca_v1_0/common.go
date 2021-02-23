package tosca_v1_0

import (
	"github.com/tliron/kutil/logging"
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/grammars/tosca_v1_1"
	"github.com/tliron/puccini/tosca/grammars/tosca_v1_2"
	"github.com/tliron/puccini/tosca/grammars/tosca_v1_3"
	"github.com/tliron/puccini/tosca/grammars/tosca_v2_0"
)

var log = logging.GetLogger("puccini.grammars.tosca_v1_0")

var Grammar = tosca.NewGrammar()

func init() {
	Grammar.RegisterVersion("tosca_definitions_version", "tosca_simple_yaml_1_0", "/tosca/simple/1.0/profile.yaml")

	Grammar.RegisterReader("$Root", tosca_v1_1.ReadServiceTemplate) // 1.1
	Grammar.RegisterReader("$Unit", tosca_v1_1.ReadUnit)            // 1.1

	Grammar.RegisterReader("Artifact", tosca_v1_2.ReadArtifact)                     // 1.2
	Grammar.RegisterReader("ArtifactDefinition", tosca_v1_2.ReadArtifactDefinition) // 1.2
	Grammar.RegisterReader("ArtifactType", tosca_v2_0.ReadArtifactType)
	Grammar.RegisterReader("AttributeDefinition", tosca_v1_2.ReadAttributeDefinition) // 1.2
	Grammar.RegisterReader("AttributeValue", tosca_v2_0.ReadAttributeValue)
	Grammar.RegisterReader("CapabilityAssignment", tosca_v1_2.ReadCapabilityAssignment) // 1.2
	Grammar.RegisterReader("CapabilityDefinition", tosca_v2_0.ReadCapabilityDefinition)
	Grammar.RegisterReader("CapabilityFilter", tosca_v2_0.ReadCapabilityFilter)
	Grammar.RegisterReader("CapabilityMapping", tosca_v2_0.ReadCapabilityMapping)
	Grammar.RegisterReader("CapabilityType", tosca_v2_0.ReadCapabilityType)
	Grammar.RegisterReader("ConstraintClause", tosca_v2_0.ReadConstraintClause)
	Grammar.RegisterReader("DataType", tosca_v1_2.ReadDataType)                       // 1.2
	Grammar.RegisterReader("Group", ReadGroup)                                        // override
	Grammar.RegisterReader("GroupType", tosca_v1_2.ReadGroupType)                     // 1.2
	Grammar.RegisterReader("Import", tosca_v1_3.ReadImport)                           /// 1.3
	Grammar.RegisterReader("InterfaceAssignment", tosca_v1_2.ReadInterfaceAssignment) // 1.2
	Grammar.RegisterReader("InterfaceDefinition", tosca_v1_2.ReadInterfaceDefinition) // 1.2
	Grammar.RegisterReader("InterfaceType", tosca_v1_2.ReadInterfaceType)             // 1.2
	Grammar.RegisterReader("Metadata", tosca_v2_0.ReadMetadata)
	Grammar.RegisterReader("NodeFilter", tosca_v2_0.ReadNodeFilter)
	Grammar.RegisterReader("NodeTemplate", ReadNodeTemplate) // override
	Grammar.RegisterReader("NodeType", tosca_v2_0.ReadNodeType)
	Grammar.RegisterReader("OperationAssignment", tosca_v1_2.ReadOperationAssignment)         // 1.2
	Grammar.RegisterReader("OperationDefinition", tosca_v1_2.ReadOperationDefinition)         // 1.2
	Grammar.RegisterReader("InterfaceImplementation", tosca_v1_1.ReadInterfaceImplementation) // 1.1
	Grammar.RegisterReader("ParameterDefinition", tosca_v2_0.ReadParameterDefinition)
	Grammar.RegisterReader("Policy", ReadPolicy)                                    // override
	Grammar.RegisterReader("PolicyType", ReadPolicyType)                            // override
	Grammar.RegisterReader("PropertyDefinition", tosca_v1_2.ReadPropertyDefinition) // 1.2
	Grammar.RegisterReader("PropertyFilter", tosca_v2_0.ReadPropertyFilter)
	Grammar.RegisterReader("range", tosca_v2_0.ReadRange)
	Grammar.RegisterReader("RangeEntity", tosca_v2_0.ReadRangeEntity)
	Grammar.RegisterReader("RelationshipAssignment", tosca_v2_0.ReadRelationshipAssignment)
	Grammar.RegisterReader("RelationshipDefinition", tosca_v2_0.ReadRelationshipDefinition)
	Grammar.RegisterReader("RelationshipTemplate", ReadRelationshipTemplate) // override
	Grammar.RegisterReader("RelationshipType", tosca_v2_0.ReadRelationshipType)
	Grammar.RegisterReader("Repository", tosca_v2_0.ReadRepository)
	Grammar.RegisterReader("RequirementAssignment", tosca_v1_2.ReadRequirementAssignment) // 1.2
	Grammar.RegisterReader("RequirementDefinition", tosca_v2_0.ReadRequirementDefinition)
	Grammar.RegisterReader("RequirementMapping", tosca_v2_0.ReadRequirementMapping)
	Grammar.RegisterReader("scalar-unit.frequency", tosca_v2_0.ReadScalarUnitFrequency)
	Grammar.RegisterReader("scalar-unit.size", tosca_v2_0.ReadScalarUnitSize)
	Grammar.RegisterReader("scalar-unit.time", tosca_v2_0.ReadScalarUnitTime)
	Grammar.RegisterReader("Schema", tosca_v2_0.ReadSchema)
	Grammar.RegisterReader("SubstitutionMappings", tosca_v1_1.ReadSubstitutionMappings) // 1.1
	Grammar.RegisterReader("timestamp", tosca_v2_0.ReadTimestamp)
	Grammar.RegisterReader("TopologyTemplate", ReadTopologyTemplate) // override
	Grammar.RegisterReader("Value", tosca_v2_0.ReadValue)
	Grammar.RegisterReader("version", tosca_v2_0.ReadVersion)
}
