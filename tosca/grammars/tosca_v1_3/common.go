package tosca_v1_3

import (
	"github.com/tliron/kutil/logging"
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/grammars/tosca_v2_0"
)

var log = logging.GetLogger("puccini.grammars.tosca_v1_3")

var Grammar = tosca.NewGrammar()

var DefaultScriptletNamespace = tosca.NewScriptletNamespace()

func init() {
	Grammar.RegisterVersion("tosca_definitions_version", "tosca_simple_yaml_1_3", "/tosca/simple/1.3/profile.yaml")

	Grammar.RegisterReader("$Root", ReadServiceTemplate)
	Grammar.RegisterReader("$Unit", ReadUnit)

	Grammar.RegisterReader("Artifact", tosca_v2_0.ReadArtifact)
	Grammar.RegisterReader("ArtifactDefinition", tosca_v2_0.ReadArtifactDefinition)
	Grammar.RegisterReader("ArtifactType", tosca_v2_0.ReadArtifactType)
	Grammar.RegisterReader("AttributeDefinition", tosca_v2_0.ReadAttributeDefinition)
	Grammar.RegisterReader("AttributeMapping", tosca_v2_0.ReadAttributeMapping) // introduced in TOSCA 1.3
	Grammar.RegisterReader("AttributeValue", tosca_v2_0.ReadAttributeValue)
	Grammar.RegisterReader("CapabilityAssignment", tosca_v2_0.ReadCapabilityAssignment)
	Grammar.RegisterReader("CapabilityDefinition", tosca_v2_0.ReadCapabilityDefinition)
	Grammar.RegisterReader("CapabilityFilter", tosca_v2_0.ReadCapabilityFilter)
	Grammar.RegisterReader("CapabilityMapping", tosca_v2_0.ReadCapabilityMapping)
	Grammar.RegisterReader("CapabilityType", tosca_v2_0.ReadCapabilityType)
	Grammar.RegisterReader("ConditionClause", tosca_v2_0.ReadConditionClause)
	Grammar.RegisterReader("ConstraintClause", tosca_v2_0.ReadConstraintClause)
	Grammar.RegisterReader("DataType", tosca_v2_0.ReadDataType)
	Grammar.RegisterReader("EventFilter", tosca_v2_0.ReadEventFilter)
	Grammar.RegisterReader("Group", tosca_v2_0.ReadGroup)
	Grammar.RegisterReader("GroupType", tosca_v2_0.ReadGroupType)
	Grammar.RegisterReader("Import", ReadImport) /// override
	Grammar.RegisterReader("InterfaceAssignment", tosca_v2_0.ReadInterfaceAssignment)
	Grammar.RegisterReader("InterfaceDefinition", tosca_v2_0.ReadInterfaceDefinition)
	Grammar.RegisterReader("InterfaceMapping", tosca_v2_0.ReadInterfaceMapping) // introduced in TOSCA 1.2
	Grammar.RegisterReader("InterfaceType", tosca_v2_0.ReadInterfaceType)
	Grammar.RegisterReader("Metadata", tosca_v2_0.ReadMetadata)
	Grammar.RegisterReader("NodeFilter", tosca_v2_0.ReadNodeFilter)
	Grammar.RegisterReader("NodeTemplate", tosca_v2_0.ReadNodeTemplate)
	Grammar.RegisterReader("NodeType", tosca_v2_0.ReadNodeType)
	Grammar.RegisterReader("NotificationAssignment", tosca_v2_0.ReadNotificationAssignment) // introduced in TOSCA 1.3
	Grammar.RegisterReader("NotificationDefinition", tosca_v2_0.ReadNotificationDefinition) // introduced in TOSCA 1.3
	Grammar.RegisterReader("OperationAssignment", tosca_v2_0.ReadOperationAssignment)
	Grammar.RegisterReader("OperationDefinition", tosca_v2_0.ReadOperationDefinition)
	Grammar.RegisterReader("OutputMapping", tosca_v2_0.ReadOutputMapping) // introduced in TOSCA 1.3
	Grammar.RegisterReader("InterfaceImplementation", tosca_v2_0.ReadInterfaceImplementation)
	Grammar.RegisterReader("ParameterDefinition", tosca_v2_0.ReadParameterDefinition)
	Grammar.RegisterReader("Policy", tosca_v2_0.ReadPolicy)
	Grammar.RegisterReader("PolicyType", tosca_v2_0.ReadPolicyType)
	Grammar.RegisterReader("PropertyDefinition", tosca_v2_0.ReadPropertyDefinition)
	Grammar.RegisterReader("PropertyFilter", tosca_v2_0.ReadPropertyFilter)
	Grammar.RegisterReader("PropertyMapping", tosca_v2_0.ReadPropertyMapping) // introduced in TOSCA 1.2
	Grammar.RegisterReader("range", tosca_v2_0.ReadRange)
	Grammar.RegisterReader("RangeEntity", tosca_v2_0.ReadRangeEntity)
	Grammar.RegisterReader("RelationshipAssignment", tosca_v2_0.ReadRelationshipAssignment)
	Grammar.RegisterReader("RelationshipDefinition", tosca_v2_0.ReadRelationshipDefinition)
	Grammar.RegisterReader("RelationshipTemplate", tosca_v2_0.ReadRelationshipTemplate)
	Grammar.RegisterReader("RelationshipType", tosca_v2_0.ReadRelationshipType)
	Grammar.RegisterReader("Repository", tosca_v2_0.ReadRepository)
	Grammar.RegisterReader("RequirementAssignment", tosca_v2_0.ReadRequirementAssignment)
	Grammar.RegisterReader("RequirementDefinition", tosca_v2_0.ReadRequirementDefinition)
	Grammar.RegisterReader("RequirementMapping", tosca_v2_0.ReadRequirementMapping)
	Grammar.RegisterReader("scalar-unit.bitrate", tosca_v2_0.ReadScalarUnitBitrate) // introduced in TOSCA 1.3
	Grammar.RegisterReader("scalar-unit.frequency", tosca_v2_0.ReadScalarUnitFrequency)
	Grammar.RegisterReader("scalar-unit.size", tosca_v2_0.ReadScalarUnitSize)
	Grammar.RegisterReader("scalar-unit.time", tosca_v2_0.ReadScalarUnitTime)
	Grammar.RegisterReader("Schema", tosca_v2_0.ReadSchema)
	Grammar.RegisterReader("SubstitutionMappings", tosca_v2_0.ReadSubstitutionMappings)
	Grammar.RegisterReader("timestamp", tosca_v2_0.ReadTimestamp)
	Grammar.RegisterReader("TopologyTemplate", tosca_v2_0.ReadTopologyTemplate)
	Grammar.RegisterReader("TriggerDefinition", tosca_v2_0.ReadTriggerDefinition)
	Grammar.RegisterReader("TriggerDefinitionCondition", tosca_v2_0.ReadTriggerDefinitionCondition)
	Grammar.RegisterReader("Value", tosca_v2_0.ReadValue)
	Grammar.RegisterReader("version", tosca_v2_0.ReadVersion)
	Grammar.RegisterReader("WorkflowActivityCallOperation", tosca_v2_0.ReadWorkflowActivityCallOperation)   // introduced in TOSCA 1.3
	Grammar.RegisterReader("WorkflowActivityDefinition", tosca_v2_0.ReadWorkflowActivityDefinition)         // introduced in TOSCA 1.1
	Grammar.RegisterReader("WorkflowDefinition", tosca_v2_0.ReadWorkflowDefinition)                         // introduced in TOSCA 1.1
	Grammar.RegisterReader("WorkflowPreconditionDefinition", tosca_v2_0.ReadWorkflowPreconditionDefinition) // introduced in TOSCA 1.1
	Grammar.RegisterReader("WorkflowStepDefinition", tosca_v2_0.ReadWorkflowStepDefinition)                 // introduced in TOSCA 1.1

	DefaultScriptletNamespace.RegisterScriptlets(tosca_v2_0.FunctionScriptlets, nil)
	DefaultScriptletNamespace.RegisterScriptlets(tosca_v2_0.ConstraintClauseScriptlets, tosca_v2_0.ConstraintClauseNativeArgumentIndexes)
}
