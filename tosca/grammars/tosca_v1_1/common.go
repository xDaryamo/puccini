package tosca_v1_1

import (
	"github.com/op/go-logging"
	"github.com/tliron/puccini/js"
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/grammars/tosca_v1_2"
)

var log = logging.MustGetLogger("grammars.tosca_v1_1")

var Grammar = make(tosca.Grammar)

var DefaultScriptNamespace = make(tosca.ScriptNamespace)

func init() {
	Grammar["Artifact"] = tosca_v1_2.ReadArtifact
	Grammar["ArtifactDefinition"] = tosca_v1_2.ReadArtifactDefinition
	Grammar["ArtifactType"] = tosca_v1_2.ReadArtifactType
	Grammar["AttributeDefinition"] = tosca_v1_2.ReadAttributeDefinition
	Grammar["AttributeValue"] = tosca_v1_2.ReadAttributeValue
	Grammar["CapabilityAssignment"] = tosca_v1_2.ReadCapabilityAssignment
	Grammar["CapabilityDefinition"] = tosca_v1_2.ReadCapabilityDefinition
	Grammar["CapabilityFilter"] = tosca_v1_2.ReadCapabilityFilter
	Grammar["CapabilityMapping"] = tosca_v1_2.ReadCapabilityMapping
	Grammar["CapabilityType"] = tosca_v1_2.ReadCapabilityType
	Grammar["ConditionClause"] = tosca_v1_2.ReadConditionClause
	Grammar["ConstraintClause"] = tosca_v1_2.ReadConstraintClause
	Grammar["DataType"] = tosca_v1_2.ReadDataType
	Grammar["EntrySchema"] = tosca_v1_2.ReadEntrySchema
	Grammar["EventFilter"] = tosca_v1_2.ReadEventFilter
	Grammar["Group"] = tosca_v1_2.ReadGroup
	Grammar["GroupType"] = tosca_v1_2.ReadGroupType
	Grammar["Import"] = tosca_v1_2.ReadImport
	Grammar["InterfaceAssignment"] = tosca_v1_2.ReadInterfaceAssignment
	Grammar["InterfaceDefinition"] = tosca_v1_2.ReadInterfaceDefinition
	Grammar["InterfaceType"] = tosca_v1_2.ReadInterfaceType
	Grammar["Metadata"] = tosca_v1_2.ReadMetadata
	Grammar["NodeFilter"] = tosca_v1_2.ReadNodeFilter
	Grammar["NodeTemplate"] = tosca_v1_2.ReadNodeTemplate
	Grammar["NodeType"] = tosca_v1_2.ReadNodeType
	Grammar["OperationAssignment"] = tosca_v1_2.ReadOperationAssignment
	Grammar["OperationDefinition"] = tosca_v1_2.ReadOperationDefinition
	Grammar["OperationImplementation"] = ReadOperationImplementation // override
	Grammar["ParameterDefinition"] = tosca_v1_2.ReadParameterDefinition
	Grammar["Policy"] = tosca_v1_2.ReadPolicy
	Grammar["PolicyType"] = tosca_v1_2.ReadPolicyType
	Grammar["PropertyDefinition"] = tosca_v1_2.ReadPropertyDefinition
	Grammar["PropertyFilter"] = tosca_v1_2.ReadPropertyFilter
	Grammar["range"] = tosca_v1_2.ReadRange
	Grammar["RangeEntity"] = tosca_v1_2.ReadRangeEntity
	Grammar["RelationshipAssignment"] = tosca_v1_2.ReadRelationshipAssignment
	Grammar["RelationshipDefinition"] = tosca_v1_2.ReadRelationshipDefinition
	Grammar["RelationshipTemplate"] = tosca_v1_2.ReadRelationshipTemplate
	Grammar["RelationshipType"] = tosca_v1_2.ReadRelationshipType
	Grammar["Repository"] = tosca_v1_2.ReadRepository
	Grammar["RequirementAssignment"] = tosca_v1_2.ReadRequirementAssignment
	Grammar["RequirementDefinition"] = tosca_v1_2.ReadRequirementDefinition
	Grammar["RequirementMapping"] = tosca_v1_2.ReadRequirementMapping
	Grammar["scalar-unit.size"] = tosca_v1_2.ReadScalarUnitSize
	Grammar["scalar-unit.time"] = tosca_v1_2.ReadScalarUnitTime
	Grammar["scalar-unit.frequency"] = tosca_v1_2.ReadScalarUnitFrequency
	Grammar["ServiceTemplate"] = ReadServiceTemplate           // override
	Grammar["SubstitutionMappings"] = ReadSubstitutionMappings // override
	Grammar["timestamp"] = tosca_v1_2.ReadTimestamp
	Grammar["TopologyTemplate"] = tosca_v1_2.ReadTopologyTemplate
	Grammar["TriggerDefinition"] = tosca_v1_2.ReadTriggerDefinition
	Grammar["TriggerDefinitionCondition"] = tosca_v1_2.ReadTriggerDefinitionCondition
	Grammar["Unit"] = ReadUnit // override
	Grammar["Value"] = tosca_v1_2.ReadValue
	Grammar["version"] = tosca_v1_2.ReadVersion
	Grammar["WorkflowActivityDefinition"] = tosca_v1_2.ReadWorkflowActivityDefinition
	Grammar["WorkflowDefinition"] = tosca_v1_2.ReadWorkflowDefinition
	Grammar["WorkflowPreconditionDefinition"] = tosca_v1_2.ReadWorkflowPreconditionDefinition
	Grammar["WorkflowStepDefinition"] = tosca_v1_2.ReadWorkflowStepDefinition

	for name, sourceCode := range tosca_v1_2.FunctionSourceCode {
		// Unsupported functions
		if name == "join" {
			continue
		}

		DefaultScriptNamespace[name] = &tosca.Script{
			SourceCode: js.Cleanup(sourceCode),
		}
	}

	for name, sourceCode := range tosca_v1_2.ConstraintClauseSourceCode {
		nativeArgumentIndexes, _ := tosca_v1_2.ConstraintClauseNativeArgumentIndexes[name]
		DefaultScriptNamespace[name] = &tosca.Script{
			SourceCode:            js.Cleanup(sourceCode),
			NativeArgumentIndexes: nativeArgumentIndexes,
		}
	}
}
