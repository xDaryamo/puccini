package v1_1

import (
	"github.com/op/go-logging"
	"github.com/tliron/puccini/js"
	"github.com/tliron/puccini/tosca"
)

var log = logging.MustGetLogger("grammars.v1_1")

var Readers = make(map[string]tosca.Reader)

var DefaultScriptNamespace = make(tosca.ScriptNamespace)

func init() {
	Readers["Artifact"] = ReadArtifact
	Readers["ArtifactDefinition"] = ReadArtifactDefinition
	Readers["ArtifactType"] = ReadArtifactType
	Readers["AttributeDefinition"] = ReadAttributeDefinition
	Readers["CapabilityAssignment"] = ReadCapabilityAssignment
	Readers["CapabilityDefinition"] = ReadCapabilityDefinition
	Readers["CapabilityFilter"] = ReadCapabilityFilter
	Readers["CapabilityMapping"] = ReadCapabilityMapping
	Readers["CapabilityType"] = ReadCapabilityType
	Readers["ConditionClause"] = ReadConditionClause
	Readers["ConstraintClause"] = ReadConstraintClause
	Readers["DataType"] = ReadDataType
	Readers["EntrySchema"] = ReadEntrySchema
	Readers["EventFilter"] = ReadEventFilter
	Readers["Group"] = ReadGroup
	Readers["GroupType"] = ReadGroupType
	Readers["Import"] = ReadImport
	Readers["InterfaceAssignment"] = ReadInterfaceAssignment
	Readers["InterfaceDefinition"] = ReadInterfaceDefinition
	Readers["InterfaceType"] = ReadInterfaceType
	Readers["Metadata"] = ReadMetadata
	Readers["NodeFilter"] = ReadNodeFilter
	Readers["NodeTemplate"] = ReadNodeTemplate
	Readers["NodeType"] = ReadNodeType
	Readers["OperationAssignment"] = ReadOperationAssignment
	Readers["OperationDefinition"] = ReadOperationDefinition
	Readers["OperationImplementation"] = ReadOperationImplementation
	Readers["ParameterDefinition"] = ReadParameterDefinition
	Readers["Policy"] = ReadPolicy
	Readers["PolicyType"] = ReadPolicyType
	Readers["PropertyDefinition"] = ReadPropertyDefinition
	Readers["range"] = ReadRange
	Readers["RangeEntity"] = ReadRangeEntity
	Readers["RelationshipAssignment"] = ReadRelationshipAssignment
	Readers["RelationshipDefinition"] = ReadRelationshipDefinition
	Readers["RelationshipTemplate"] = ReadRelationshipTemplate
	Readers["RelationshipType"] = ReadRelationshipType
	Readers["Repository"] = ReadRepository
	Readers["RequirementAssignment"] = ReadRequirementAssignment
	Readers["RequirementDefinition"] = ReadRequirementDefinition
	Readers["RequirementMapping"] = ReadRequirementMapping
	Readers["scalar-unit.size"] = ReadScalarUnitSize
	Readers["scalar-unit.time"] = ReadScalarUnitTime
	Readers["scalar-unit.frequency"] = ReadScalarUnitFrequency
	Readers["ServiceTemplate"] = ReadServiceTemplate
	Readers["SubstitutionMappings"] = ReadSubstitutionMappings
	Readers["timestamp"] = ReadTimestamp
	Readers["TopologyTemplate"] = ReadTopologyTemplate
	Readers["TriggerDefinition"] = ReadTriggerDefinition
	Readers["TriggerDefinitionCondition"] = ReadTriggerDefinitionCondition
	Readers["Unit"] = ReadUnit
	Readers["Value"] = ReadValue
	Readers["version"] = ReadVersion
	Readers["WorkflowActivityDefinition"] = ReadWorkflowActivityDefinition
	Readers["WorkflowDefinition"] = ReadWorkflowDefinition
	Readers["WorkflowPreconditionDefinition"] = ReadWorkflowPreconditionDefinition
	Readers["WorkflowStepDefinition"] = ReadWorkflowStepDefinition

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
