package tosca_v1_1

import (
	"github.com/op/go-logging"
	"github.com/tliron/puccini/js"
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/grammars/tosca_v1_2"
	"github.com/tliron/puccini/tosca/grammars/tosca_v1_3"
)

var log = logging.MustGetLogger("grammars.tosca_v1_1")

var Grammar = make(tosca.Grammar)

var DefaultScriptletNamespace = make(tosca.ScriptletNamespace)

func init() {
	Grammar["Artifact"] = tosca_v1_2.ReadArtifact                     // 1.2
	Grammar["ArtifactDefinition"] = tosca_v1_2.ReadArtifactDefinition // 1.2
	Grammar["ArtifactType"] = tosca_v1_3.ReadArtifactType
	Grammar["AttributeDefinition"] = tosca_v1_2.ReadAttributeDefinition // 1.2
	Grammar["AttributeValue"] = tosca_v1_3.ReadAttributeValue
	Grammar["CapabilityAssignment"] = tosca_v1_3.ReadCapabilityAssignment
	Grammar["CapabilityDefinition"] = tosca_v1_3.ReadCapabilityDefinition
	Grammar["CapabilityFilter"] = tosca_v1_3.ReadCapabilityFilter
	Grammar["CapabilityMapping"] = tosca_v1_3.ReadCapabilityMapping
	Grammar["CapabilityType"] = tosca_v1_3.ReadCapabilityType
	Grammar["ConditionClause"] = tosca_v1_3.ReadConditionClause
	Grammar["ConstraintClause"] = tosca_v1_3.ReadConstraintClause
	Grammar["DataType"] = tosca_v1_3.ReadDataType
	Grammar["EventFilter"] = tosca_v1_3.ReadEventFilter
	Grammar["Group"] = tosca_v1_2.ReadGroup         // 1.2
	Grammar["GroupType"] = tosca_v1_2.ReadGroupType // 1.2
	Grammar["Import"] = tosca_v1_3.ReadImport
	Grammar["InterfaceAssignment"] = tosca_v1_2.ReadInterfaceAssignment // 1.2
	Grammar["InterfaceDefinition"] = tosca_v1_2.ReadInterfaceDefinition // 1.2
	Grammar["InterfaceType"] = tosca_v1_2.ReadInterfaceType             // 1.2
	Grammar["Metadata"] = tosca_v1_3.ReadMetadata
	Grammar["NodeFilter"] = tosca_v1_3.ReadNodeFilter
	Grammar["NodeTemplate"] = tosca_v1_3.ReadNodeTemplate
	Grammar["NodeType"] = tosca_v1_3.ReadNodeType
	Grammar["NotificationDefinition"] = tosca_v1_3.ReadNotificationDefinition // unused
	Grammar["OperationAssignment"] = tosca_v1_3.ReadOperationAssignment
	Grammar["OperationDefinition"] = tosca_v1_3.ReadOperationDefinition
	Grammar["InterfaceImplementation"] = ReadInterfaceImplementation // override
	Grammar["ParameterDefinition"] = tosca_v1_3.ReadParameterDefinition
	Grammar["Policy"] = tosca_v1_3.ReadPolicy
	Grammar["PolicyType"] = tosca_v1_3.ReadPolicyType
	Grammar["PropertyDefinition"] = tosca_v1_2.ReadPropertyDefinition // 1.2
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
	Grammar["scalar-unit.frequency"] = tosca_v1_3.ReadScalarUnitFrequency
	Grammar["scalar-unit.size"] = tosca_v1_3.ReadScalarUnitSize
	Grammar["scalar-unit.time"] = tosca_v1_3.ReadScalarUnitTime
	Grammar["Schema"] = tosca_v1_3.ReadSchema
	Grammar["ServiceTemplate"] = ReadServiceTemplate           // override
	Grammar["SubstitutionMappings"] = ReadSubstitutionMappings // override
	Grammar["timestamp"] = tosca_v1_3.ReadTimestamp
	Grammar["TopologyTemplate"] = tosca_v1_3.ReadTopologyTemplate
	Grammar["TriggerDefinition"] = tosca_v1_3.ReadTriggerDefinition
	Grammar["TriggerDefinitionCondition"] = tosca_v1_3.ReadTriggerDefinitionCondition
	Grammar["Unit"] = ReadUnit // override
	Grammar["Value"] = tosca_v1_3.ReadValue
	Grammar["version"] = tosca_v1_3.ReadVersion
	Grammar["WorkflowActivityDefinition"] = tosca_v1_3.ReadWorkflowActivityDefinition
	Grammar["WorkflowDefinition"] = tosca_v1_3.ReadWorkflowDefinition
	Grammar["WorkflowPreconditionDefinition"] = tosca_v1_3.ReadWorkflowPreconditionDefinition
	Grammar["WorkflowStepDefinition"] = tosca_v1_3.ReadWorkflowStepDefinition

	for name, scriptlet := range tosca_v1_3.FunctionScriptlets {
		// Unsupported functions
		if name == "join" {
			continue
		}

		DefaultScriptletNamespace[name] = &tosca.Scriptlet{
			Scriptlet: js.CleanupScriptlet(scriptlet),
		}
	}

	for name, scriptlet := range tosca_v1_3.ConstraintClauseScriptlets {
		// Unsupported constraints
		if name == "schema" {
			continue
		}

		nativeArgumentIndexes, _ := tosca_v1_3.ConstraintClauseNativeArgumentIndexes[name]
		DefaultScriptletNamespace[name] = &tosca.Scriptlet{
			Scriptlet:             js.CleanupScriptlet(scriptlet),
			NativeArgumentIndexes: nativeArgumentIndexes,
		}
	}
}
