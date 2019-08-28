package tosca_v1_0

// TODO: This is currently a copy of the v1_1 grammar,
// but we should properly remove unsupported features

import (
	"github.com/op/go-logging"
	"github.com/tliron/puccini/js"
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/grammars/tosca_v1_1"
	"github.com/tliron/puccini/tosca/grammars/tosca_v1_2"
	"github.com/tliron/puccini/tosca/grammars/tosca_v1_3"
)

var log = logging.MustGetLogger("grammars.tosca_v1_0")

var Grammar = make(tosca.Grammar)

var DefaultScriptNamespace = make(tosca.ScriptNamespace)

func init() {
	Grammar["Artifact"] = tosca_v1_3.ReadArtifact
	Grammar["ArtifactDefinition"] = tosca_v1_3.ReadArtifactDefinition
	Grammar["ArtifactType"] = tosca_v1_3.ReadArtifactType
	Grammar["AttributeDefinition"] = tosca_v1_3.ReadAttributeDefinition
	Grammar["AttributeValue"] = tosca_v1_3.ReadAttributeValue
	Grammar["CapabilityAssignment"] = tosca_v1_3.ReadCapabilityAssignment
	Grammar["CapabilityDefinition"] = tosca_v1_3.ReadCapabilityDefinition
	Grammar["CapabilityFilter"] = tosca_v1_3.ReadCapabilityFilter
	Grammar["CapabilityMapping"] = tosca_v1_3.ReadCapabilityMapping
	Grammar["CapabilityType"] = tosca_v1_3.ReadCapabilityType
	Grammar["ConditionClause"] = tosca_v1_3.ReadConditionClause
	Grammar["ConstraintClause"] = tosca_v1_3.ReadConstraintClause
	Grammar["DataType"] = tosca_v1_3.ReadDataType
	Grammar["EntrySchema"] = tosca_v1_3.ReadEntrySchema
	Grammar["EventFilter"] = tosca_v1_3.ReadEventFilter
	Grammar["Group"] = tosca_v1_3.ReadGroup
	Grammar["GroupType"] = tosca_v1_3.ReadGroupType
	Grammar["Import"] = tosca_v1_3.ReadImport
	Grammar["InterfaceAssignment"] = tosca_v1_2.ReadInterfaceAssignment // 1.2
	Grammar["InterfaceDefinition"] = tosca_v1_2.ReadInterfaceDefinition // 1.2
	Grammar["InterfaceType"] = tosca_v1_2.ReadInterfaceType             // 1.2
	Grammar["Metadata"] = tosca_v1_3.ReadMetadata
	Grammar["NodeFilter"] = tosca_v1_3.ReadNodeFilter
	Grammar["NodeTemplate"] = tosca_v1_3.ReadNodeTemplate
	Grammar["NodeType"] = tosca_v1_3.ReadNodeType
	Grammar["NotificationDefinition"] = tosca_v1_3.ReadNotificationDefinition // not used
	Grammar["OperationAssignment"] = tosca_v1_3.ReadOperationAssignment
	Grammar["OperationDefinition"] = tosca_v1_3.ReadOperationDefinition
	Grammar["InterfaceImplementation"] = tosca_v1_1.ReadInterfaceImplementation // 1.1
	Grammar["ParameterDefinition"] = tosca_v1_3.ReadParameterDefinition
	Grammar["Policy"] = tosca_v1_3.ReadPolicy
	Grammar["PolicyType"] = tosca_v1_3.ReadPolicyType
	Grammar["PropertyDefinition"] = tosca_v1_3.ReadPropertyDefinition
	Grammar["PropertyFilter"] = tosca_v1_3.ReadPropertyFilter
	Grammar["range"] = tosca_v1_3.ReadRange
	Grammar["RangeEntity"] = tosca_v1_3.ReadRangeEntity
	Grammar["RelationshipAssignment"] = tosca_v1_3.ReadRelationshipAssignment
	Grammar["RelationshipDefinition"] = tosca_v1_3.ReadRelationshipDefinition
	Grammar["RelationshipTemplate"] = tosca_v1_3.ReadRelationshipTemplate
	Grammar["RelationshipType"] = tosca_v1_3.ReadRelationshipType
	Grammar["Repository"] = tosca_v1_3.ReadRepository
	Grammar["RequirementAssignment"] = tosca_v1_3.ReadRequirementAssignment
	Grammar["RequirementDefinition"] = tosca_v1_3.ReadRequirementDefinition
	Grammar["RequirementMapping"] = tosca_v1_3.ReadRequirementMapping
	Grammar["scalar-unit.size"] = tosca_v1_3.ReadScalarUnitSize
	Grammar["scalar-unit.time"] = tosca_v1_3.ReadScalarUnitTime
	Grammar["scalar-unit.frequency"] = tosca_v1_3.ReadScalarUnitFrequency
	Grammar["ServiceTemplate"] = tosca_v1_1.ReadServiceTemplate           // 1.1
	Grammar["SubstitutionMappings"] = tosca_v1_1.ReadSubstitutionMappings // 1.1
	Grammar["timestamp"] = tosca_v1_3.ReadTimestamp
	Grammar["TopologyTemplate"] = tosca_v1_3.ReadTopologyTemplate
	Grammar["TriggerDefinition"] = tosca_v1_3.ReadTriggerDefinition
	Grammar["TriggerDefinitionCondition"] = tosca_v1_3.ReadTriggerDefinitionCondition
	Grammar["Unit"] = tosca_v1_1.ReadUnit // 1.1
	Grammar["Value"] = tosca_v1_3.ReadValue
	Grammar["version"] = tosca_v1_3.ReadVersion
	Grammar["WorkflowActivityDefinition"] = tosca_v1_3.ReadWorkflowActivityDefinition
	Grammar["WorkflowDefinition"] = tosca_v1_3.ReadWorkflowDefinition
	Grammar["WorkflowPreconditionDefinition"] = tosca_v1_3.ReadWorkflowPreconditionDefinition
	Grammar["WorkflowStepDefinition"] = tosca_v1_3.ReadWorkflowStepDefinition

	for name, sourceCode := range tosca_v1_3.FunctionSourceCode {
		// Unsupported functions
		if name == "join" {
			continue
		}

		DefaultScriptNamespace[name] = &tosca.Script{
			SourceCode: js.Cleanup(sourceCode),
		}
	}

	for name, sourceCode := range tosca_v1_3.ConstraintClauseSourceCode {
		nativeArgumentIndexes, _ := tosca_v1_3.ConstraintClauseNativeArgumentIndexes[name]
		DefaultScriptNamespace[name] = &tosca.Script{
			SourceCode:            js.Cleanup(sourceCode),
			NativeArgumentIndexes: nativeArgumentIndexes,
		}
	}
}
