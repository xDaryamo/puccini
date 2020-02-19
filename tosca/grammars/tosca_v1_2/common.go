package tosca_v1_2

import (
	"github.com/op/go-logging"
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/grammars/tosca_v1_3"
)

var log = logging.MustGetLogger("grammars.tosca_v1_2")

var Grammar = tosca.NewGrammar()

var DefaultScriptletNamespace = make(tosca.ScriptletNamespace)

func init() {
	Grammar.RegisterVersion("tosca_definitions_version", "tosca_simple_yaml_1_2", "/tosca/simple/1.2/profile.yaml")
	Grammar.RegisterVersion("tosca_definitions_version", "tosca_simple_profile_for_nfv_1_0", "/tosca/simple-for-nfv/1.0/profile.yaml")

	Grammar.RegisterReader("$Root", ReadServiceTemplate) // override
	Grammar.RegisterReader("$Unit", ReadUnit)            // override

	Grammar.RegisterReader("Artifact", ReadArtifact)                     // override
	Grammar.RegisterReader("ArtifactDefinition", ReadArtifactDefinition) // override
	Grammar.RegisterReader("ArtifactType", tosca_v1_3.ReadArtifactType)
	Grammar.RegisterReader("AttributeDefinition", ReadAttributeDefinition) // override
	Grammar.RegisterReader("AttributeValue", tosca_v1_3.ReadAttributeValue)
	Grammar.RegisterReader("CapabilityAssignment", tosca_v1_3.ReadCapabilityAssignment)
	Grammar.RegisterReader("CapabilityDefinition", tosca_v1_3.ReadCapabilityDefinition)
	Grammar.RegisterReader("CapabilityFilter", tosca_v1_3.ReadCapabilityFilter)
	Grammar.RegisterReader("CapabilityMapping", tosca_v1_3.ReadCapabilityMapping)
	Grammar.RegisterReader("CapabilityType", tosca_v1_3.ReadCapabilityType)
	Grammar.RegisterReader("ConditionClause", tosca_v1_3.ReadConditionClause)
	Grammar.RegisterReader("ConstraintClause", tosca_v1_3.ReadConstraintClause)
	Grammar.RegisterReader("DataType", tosca_v1_3.ReadDataType)
	Grammar.RegisterReader("EventFilter", tosca_v1_3.ReadEventFilter)
	Grammar.RegisterReader("Group", ReadGroup)         // override
	Grammar.RegisterReader("GroupType", ReadGroupType) // override
	Grammar.RegisterReader("Import", tosca_v1_3.ReadImport)
	Grammar.RegisterReader("InterfaceAssignment", ReadInterfaceAssignment)      // override
	Grammar.RegisterReader("InterfaceDefinition", ReadInterfaceDefinition)      // override
	Grammar.RegisterReader("InterfaceMapping", tosca_v1_3.ReadInterfaceMapping) // introduced in TOSCA 1.2
	Grammar.RegisterReader("InterfaceType", ReadInterfaceType)                  // override
	Grammar.RegisterReader("Metadata", tosca_v1_3.ReadMetadata)
	Grammar.RegisterReader("NodeFilter", tosca_v1_3.ReadNodeFilter)
	Grammar.RegisterReader("NodeTemplate", tosca_v1_3.ReadNodeTemplate)
	Grammar.RegisterReader("NodeType", tosca_v1_3.ReadNodeType)
	Grammar.RegisterReader("OperationAssignment", tosca_v1_3.ReadOperationAssignment)
	Grammar.RegisterReader("OperationDefinition", tosca_v1_3.ReadOperationDefinition)
	Grammar.RegisterReader("InterfaceImplementation", tosca_v1_3.ReadInterfaceImplementation)
	Grammar.RegisterReader("ParameterDefinition", tosca_v1_3.ReadParameterDefinition)
	Grammar.RegisterReader("Policy", tosca_v1_3.ReadPolicy)
	Grammar.RegisterReader("PolicyType", tosca_v1_3.ReadPolicyType)
	Grammar.RegisterReader("PropertyDefinition", ReadPropertyDefinition) // override
	Grammar.RegisterReader("PropertyFilter", tosca_v1_3.ReadPropertyFilter)
	Grammar.RegisterReader("PropertyMapping", tosca_v1_3.ReadPropertyMapping) // introduced in TOSCA 1.2
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
	Grammar.RegisterReader("SubstitutionMappings", ReadSubstitutionMappings) // override
	Grammar.RegisterReader("timestamp", tosca_v1_3.ReadTimestamp)
	Grammar.RegisterReader("TopologyTemplate", tosca_v1_3.ReadTopologyTemplate)
	Grammar.RegisterReader("TriggerDefinition", tosca_v1_3.ReadTriggerDefinition)
	Grammar.RegisterReader("TriggerDefinitionCondition", tosca_v1_3.ReadTriggerDefinitionCondition)
	Grammar.RegisterReader("Value", tosca_v1_3.ReadValue)
	Grammar.RegisterReader("version", tosca_v1_3.ReadVersion)
	Grammar.RegisterReader("WorkflowActivityDefinition", tosca_v1_3.ReadWorkflowActivityDefinition)         // introduced in TOSCA 1.1
	Grammar.RegisterReader("WorkflowDefinition", tosca_v1_3.ReadWorkflowDefinition)                         // introduced in TOSCA 1.1
	Grammar.RegisterReader("WorkflowPreconditionDefinition", tosca_v1_3.ReadWorkflowPreconditionDefinition) // introduced in TOSCA 1.1
	Grammar.RegisterReader("WorkflowStepDefinition", tosca_v1_3.ReadWorkflowStepDefinition)                 // introduced in TOSCA 1.1

	DefaultScriptletNamespace.RegisterScriptlets(tosca_v1_3.FunctionScriptlets, nil, "tosca.function.join")
	DefaultScriptletNamespace.RegisterScriptlets(tosca_v1_3.ConstraintClauseScriptlets, tosca_v1_3.ConstraintClauseNativeArgumentIndexes)
}
