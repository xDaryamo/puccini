package tosca_v1_2

import (
	"github.com/tliron/kutil/logging"
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/grammars/tosca_v1_3"
	"github.com/tliron/puccini/tosca/grammars/tosca_v2_0"
)

var log = logging.GetLogger("puccini.grammars.tosca_v1_2")
var logRender = logging.NewSubLogger(log, "render")

var Grammar = tosca.NewGrammar()

var DefaultScriptletNamespace = tosca.NewScriptletNamespace()

func init() {
	Grammar.RegisterVersion("tosca_definitions_version", "tosca_simple_yaml_1_2", "/tosca/simple/1.2/profile.yaml")
	Grammar.RegisterVersion("tosca_definitions_version", "tosca_simple_profile_for_nfv_1_0", "/tosca/simple-for-nfv/1.0/profile.yaml")

	Grammar.RegisterReader("$Root", ReadServiceTemplate) // override
	Grammar.RegisterReader("$Unit", ReadUnit)            // override

	Grammar.RegisterReader("Artifact", ReadArtifact)                     // override
	Grammar.RegisterReader("ArtifactDefinition", ReadArtifactDefinition) // override
	Grammar.RegisterReader("ArtifactType", tosca_v2_0.ReadArtifactType)
	Grammar.RegisterReader("AttributeDefinition", ReadAttributeDefinition) // override
	Grammar.RegisterReader("AttributeValue", tosca_v2_0.ReadAttributeValue)
	Grammar.RegisterReader("CapabilityAssignment", ReadCapabilityAssignment) // override
	Grammar.RegisterReader("CapabilityDefinition", tosca_v2_0.ReadCapabilityDefinition)
	Grammar.RegisterReader("CapabilityFilter", tosca_v2_0.ReadCapabilityFilter)
	Grammar.RegisterReader("CapabilityMapping", tosca_v2_0.ReadCapabilityMapping)
	Grammar.RegisterReader("CapabilityType", tosca_v2_0.ReadCapabilityType)
	Grammar.RegisterReader("ConditionClause", tosca_v2_0.ReadConditionClause)
	Grammar.RegisterReader("ConstraintClause", tosca_v2_0.ReadConstraintClause)
	Grammar.RegisterReader("DataType", ReadDataType) // override
	Grammar.RegisterReader("EventFilter", tosca_v2_0.ReadEventFilter)
	Grammar.RegisterReader("Group", ReadGroup)                                  // override
	Grammar.RegisterReader("GroupType", ReadGroupType)                          // override
	Grammar.RegisterReader("Import", tosca_v1_3.ReadImport)                     // 1.3
	Grammar.RegisterReader("InterfaceAssignment", ReadInterfaceAssignment)      // override
	Grammar.RegisterReader("InterfaceDefinition", ReadInterfaceDefinition)      // override
	Grammar.RegisterReader("InterfaceMapping", tosca_v2_0.ReadInterfaceMapping) // introduced in TOSCA 1.2
	Grammar.RegisterReader("InterfaceType", ReadInterfaceType)                  // override
	Grammar.RegisterReader("Metadata", tosca_v2_0.ReadMetadata)
	Grammar.RegisterReader("NodeFilter", tosca_v2_0.ReadNodeFilter)
	Grammar.RegisterReader("NodeTemplate", tosca_v2_0.ReadNodeTemplate)
	Grammar.RegisterReader("NodeType", tosca_v2_0.ReadNodeType)
	Grammar.RegisterReader("OperationAssignment", ReadOperationAssignment) // override
	Grammar.RegisterReader("OperationDefinition", ReadOperationDefinition) // override
	Grammar.RegisterReader("InterfaceImplementation", tosca_v2_0.ReadInterfaceImplementation)
	Grammar.RegisterReader("ParameterDefinition", tosca_v2_0.ReadParameterDefinition)
	Grammar.RegisterReader("Policy", tosca_v2_0.ReadPolicy)
	Grammar.RegisterReader("PolicyType", tosca_v2_0.ReadPolicyType)
	Grammar.RegisterReader("PropertyDefinition", ReadPropertyDefinition) // override
	Grammar.RegisterReader("PropertyFilter", tosca_v2_0.ReadPropertyFilter)
	Grammar.RegisterReader("PropertyMapping", tosca_v2_0.ReadPropertyMapping) // introduced in TOSCA 1.2
	Grammar.RegisterReader("range", tosca_v2_0.ReadRange)
	Grammar.RegisterReader("RangeEntity", tosca_v2_0.ReadRangeEntity)
	Grammar.RegisterReader("RelationshipAssignment", tosca_v2_0.ReadRelationshipAssignment)
	Grammar.RegisterReader("RelationshipDefinition", tosca_v2_0.ReadRelationshipDefinition)
	Grammar.RegisterReader("RelationshipTemplate", tosca_v2_0.ReadRelationshipTemplate)
	Grammar.RegisterReader("RelationshipType", tosca_v2_0.ReadRelationshipType)
	Grammar.RegisterReader("Repository", tosca_v2_0.ReadRepository)
	Grammar.RegisterReader("RequirementAssignment", ReadRequirementAssignment) // override
	Grammar.RegisterReader("RequirementDefinition", tosca_v2_0.ReadRequirementDefinition)
	Grammar.RegisterReader("RequirementMapping", tosca_v2_0.ReadRequirementMapping)
	Grammar.RegisterReader("scalar-unit.frequency", tosca_v2_0.ReadScalarUnitFrequency)
	Grammar.RegisterReader("scalar-unit.size", tosca_v2_0.ReadScalarUnitSize)
	Grammar.RegisterReader("scalar-unit.time", tosca_v2_0.ReadScalarUnitTime)
	Grammar.RegisterReader("Schema", tosca_v2_0.ReadSchema)
	Grammar.RegisterReader("SubstitutionMappings", ReadSubstitutionMappings) // override
	Grammar.RegisterReader("timestamp", tosca_v2_0.ReadTimestamp)
	Grammar.RegisterReader("TopologyTemplate", tosca_v2_0.ReadTopologyTemplate)
	Grammar.RegisterReader("TriggerDefinition", ReadTriggerDefinition) // override
	Grammar.RegisterReader("TriggerDefinitionCondition", tosca_v2_0.ReadTriggerDefinitionCondition)
	Grammar.RegisterReader("Value", tosca_v2_0.ReadValue)
	Grammar.RegisterReader("version", tosca_v2_0.ReadVersion)
	Grammar.RegisterReader("WorkflowActivityDefinition", tosca_v2_0.ReadWorkflowActivityDefinition)         // introduced in TOSCA 1.1
	Grammar.RegisterReader("WorkflowDefinition", tosca_v2_0.ReadWorkflowDefinition)                         // introduced in TOSCA 1.1
	Grammar.RegisterReader("WorkflowPreconditionDefinition", tosca_v2_0.ReadWorkflowPreconditionDefinition) // introduced in TOSCA 1.1
	Grammar.RegisterReader("WorkflowStepDefinition", tosca_v2_0.ReadWorkflowStepDefinition)                 // introduced in TOSCA 1.1

	DefaultScriptletNamespace.RegisterScriptlets(tosca_v2_0.FunctionScriptlets, nil, "tosca.function.join")
	DefaultScriptletNamespace.RegisterScriptlets(tosca_v2_0.ConstraintClauseScriptlets, tosca_v2_0.ConstraintClauseNativeArgumentIndexes)
}
