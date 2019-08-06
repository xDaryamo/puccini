package tosca_v1_3

import (
	"github.com/op/go-logging"
	"github.com/tliron/puccini/js"
	"github.com/tliron/puccini/tosca"
)

var log = logging.MustGetLogger("grammars.tosca_v1_3")

var Grammar = make(tosca.Grammar)

var DefaultScriptNamespace = make(tosca.ScriptNamespace)

func init() {
	Grammar["Artifact"] = ReadArtifact
	Grammar["ArtifactDefinition"] = ReadArtifactDefinition
	Grammar["ArtifactType"] = ReadArtifactType
	Grammar["AttributeDefinition"] = ReadAttributeDefinition
	Grammar["AttributeValue"] = ReadAttributeValue
	Grammar["CapabilityAssignment"] = ReadCapabilityAssignment
	Grammar["CapabilityDefinition"] = ReadCapabilityDefinition
	Grammar["CapabilityFilter"] = ReadCapabilityFilter
	Grammar["CapabilityMapping"] = ReadCapabilityMapping
	Grammar["CapabilityType"] = ReadCapabilityType
	Grammar["ConditionClause"] = ReadConditionClause
	Grammar["ConstraintClause"] = ReadConstraintClause
	Grammar["DataType"] = ReadDataType
	Grammar["EntrySchema"] = ReadEntrySchema
	Grammar["EventFilter"] = ReadEventFilter
	Grammar["Group"] = ReadGroup
	Grammar["GroupType"] = ReadGroupType
	Grammar["Import"] = ReadImport
	Grammar["InterfaceAssignment"] = ReadInterfaceAssignment
	Grammar["InterfaceDefinition"] = ReadInterfaceDefinition
	Grammar["InterfaceMapping"] = ReadInterfaceMapping // introduced in 1.2
	Grammar["InterfaceType"] = ReadInterfaceType
	Grammar["Metadata"] = ReadMetadata
	Grammar["NodeFilter"] = ReadNodeFilter
	Grammar["NodeTemplate"] = ReadNodeTemplate
	Grammar["NodeType"] = ReadNodeType
	Grammar["OperationAssignment"] = ReadOperationAssignment
	Grammar["OperationDefinition"] = ReadOperationDefinition
	Grammar["OperationImplementation"] = ReadOperationImplementation
	Grammar["ParameterDefinition"] = ReadParameterDefinition
	Grammar["Policy"] = ReadPolicy
	Grammar["PolicyType"] = ReadPolicyType
	Grammar["PropertyDefinition"] = ReadPropertyDefinition
	Grammar["PropertyFilter"] = ReadPropertyFilter
	Grammar["PropertyMapping"] = ReadPropertyMapping // introduced in 1.2
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
	Grammar["scalar-unit.size"] = ReadScalarUnitSize
	Grammar["scalar-unit.time"] = ReadScalarUnitTime
	Grammar["scalar-unit.frequency"] = ReadScalarUnitFrequency
	Grammar["ServiceTemplate"] = ReadServiceTemplate
	Grammar["SubstitutionMappings"] = ReadSubstitutionMappings
	Grammar["timestamp"] = ReadTimestamp
	Grammar["TopologyTemplate"] = ReadTopologyTemplate
	Grammar["TriggerDefinition"] = ReadTriggerDefinition
	Grammar["TriggerDefinitionCondition"] = ReadTriggerDefinitionCondition
	Grammar["Unit"] = ReadUnit
	Grammar["Value"] = ReadValue
	Grammar["version"] = ReadVersion
	Grammar["WorkflowActivityCallOperation"] = ReadWorkflowActivityCallOperation // introduced in 1.3
	Grammar["WorkflowActivityDefinition"] = ReadWorkflowActivityDefinition
	Grammar["WorkflowDefinition"] = ReadWorkflowDefinition
	Grammar["WorkflowPreconditionDefinition"] = ReadWorkflowPreconditionDefinition
	Grammar["WorkflowStepDefinition"] = ReadWorkflowStepDefinition

	for name, sourceCode := range FunctionSourceCode {
		DefaultScriptNamespace[name] = &tosca.Script{
			SourceCode: js.Cleanup(sourceCode),
		}
	}

	for name, sourceCode := range ConstraintClauseSourceCode {
		nativeArgumentIndexes, _ := ConstraintClauseNativeArgumentIndexes[name]
		DefaultScriptNamespace[name] = &tosca.Script{
			SourceCode:            js.Cleanup(sourceCode),
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
