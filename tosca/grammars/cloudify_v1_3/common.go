package cloudify_v1_3

import (
	"github.com/op/go-logging"
	"github.com/tliron/puccini/js"
	"github.com/tliron/puccini/tosca"
)

var log = logging.MustGetLogger("grammars.cloudify_v1_3")

var Grammar = make(tosca.Grammar)

var DefaultScriptNamespace = make(tosca.ScriptNamespace)

func init() {
	Grammar["ServiceTemplate"] = ReadBlueprint

	Grammar["Blueprint"] = ReadBlueprint
	Grammar["DataType"] = ReadDataType
	Grammar["Group"] = ReadGroup
	Grammar["GroupPolicy"] = ReadGroupPolicy
	Grammar["GroupPolicyTrigger"] = ReadGroupPolicyTrigger
	Grammar["DSLResource"] = ReadDSLResource
	Grammar["Import"] = ReadImport
	Grammar["Input"] = ReadInput
	Grammar["InterfaceAssignment"] = ReadInterfaceAssignment
	Grammar["InterfaceDefinition"] = ReadInterfaceDefinition
	Grammar["NodeTemplate"] = ReadNodeTemplate
	Grammar["NodeTemplateCapability"] = ReadNodeTemplateCapability
	Grammar["NodeTemplateInstances"] = ReadNodeTemplateInstances
	Grammar["NodeType"] = ReadNodeType
	Grammar["OperationDefinition"] = ReadOperationDefinition
	Grammar["OperationAssignment"] = ReadOperationAssignment
	Grammar["ParameterDefinition"] = ReadParameterDefinition
	Grammar["Plugin"] = ReadPlugin
	Grammar["Policy"] = ReadPolicy
	Grammar["PolicyTriggerType"] = ReadPolicyTriggerType
	Grammar["PolicyType"] = ReadPolicyType
	Grammar["PropertyDefinition"] = ReadPropertyDefinition
	Grammar["RelationshipType"] = ReadRelationshipType
	Grammar["RelationshipAssignment"] = ReadRelationshipAssignment
	Grammar["UploadResources"] = ReadUploadResources
	Grammar["Unit"] = ReadUnit
	Grammar["Value"] = ReadValue
	Grammar["ValueDefinition"] = ReadValueDefinition
	Grammar["Workflow"] = ReadWorkflow

	for name, sourceCode := range FunctionSourceCode {
		DefaultScriptNamespace[name] = &tosca.Script{
			SourceCode: js.Cleanup(sourceCode),
		}
	}
}
