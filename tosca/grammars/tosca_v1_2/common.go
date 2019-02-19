package tosca_v1_2

import (
	"github.com/op/go-logging"
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/grammars/tosca_v1_1"
)

var log = logging.MustGetLogger("grammars.tosca_v1_2")

var Readers = make(map[string]tosca.Reader)

func init() {
	Readers["Artifact"] = tosca_v1_1.ReadArtifact
	Readers["ArtifactDefinition"] = tosca_v1_1.ReadArtifactDefinition
	Readers["ArtifactType"] = tosca_v1_1.ReadArtifactType
	Readers["AttributeDefinition"] = tosca_v1_1.ReadAttributeDefinition
	Readers["CapabilityAssignment"] = tosca_v1_1.ReadCapabilityAssignment
	Readers["CapabilityDefinition"] = tosca_v1_1.ReadCapabilityDefinition
	Readers["CapabilityFilter"] = tosca_v1_1.ReadCapabilityFilter
	Readers["CapabilityMapping"] = tosca_v1_1.ReadCapabilityMapping
	Readers["CapabilityType"] = tosca_v1_1.ReadCapabilityType
	Readers["ConditionClause"] = tosca_v1_1.ReadConditionClause
	Readers["ConstraintClause"] = tosca_v1_1.ReadConstraintClause
	Readers["DataType"] = tosca_v1_1.ReadDataType
	Readers["EntrySchema"] = tosca_v1_1.ReadEntrySchema
	Readers["EventFilter"] = tosca_v1_1.ReadEventFilter
	Readers["Group"] = tosca_v1_1.ReadGroup
	Readers["GroupType"] = tosca_v1_1.ReadGroupType
	Readers["Import"] = tosca_v1_1.ReadImport
	Readers["InterfaceAssignment"] = tosca_v1_1.ReadInterfaceAssignment
	Readers["InterfaceDefinition"] = tosca_v1_1.ReadInterfaceDefinition
	Readers["InterfaceMapping"] = ReadInterfaceMapping // new
	Readers["InterfaceType"] = tosca_v1_1.ReadInterfaceType
	Readers["Metadata"] = tosca_v1_1.ReadMetadata
	Readers["NodeFilter"] = tosca_v1_1.ReadNodeFilter
	Readers["NodeTemplate"] = tosca_v1_1.ReadNodeTemplate
	Readers["NodeType"] = tosca_v1_1.ReadNodeType
	Readers["OperationAssignment"] = tosca_v1_1.ReadOperationAssignment
	Readers["OperationDefinition"] = tosca_v1_1.ReadOperationDefinition
	Readers["OperationImplementation"] = tosca_v1_1.ReadOperationImplementation
	Readers["ParameterDefinition"] = tosca_v1_1.ReadParameterDefinition
	Readers["Policy"] = tosca_v1_1.ReadPolicy
	Readers["PolicyType"] = tosca_v1_1.ReadPolicyType
	Readers["PropertyDefinition"] = tosca_v1_1.ReadPropertyDefinition
	Readers["PropertyFilter"] = tosca_v1_1.ReadPropertyFilter
	Readers["PropertyMapping"] = ReadPropertyMapping // new
	Readers["range"] = tosca_v1_1.ReadRange
	Readers["RangeEntity"] = tosca_v1_1.ReadRangeEntity
	Readers["RelationshipAssignment"] = tosca_v1_1.ReadRelationshipAssignment
	Readers["RelationshipDefinition"] = tosca_v1_1.ReadRelationshipDefinition
	Readers["RelationshipTemplate"] = tosca_v1_1.ReadRelationshipTemplate
	Readers["RelationshipType"] = tosca_v1_1.ReadRelationshipType
	Readers["Repository"] = tosca_v1_1.ReadRepository
	Readers["RequirementAssignment"] = tosca_v1_1.ReadRequirementAssignment
	Readers["RequirementDefinition"] = tosca_v1_1.ReadRequirementDefinition
	Readers["RequirementMapping"] = tosca_v1_1.ReadRequirementMapping
	Readers["scalar-unit.size"] = tosca_v1_1.ReadScalarUnitSize
	Readers["scalar-unit.time"] = tosca_v1_1.ReadScalarUnitTime
	Readers["scalar-unit.frequency"] = tosca_v1_1.ReadScalarUnitFrequency
	Readers["ServiceTemplate"] = ReadServiceTemplate           // override
	Readers["SubstitutionMappings"] = ReadSubstitutionMappings // override
	Readers["timestamp"] = tosca_v1_1.ReadTimestamp
	Readers["TopologyTemplate"] = ReadTopologyTemplate // override
	Readers["TriggerDefinition"] = tosca_v1_1.ReadTriggerDefinition
	Readers["TriggerDefinitionCondition"] = tosca_v1_1.ReadTriggerDefinitionCondition
	Readers["Unit"] = tosca_v1_1.ReadUnit
	Readers["Value"] = tosca_v1_1.ReadValue
	Readers["version"] = tosca_v1_1.ReadVersion
	Readers["WorkflowActivityDefinition"] = tosca_v1_1.ReadWorkflowActivityDefinition
	Readers["WorkflowDefinition"] = tosca_v1_1.ReadWorkflowDefinition
	Readers["WorkflowPreconditionDefinition"] = tosca_v1_1.ReadWorkflowPreconditionDefinition
	Readers["WorkflowStepDefinition"] = tosca_v1_1.ReadWorkflowStepDefinition
}
