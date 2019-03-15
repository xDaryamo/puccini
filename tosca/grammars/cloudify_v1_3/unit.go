package cloudify_v1_3

import (
	"github.com/tliron/puccini/tosca"
)

//
// Unit
//
// See Blueprint
//

type Unit struct {
	*Entity `name:"unit"`

	ToscaDefinitionsVersion *string              `read:"tosca_definitions_version" require:"tosca_definitions_version"`
	Imports                 []*Import            `read:"imports,[]Import"`
	Inputs                  []*Input             `read:"inputs,Input"`
	NodeTemplates           []*NodeTemplate      `read:"node_templates,NodeTemplate"`
	NodeTypes               []*NodeType          `read:"node_types,NodeType" hierarchy:""`
	Capabilities            []*ValueDefinition   `read:"capability,ValueDefinition"`
	Outputs                 []*ValueDefinition   `read:"outputs,ValueDefinition"`
	RelationshipTypes       []*RelationshipType  `read:"relationships,RelationshipType" hierarchy:""`
	Plugins                 []*Plugin            `read:"plugins,Plugin"`
	Workflows               []*Workflow          `read:"workflows,Workflow"`
	DataTypes               []*DataType          `read:"data_types,DataType" hierarchy:""`
	Policies                []*Policy            `read:"policies,Policy"`
	PolicyTypes             []*PolicyType        `read:"policy_types,PolicyType"`
	PolicyTriggerTypes      []*PolicyTriggerType `read:"policy_triggers,PolicyTriggerType"`
	UploadResources         *UploadResources     `read:"upload_resources,UploadResources"`
}

func NewUnit(context *tosca.Context) *Unit {
	return &Unit{
		Entity: NewEntity(context),
	}
}

// tosca.Reader signature
func ReadUnit(context *tosca.Context) interface{} {
	self := NewUnit(context)
	context.ScriptNamespace.Merge(DefaultScriptNamespace)
	context.ValidateUnsupportedFields(append(context.ReadFields(self), "dsl_definitions"))
	return self
}

// parser.Importer interface
func (self *Unit) GetImportSpecs() []*tosca.ImportSpec {
	var importSpecs = make([]*tosca.ImportSpec, 0, len(self.Imports))
	for _, import_ := range self.Imports {
		if importSpec, ok := import_.NewImportSpec(self); ok {
			importSpecs = append(importSpecs, importSpec)
		}
	}
	return importSpecs
}
