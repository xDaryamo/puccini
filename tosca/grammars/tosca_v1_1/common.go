package tosca_v1_1

import (
	"github.com/tliron/commonlog"
	"github.com/tliron/puccini/tosca/grammars/tosca_v1_2"
	"github.com/tliron/puccini/tosca/grammars/tosca_v1_3"
	"github.com/tliron/puccini/tosca/grammars/tosca_v2_0"
	"github.com/tliron/puccini/tosca/parsing"
)

var log = commonlog.GetLogger("puccini.grammars.tosca_v1_1")

var Grammar = parsing.NewGrammar()

var DefaultScriptletNamespace = parsing.NewScriptletNamespace()

func init() {
	Grammar.RegisterVersion("tosca_definitions_version", "tosca_simple_yaml_1_1", "/profiles/simple/1.1/profile.yaml")

	Grammar.RegisterReader("$Root", ReadServiceFile) // override
	Grammar.RegisterReader("$File", ReadFile)        // override

	Grammar.RegisterReader("Artifact", tosca_v1_2.ReadArtifact)                     // 1.2
	Grammar.RegisterReader("ArtifactDefinition", tosca_v1_2.ReadArtifactDefinition) // 1.2
	Grammar.RegisterReader("ArtifactType", tosca_v2_0.ReadArtifactType)
	Grammar.RegisterReader("AttributeDefinition", tosca_v1_2.ReadAttributeDefinition)   // 1.2
	Grammar.RegisterReader("AttributeValue", tosca_v1_3.ReadAttributeValue)             // 1.3
	Grammar.RegisterReader("CapabilityAssignment", tosca_v1_2.ReadCapabilityAssignment) // 1.2
	Grammar.RegisterReader("CapabilityDefinition", tosca_v2_0.ReadCapabilityDefinition)
	Grammar.RegisterReader("CapabilityFilter", tosca_v2_0.ReadCapabilityFilter)
	Grammar.RegisterReader("CapabilityMapping", tosca_v2_0.ReadCapabilityMapping)
	Grammar.RegisterReader("CapabilityType", tosca_v2_0.ReadCapabilityType)
	Grammar.RegisterReader("ConditionClause", tosca_v2_0.ReadConditionClause)
	Grammar.RegisterReader("ConditionClauseAnd", tosca_v2_0.ReadConditionClauseAnd)
	Grammar.RegisterReader("ConstraintClause", tosca_v1_3.ReadConstraintClause)
	Grammar.RegisterReader("DataType", tosca_v1_2.ReadDataType) // 1.2
	Grammar.RegisterReader("EventFilter", tosca_v2_0.ReadEventFilter)
	Grammar.RegisterReader("Group", tosca_v1_2.ReadGroup)                             // 1.2
	Grammar.RegisterReader("GroupType", tosca_v1_2.ReadGroupType)                     // 1.2
	Grammar.RegisterReader("Import", tosca_v1_3.ReadImport)                           /// 1.3
	Grammar.RegisterReader("InterfaceAssignment", tosca_v1_2.ReadInterfaceAssignment) // 1.2
	Grammar.RegisterReader("InterfaceDefinition", tosca_v1_2.ReadInterfaceDefinition) // 1.2
	Grammar.RegisterReader("InterfaceType", tosca_v1_2.ReadInterfaceType)             // 1.2
	Grammar.RegisterReader("Metadata", tosca_v2_0.ReadMetadata)
	Grammar.RegisterReader("NodeFilter", tosca_v2_0.ReadNodeFilter)
	Grammar.RegisterReader("NodeTemplate", tosca_v1_3.ReadNodeTemplate) // 1.3
	Grammar.RegisterReader("NodeType", tosca_v2_0.ReadNodeType)
	Grammar.RegisterReader("OperationAssignment", tosca_v1_2.ReadOperationAssignment) // 1.2
	Grammar.RegisterReader("OperationDefinition", tosca_v1_2.ReadOperationDefinition) // 1.2
	Grammar.RegisterReader("InterfaceImplementation", ReadInterfaceImplementation)    // override
	Grammar.RegisterReader("ParameterDefinition", tosca_v2_0.ReadParameterDefinition)
	Grammar.RegisterReader("Policy", tosca_v2_0.ReadPolicy)
	Grammar.RegisterReader("PolicyType", tosca_v2_0.ReadPolicyType)
	Grammar.RegisterReader("PropertyDefinition", tosca_v1_2.ReadPropertyDefinition) // 1.2
	Grammar.RegisterReader("PropertyFilter", tosca_v2_0.ReadPropertyFilter)
	Grammar.RegisterReader("range", tosca_v2_0.ReadRange)
	Grammar.RegisterReader("RangeEntity", tosca_v2_0.ReadRangeEntity)
	Grammar.RegisterReader("RelationshipAssignment", tosca_v2_0.ReadRelationshipAssignment)
	Grammar.RegisterReader("RelationshipDefinition", tosca_v2_0.ReadRelationshipDefinition)
	Grammar.RegisterReader("RelationshipTemplate", tosca_v2_0.ReadRelationshipTemplate)
	Grammar.RegisterReader("RelationshipType", tosca_v2_0.ReadRelationshipType)
	Grammar.RegisterReader("Repository", tosca_v2_0.ReadRepository)
	Grammar.RegisterReader("RequirementAssignment", tosca_v1_2.ReadRequirementAssignment) // 1.2
	Grammar.RegisterReader("RequirementDefinition", tosca_v1_3.ReadRequirementDefinition) // 1.3
	Grammar.RegisterReader("RequirementMapping", tosca_v2_0.ReadRequirementMapping)
	Grammar.RegisterReader("ServiceTemplate", tosca_v2_0.ReadServiceTemplate)
	Grammar.RegisterReader("scalar-unit.frequency", tosca_v1_3.ReadScalarUnitFrequency)
	Grammar.RegisterReader("scalar-unit.size", tosca_v1_3.ReadScalarUnitSize)
	Grammar.RegisterReader("scalar-unit.time", tosca_v1_3.ReadScalarUnitTime)
	Grammar.RegisterReader("Schema", tosca_v2_0.ReadSchema)
	Grammar.RegisterReader("SubstitutionMappings", ReadSubstitutionMappings) // override
	Grammar.RegisterReader("timestamp", tosca_v2_0.ReadTimestamp)
	Grammar.RegisterReader("TriggerDefinition", tosca_v1_2.ReadTriggerDefinition) // 1.2
	Grammar.RegisterReader("TriggerDefinitionCondition", tosca_v2_0.ReadTriggerDefinitionCondition)
	Grammar.RegisterReader("Value", tosca_v2_0.ReadValue)
	Grammar.RegisterReader("version", tosca_v2_0.ReadVersion)
	Grammar.RegisterReader("WorkflowActivityCallOperation", tosca_v1_2.ReadWorkflowActivityCallOperation)   // 1.2
	Grammar.RegisterReader("WorkflowActivityDefinition", tosca_v2_0.ReadWorkflowActivityDefinition)         // introduced in TOSCA 1.1
	Grammar.RegisterReader("WorkflowDefinition", tosca_v2_0.ReadWorkflowDefinition)                         // introduced in TOSCA 1.1
	Grammar.RegisterReader("WorkflowPreconditionDefinition", tosca_v2_0.ReadWorkflowPreconditionDefinition) // introduced in TOSCA 1.1
	Grammar.RegisterReader("WorkflowStepDefinition", tosca_v2_0.ReadWorkflowStepDefinition)                 // introduced in TOSCA 1.1

	DefaultScriptletNamespace.RegisterScriptlets(tosca_v2_0.FunctionScriptlets, nil, parsing.MetadataFunctionPrefix+"join")
	DefaultScriptletNamespace.RegisterScriptlets(tosca_v1_3.ConstraintClauseScriptlets, tosca_v1_3.ConstraintClauseNativeArgumentIndexes, parsing.MetadataContraintPrefix+"schema")

	Grammar.InvalidNamespaceCharacters = ":"
}
