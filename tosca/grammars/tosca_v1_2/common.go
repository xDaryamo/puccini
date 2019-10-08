package tosca_v1_2

import (
	"github.com/op/go-logging"
	"github.com/tliron/puccini/js"
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/grammars/tosca_v1_3"
)

var log = logging.MustGetLogger("grammars.tosca_v1_2")

var Grammar = make(tosca.Grammar)

var DefaultScriptletNamespace = make(tosca.ScriptletNamespace)

func init() {
	Grammar["Artifact"] = ReadArtifact                     // override
	Grammar["ArtifactDefinition"] = ReadArtifactDefinition // override
	Grammar["ArtifactType"] = tosca_v1_3.ReadArtifactType
	Grammar["AttributeDefinition"] = ReadAttributeDefinition // override
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
	Grammar["Group"] = ReadGroup         // override
	Grammar["GroupType"] = ReadGroupType // override
	Grammar["Import"] = tosca_v1_3.ReadImport
	Grammar["InterfaceAssignment"] = ReadInterfaceAssignment      // override
	Grammar["InterfaceDefinition"] = ReadInterfaceDefinition      // override
	Grammar["InterfaceMapping"] = tosca_v1_3.ReadInterfaceMapping // introduced in TOSCA 1.2
	Grammar["InterfaceType"] = ReadInterfaceType                  // override
	Grammar["Metadata"] = tosca_v1_3.ReadMetadata
	Grammar["NodeFilter"] = tosca_v1_3.ReadNodeFilter
	Grammar["NodeTemplate"] = tosca_v1_3.ReadNodeTemplate
	Grammar["NodeType"] = tosca_v1_3.ReadNodeType
	Grammar["NotificationDefinition"] = tosca_v1_3.ReadNotificationDefinition // not used
	Grammar["OperationAssignment"] = tosca_v1_3.ReadOperationAssignment
	Grammar["OperationDefinition"] = tosca_v1_3.ReadOperationDefinition
	Grammar["InterfaceImplementation"] = tosca_v1_3.ReadInterfaceImplementation
	Grammar["ParameterDefinition"] = tosca_v1_3.ReadParameterDefinition
	Grammar["Policy"] = tosca_v1_3.ReadPolicy
	Grammar["PolicyType"] = tosca_v1_3.ReadPolicyType
	Grammar["PropertyDefinition"] = ReadPropertyDefinition // override
	Grammar["PropertyFilter"] = tosca_v1_3.ReadPropertyFilter
	Grammar["PropertyMapping"] = tosca_v1_3.ReadPropertyMapping // introduced in TOSCA 1.2
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
	Grammar["ServiceTemplate"] = tosca_v1_3.ReadServiceTemplate
	Grammar["SubstitutionMappings"] = tosca_v1_3.ReadSubstitutionMappings
	Grammar["timestamp"] = tosca_v1_3.ReadTimestamp
	Grammar["TopologyTemplate"] = tosca_v1_3.ReadTopologyTemplate
	Grammar["TriggerDefinition"] = tosca_v1_3.ReadTriggerDefinition
	Grammar["TriggerDefinitionCondition"] = tosca_v1_3.ReadTriggerDefinitionCondition
	Grammar["Unit"] = tosca_v1_3.ReadUnit
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
