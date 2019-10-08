package tosca_v1_3

import (
	"github.com/op/go-logging"
	"github.com/tliron/puccini/js"
	"github.com/tliron/puccini/tosca"
)

var log = logging.MustGetLogger("grammars.tosca_v1_3")

var Grammar = make(tosca.Grammar)

var DefaultScriptletNamespace = make(tosca.ScriptletNamespace)

func init() {
	Grammar["Artifact"] = ReadArtifact
	Grammar["ArtifactDefinition"] = ReadArtifactDefinition
	Grammar["ArtifactType"] = ReadArtifactType
	Grammar["AttributeDefinition"] = ReadAttributeDefinition
	Grammar["AttributeMapping"] = ReadAttributeMapping // introduced in TOSCA 1.3
	Grammar["AttributeValue"] = ReadAttributeValue
	Grammar["CapabilityAssignment"] = ReadCapabilityAssignment
	Grammar["CapabilityDefinition"] = ReadCapabilityDefinition
	Grammar["CapabilityFilter"] = ReadCapabilityFilter
	Grammar["CapabilityMapping"] = ReadCapabilityMapping
	Grammar["CapabilityType"] = ReadCapabilityType
	Grammar["ConditionClause"] = ReadConditionClause
	Grammar["ConstraintClause"] = ReadConstraintClause
	Grammar["DataType"] = ReadDataType
	Grammar["EventFilter"] = ReadEventFilter
	Grammar["Group"] = ReadGroup
	Grammar["GroupType"] = ReadGroupType
	Grammar["Import"] = ReadImport
	Grammar["InterfaceAssignment"] = ReadInterfaceAssignment
	Grammar["InterfaceDefinition"] = ReadInterfaceDefinition
	Grammar["InterfaceMapping"] = ReadInterfaceMapping // introduced in TOSCA 1.2
	Grammar["InterfaceType"] = ReadInterfaceType
	Grammar["Metadata"] = ReadMetadata
	Grammar["NodeFilter"] = ReadNodeFilter
	Grammar["NodeTemplate"] = ReadNodeTemplate
	Grammar["NodeType"] = ReadNodeType
	Grammar["NotificationAssignment"] = ReadNotificationAssignment // introduced in TOSCA 1.3
	Grammar["NotificationDefinition"] = ReadNotificationDefinition // introduced in TOSCA 1.3
	Grammar["OperationAssignment"] = ReadOperationAssignment
	Grammar["OperationDefinition"] = ReadOperationDefinition
	Grammar["InterfaceImplementation"] = ReadInterfaceImplementation
	Grammar["ParameterDefinition"] = ReadParameterDefinition
	Grammar["Policy"] = ReadPolicy
	Grammar["PolicyType"] = ReadPolicyType
	Grammar["PropertyDefinition"] = ReadPropertyDefinition
	Grammar["PropertyFilter"] = ReadPropertyFilter
	Grammar["PropertyMapping"] = ReadPropertyMapping // introduced in TOSCA 1.2
	Grammar["range"] = ReadRange
	Grammar["RangeEntity"] = ReadRangeEntity
	Grammar["RelationshipAssignment"] = ReadRelationshipAssignment
	Grammar["RelationshipDefinition"] = ReadRelationshipDefinition
	Grammar["RelationshipTemplate"] = ReadRelationshipTemplate
	Grammar["RelationshipType"] = ReadRelationshipType
	Grammar["Repository"] = ReadRepository
	Grammar["RequirementAssignment"] = ReadRequirementAssignment
	Grammar["RequirementDefinition"] = ReadRequirementDefinition
	Grammar["RequirementMapping"] = ReadRequirementMapping
	Grammar["scalar-unit.bitrate"] = ReadScalarUnitBitrate // introduced in TOSCA 1.3
	Grammar["scalar-unit.frequency"] = ReadScalarUnitFrequency
	Grammar["scalar-unit.size"] = ReadScalarUnitSize
	Grammar["scalar-unit.time"] = ReadScalarUnitTime
	Grammar["Schema"] = ReadSchema
	Grammar["ServiceTemplate"] = ReadServiceTemplate
	Grammar["SubstitutionMappings"] = ReadSubstitutionMappings
	Grammar["timestamp"] = ReadTimestamp
	Grammar["TopologyTemplate"] = ReadTopologyTemplate
	Grammar["TriggerDefinition"] = ReadTriggerDefinition
	Grammar["TriggerDefinitionCondition"] = ReadTriggerDefinitionCondition
	Grammar["Unit"] = ReadUnit
	Grammar["Value"] = ReadValue
	Grammar["version"] = ReadVersion
	Grammar["WorkflowActivityCallOperation"] = ReadWorkflowActivityCallOperation // introduced in TOSCA 1.3
	Grammar["WorkflowActivityDefinition"] = ReadWorkflowActivityDefinition
	Grammar["WorkflowDefinition"] = ReadWorkflowDefinition
	Grammar["WorkflowPreconditionDefinition"] = ReadWorkflowPreconditionDefinition
	Grammar["WorkflowStepDefinition"] = ReadWorkflowStepDefinition

	for name, scriptlet := range FunctionScriptlets {
		DefaultScriptletNamespace[name] = &tosca.Scriptlet{
			Scriptlet: js.CleanupScriptlet(scriptlet),
		}
	}

	for name, scriptlet := range ConstraintClauseScriptlets {
		nativeArgumentIndexes, _ := ConstraintClauseNativeArgumentIndexes[name]
		DefaultScriptletNamespace[name] = &tosca.Scriptlet{
			Scriptlet:             js.CleanupScriptlet(scriptlet),
			NativeArgumentIndexes: nativeArgumentIndexes,
		}
	}
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
