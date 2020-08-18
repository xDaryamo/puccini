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

	ToscaDefinitionsVersion *string            `read:"tosca_definitions_version" require:""`
	Metadata                Metadata           `read:"metadata,!Metadata"` // not in spec, but in code
	Imports                 Imports            `read:"imports,[]Import"`
	Inputs                  Inputs             `read:"inputs,Input"`
	NodeTemplates           NodeTemplates      `read:"node_templates,NodeTemplate"`
	NodeTypes               NodeTypes          `read:"node_types,NodeType" hierarchy:""`
	Capabilities            []*ValueDefinition `read:"capability,ValueDefinition"`
	Outputs                 ValueDefinitions   `read:"outputs,ValueDefinition"`
	RelationshipTypes       RelationshipTypes  `read:"relationships,RelationshipType" hierarchy:""`
	Plugins                 Plugins            `read:"plugins,Plugin"`
	Workflows               Workflows          `read:"workflows,Workflow"`
	DataTypes               DataTypes          `read:"data_types,DataType" hierarchy:""`
	Policies                Policies           `read:"policies,Policy"`
	PolicyTypes             PolicyTypes        `read:"policy_types,PolicyType" hierarchy:""`
	PolicyTriggerTypes      PolicyTriggerTypes `read:"policy_triggers,PolicyTriggerType" hierarchy:""`
	UploadResources         *UploadResources   `read:"upload_resources,UploadResources"`
}

func NewUnit(context *tosca.Context) *Unit {
	return &Unit{
		Entity:  NewEntity(context),
		Inputs:  make(Inputs),
		Outputs: make(ValueDefinitions),
	}
}

// tosca.Reader signature
func ReadUnit(context *tosca.Context) tosca.EntityPtr {
	self := NewUnit(context)
	context.ScriptletNamespace.Merge(DefaultScriptletNamespace)
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
