package v1_2

import (
	"github.com/op/go-logging"
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/grammars/v1_1"
)

var log = logging.MustGetLogger("grammars.v1_2")

var Readers = make(map[string]tosca.Reader)

func init() {
	Readers["Artifact"] = v1_1.ReadArtifact
	Readers["ArtifactDefinition"] = v1_1.ReadArtifactDefinition
	Readers["ArtifactType"] = v1_1.ReadArtifactType
	Readers["AttributeDefinition"] = v1_1.ReadAttributeDefinition
	Readers["CapabilityAssignment"] = v1_1.ReadCapabilityAssignment
	Readers["CapabilityDefinition"] = v1_1.ReadCapabilityDefinition
	Readers["CapabilityFilter"] = v1_1.ReadCapabilityFilter
	Readers["CapabilityMapping"] = v1_1.ReadCapabilityMapping
	Readers["CapabilityType"] = v1_1.ReadCapabilityType
	Readers["ConditionClause"] = v1_1.ReadConditionClause
	Readers["ConstraintClause"] = v1_1.ReadConstraintClause
	Readers["DataType"] = v1_1.ReadDataType
	Readers["EntrySchema"] = v1_1.ReadEntrySchema
	Readers["EventFilter"] = v1_1.ReadEventFilter
	Readers["Group"] = v1_1.ReadGroup
	Readers["GroupType"] = v1_1.ReadGroupType
	Readers["Import"] = v1_1.ReadImport
	Readers["InterfaceAssignment"] = v1_1.ReadInterfaceAssignment
	Readers["InterfaceDefinition"] = v1_1.ReadInterfaceDefinition
	Readers["InterfaceMapping"] = ReadInterfaceMapping // new
	Readers["InterfaceType"] = v1_1.ReadInterfaceType
	Readers["Metadata"] = v1_1.ReadMetadata
	Readers["NodeFilter"] = v1_1.ReadNodeFilter
	Readers["NodeTemplate"] = v1_1.ReadNodeTemplate
	Readers["NodeType"] = v1_1.ReadNodeType
	Readers["OperationAssignment"] = v1_1.ReadOperationAssignment
	Readers["OperationDefinition"] = v1_1.ReadOperationDefinition
	Readers["OperationImplementation"] = v1_1.ReadOperationImplementation
	Readers["ParameterDefinition"] = v1_1.ReadParameterDefinition
	Readers["Policy"] = v1_1.ReadPolicy
	Readers["PolicyType"] = v1_1.ReadPolicyType
	Readers["PropertyDefinition"] = v1_1.ReadPropertyDefinition
	Readers["PropertyFilter"] = v1_1.ReadPropertyFilter
	Readers["PropertyMapping"] = ReadPropertyMapping // new
	Readers["range"] = v1_1.ReadRange
	Readers["RangeEntity"] = v1_1.ReadRangeEntity
	Readers["RelationshipAssignment"] = v1_1.ReadRelationshipAssignment
	Readers["RelationshipDefinition"] = v1_1.ReadRelationshipDefinition
	Readers["RelationshipTemplate"] = v1_1.ReadRelationshipTemplate
	Readers["RelationshipType"] = v1_1.ReadRelationshipType
	Readers["Repository"] = v1_1.ReadRepository
	Readers["RequirementAssignment"] = v1_1.ReadRequirementAssignment
	Readers["RequirementDefinition"] = v1_1.ReadRequirementDefinition
	Readers["RequirementMapping"] = v1_1.ReadRequirementMapping
	Readers["scalar-unit.size"] = v1_1.ReadScalarUnitSize
	Readers["scalar-unit.time"] = v1_1.ReadScalarUnitTime
	Readers["scalar-unit.frequency"] = v1_1.ReadScalarUnitFrequency
	Readers["ServiceTemplate"] = ReadServiceTemplate           // override
	Readers["SubstitutionMappings"] = ReadSubstitutionMappings // override
	Readers["timestamp"] = v1_1.ReadTimestamp
	Readers["TopologyTemplate"] = ReadTopologyTemplate // override
	Readers["TriggerDefinition"] = v1_1.ReadTriggerDefinition
	Readers["TriggerDefinitionCondition"] = v1_1.ReadTriggerDefinitionCondition
	Readers["Unit"] = v1_1.ReadUnit
	Readers["Value"] = v1_1.ReadValue
	Readers["version"] = v1_1.ReadVersion
	Readers["WorkflowActivityDefinition"] = v1_1.ReadWorkflowActivityDefinition
	Readers["WorkflowDefinition"] = v1_1.ReadWorkflowDefinition
	Readers["WorkflowPreconditionDefinition"] = v1_1.ReadWorkflowPreconditionDefinition
	Readers["WorkflowStepDefinition"] = v1_1.ReadWorkflowStepDefinition
}
