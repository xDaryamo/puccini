package tosca_v2_0

import (
	"github.com/tliron/kutil/logging"
	"github.com/tliron/puccini/tosca"
)

var log = logging.GetLogger("puccini.grammars.tosca_v2_0")
var logInherit = logging.NewSubLogger(log, "inherit")
var logRender = logging.NewSubLogger(log, "render")
var logNormalize = logging.NewSubLogger(log, "normalize")

var Grammar = tosca.NewGrammar()

var DefaultScriptletNamespace = tosca.NewScriptletNamespace()

func init() {
	Grammar.RegisterVersion("tosca_definitions_version", "tosca_2_0", "/tosca/implicit/2.0/profile.yaml")

	Grammar.RegisterReader("$Root", ReadServiceTemplate)
	Grammar.RegisterReader("$Unit", ReadUnit)

	Grammar.RegisterReader("Artifact", ReadArtifact)
	Grammar.RegisterReader("ArtifactDefinition", ReadArtifactDefinition)
	Grammar.RegisterReader("ArtifactType", ReadArtifactType)
	Grammar.RegisterReader("AttributeDefinition", ReadAttributeDefinition)
	Grammar.RegisterReader("AttributeMapping", ReadAttributeMapping) // introduced in TOSCA 1.3
	Grammar.RegisterReader("AttributeValue", ReadAttributeValue)
	Grammar.RegisterReader("bytes", ReadBytes) // introduced in TOSCA 2.0
	Grammar.RegisterReader("CapabilityAssignment", ReadCapabilityAssignment)
	Grammar.RegisterReader("CapabilityDefinition", ReadCapabilityDefinition)
	Grammar.RegisterReader("CapabilityFilter", ReadCapabilityFilter)
	Grammar.RegisterReader("CapabilityMapping", ReadCapabilityMapping)
	Grammar.RegisterReader("CapabilityType", ReadCapabilityType)
	Grammar.RegisterReader("ConditionClause", ReadConditionClause)
	Grammar.RegisterReader("ConstraintClause", ReadConstraintClause)
	Grammar.RegisterReader("DataType", ReadDataType)
	Grammar.RegisterReader("EventFilter", ReadEventFilter)
	Grammar.RegisterReader("Group", ReadGroup)
	Grammar.RegisterReader("GroupType", ReadGroupType)
	Grammar.RegisterReader("Import", ReadImport)
	Grammar.RegisterReader("InterfaceAssignment", ReadInterfaceAssignment)
	Grammar.RegisterReader("InterfaceDefinition", ReadInterfaceDefinition)
	Grammar.RegisterReader("InterfaceMapping", ReadInterfaceMapping) // introduced in TOSCA 1.2
	Grammar.RegisterReader("InterfaceType", ReadInterfaceType)
	Grammar.RegisterReader("Metadata", ReadMetadata)
	Grammar.RegisterReader("NodeFilter", ReadNodeFilter)
	Grammar.RegisterReader("NodeTemplate", ReadNodeTemplate)
	Grammar.RegisterReader("NodeType", ReadNodeType)
	Grammar.RegisterReader("NotificationAssignment", ReadNotificationAssignment) // introduced in TOSCA 1.3
	Grammar.RegisterReader("NotificationDefinition", ReadNotificationDefinition) // introduced in TOSCA 1.3
	Grammar.RegisterReader("OperationAssignment", ReadOperationAssignment)
	Grammar.RegisterReader("OperationDefinition", ReadOperationDefinition)
	Grammar.RegisterReader("OutputMapping", ReadOutputMapping) // introduced in TOSCA 1.3
	Grammar.RegisterReader("InterfaceImplementation", ReadInterfaceImplementation)
	Grammar.RegisterReader("ParameterDefinition", ReadParameterDefinition)
	Grammar.RegisterReader("Policy", ReadPolicy)
	Grammar.RegisterReader("PolicyType", ReadPolicyType)
	Grammar.RegisterReader("PropertyDefinition", ReadPropertyDefinition)
	Grammar.RegisterReader("PropertyFilter", ReadPropertyFilter)
	Grammar.RegisterReader("PropertyMapping", ReadPropertyMapping) // introduced in TOSCA 1.2
	Grammar.RegisterReader("range", ReadRange)
	Grammar.RegisterReader("RangeEntity", ReadRangeEntity)
	Grammar.RegisterReader("RelationshipAssignment", ReadRelationshipAssignment)
	Grammar.RegisterReader("RelationshipDefinition", ReadRelationshipDefinition)
	Grammar.RegisterReader("RelationshipTemplate", ReadRelationshipTemplate)
	Grammar.RegisterReader("RelationshipType", ReadRelationshipType)
	Grammar.RegisterReader("Repository", ReadRepository)
	Grammar.RegisterReader("RequirementAssignment", ReadRequirementAssignment)
	Grammar.RegisterReader("RequirementDefinition", ReadRequirementDefinition)
	Grammar.RegisterReader("RequirementMapping", ReadRequirementMapping)
	Grammar.RegisterReader("scalar-unit.bitrate", ReadScalarUnitBitrate) // introduced in TOSCA 1.3
	Grammar.RegisterReader("scalar-unit.frequency", ReadScalarUnitFrequency)
	Grammar.RegisterReader("scalar-unit.size", ReadScalarUnitSize)
	Grammar.RegisterReader("scalar-unit.time", ReadScalarUnitTime)
	Grammar.RegisterReader("Schema", ReadSchema)
	Grammar.RegisterReader("SubstitutionMappings", ReadSubstitutionMappings)
	Grammar.RegisterReader("timestamp", ReadTimestamp)
	Grammar.RegisterReader("TopologyTemplate", ReadTopologyTemplate)
	Grammar.RegisterReader("TriggerDefinition", ReadTriggerDefinition)
	Grammar.RegisterReader("TriggerDefinitionCondition", ReadTriggerDefinitionCondition)
	Grammar.RegisterReader("Value", ReadValue)
	Grammar.RegisterReader("version", ReadVersion)
	Grammar.RegisterReader("WorkflowActivityCallOperation", ReadWorkflowActivityCallOperation)   // TODO: introduced in TOSCA 1.3
	Grammar.RegisterReader("WorkflowActivityDefinition", ReadWorkflowActivityDefinition)         // introduced in TOSCA 1.1
	Grammar.RegisterReader("WorkflowDefinition", ReadWorkflowDefinition)                         // introduced in TOSCA 1.1
	Grammar.RegisterReader("WorkflowPreconditionDefinition", ReadWorkflowPreconditionDefinition) // introduced in TOSCA 1.1
	Grammar.RegisterReader("WorkflowStepDefinition", ReadWorkflowStepDefinition)                 // introduced in TOSCA 1.1

	DefaultScriptletNamespace.RegisterScriptlets(FunctionScriptlets, nil)
	DefaultScriptletNamespace.RegisterScriptlets(ConstraintClauseScriptlets, ConstraintClauseNativeArgumentIndexes)
}

func CompareUint32(v1 uint32, v2 uint32) int {
	if v1 < v2 {
		return -1
	} else if v2 > v1 {
		return 1
	}
	return 0
}

func CompareUint64(v1 uint64, v2 uint64) int {
	if v1 < v2 {
		return -1
	} else if v2 > v1 {
		return 1
	}
	return 0
}

func CompareInt64(v1 int64, v2 int64) int {
	if v1 < v2 {
		return -1
	} else if v2 > v1 {
		return 1
	}
	return 0
}

func CompareFloat64(v1 float64, v2 float64) int {
	if v1 < v2 {
		return -1
	} else if v2 > v1 {
		return 1
	}
	return 0
}
