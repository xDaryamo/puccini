package cloudify_v1_3

import (
	"github.com/tliron/kutil/logging"
	"github.com/tliron/puccini/tosca"
)

var log = logging.GetLogger("puccini.grammars.cloudify_v1_3")
var logInherit = logging.NewSubLogger(log, "inherit")
var logRender = logging.NewSubLogger(log, "render")
var logNormalize = logging.NewSubLogger(log, "normalize")

var Grammar = tosca.NewGrammar()

var DefaultScriptletNamespace = tosca.NewScriptletNamespace()

func init() {
	Grammar.RegisterVersion("tosca_definitions_version", "cloudify_dsl_1_3", "/cloudify/5.0.5/profile.yaml")

	Grammar.RegisterReader("$Root", ReadBlueprint)
	Grammar.RegisterReader("$Unit", ReadUnit)

	Grammar.RegisterReader("Blueprint", ReadBlueprint)
	Grammar.RegisterReader("DataType", ReadDataType)
	Grammar.RegisterReader("Group", ReadGroup)
	Grammar.RegisterReader("GroupPolicy", ReadGroupPolicy)
	Grammar.RegisterReader("GroupPolicyTrigger", ReadGroupPolicyTrigger)
	Grammar.RegisterReader("DSLResource", ReadDSLResource)
	Grammar.RegisterReader("Import", ReadImport)
	Grammar.RegisterReader("Input", ReadInput)
	Grammar.RegisterReader("InterfaceAssignment", ReadInterfaceAssignment)
	Grammar.RegisterReader("InterfaceDefinition", ReadInterfaceDefinition)
	Grammar.RegisterReader("Metadata", ReadMetadata)
	Grammar.RegisterReader("NodeTemplate", ReadNodeTemplate)
	Grammar.RegisterReader("NodeTemplateCapability", ReadNodeTemplateCapability)
	Grammar.RegisterReader("NodeTemplateInstances", ReadNodeTemplateInstances)
	Grammar.RegisterReader("NodeType", ReadNodeType)
	Grammar.RegisterReader("OperationDefinition", ReadOperationDefinition)
	Grammar.RegisterReader("OperationAssignment", ReadOperationAssignment)
	Grammar.RegisterReader("ParameterDefinition", ReadParameterDefinition)
	Grammar.RegisterReader("Plugin", ReadPlugin)
	Grammar.RegisterReader("Policy", ReadPolicy)
	Grammar.RegisterReader("PolicyTriggerType", ReadPolicyTriggerType)
	Grammar.RegisterReader("PolicyType", ReadPolicyType)
	Grammar.RegisterReader("PropertyDefinition", ReadPropertyDefinition)
	Grammar.RegisterReader("RelationshipType", ReadRelationshipType)
	Grammar.RegisterReader("RelationshipAssignment", ReadRelationshipAssignment)
	Grammar.RegisterReader("UploadResources", ReadUploadResources)
	Grammar.RegisterReader("Value", ReadValue)
	Grammar.RegisterReader("ValueDefinition", ReadValueDefinition)
	Grammar.RegisterReader("Workflow", ReadWorkflow)

	DefaultScriptletNamespace.RegisterScriptlets(FunctionScriptlets, nil)
}
