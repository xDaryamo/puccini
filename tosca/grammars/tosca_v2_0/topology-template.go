package tosca_v2_0

import (
	"github.com/tliron/kutil/ard"
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
)

//
// TopologyTemplate
//
// [TOSCA-v2.0] @ ?
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.9
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.9
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.8
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.8
//

type TopologyTemplate struct {
	*Entity `name:"topology template"`

	Description           *string               `read:"description"`
	NodeTemplates         NodeTemplates         `read:"node_templates,NodeTemplate"`
	RelationshipTemplates RelationshipTemplates `read:"relationship_templates,RelationshipTemplate"`
	Groups                Groups                `read:"groups,Group"`
	Policies              Policies              `read:"policies,<>Policy"`
	InputDefinitions      ParameterDefinitions  `read:"inputs,ParameterDefinition"`
	OutputDefinitions     ParameterDefinitions  `read:"outputs,ParameterDefinition"`
	WorkflowDefinitions   WorkflowDefinitions   `read:"workflows,WorkflowDefinition"`
	SubstitutionMappings  *SubstitutionMappings `read:"substitution_mappings,SubstitutionMappings"`
}

func NewTopologyTemplate(context *tosca.Context) *TopologyTemplate {
	return &TopologyTemplate{
		Entity:              NewEntity(context),
		InputDefinitions:    make(ParameterDefinitions),
		OutputDefinitions:   make(ParameterDefinitions),
		WorkflowDefinitions: make(WorkflowDefinitions),
	}
}

// tosca.Reader signature
func ReadTopologyTemplate(context *tosca.Context) tosca.EntityPtr {
	self := NewTopologyTemplate(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))
	return self
}

func (self *TopologyTemplate) GetNodeTemplatesOfType(nodeType *NodeType) []*NodeTemplate {
	var nodeTemplates []*NodeTemplate
	for _, nodeTemplate := range self.NodeTemplates {
		if (nodeTemplate.NodeType != nil) && self.Context.Hierarchy.IsCompatible(nodeType, nodeTemplate.NodeType) {
			nodeTemplates = append(nodeTemplates, nodeTemplate)
		}
	}
	return nodeTemplates
}

// parser.HasInputs interface
func (self *TopologyTemplate) SetInputs(inputs map[string]ard.Value) {
	context := self.Context.FieldChild("inputs", nil)
	for name, data := range inputs {
		childContext := context.MapChild(name, data)
		if definition, ok := self.InputDefinitions[name]; ok {
			if definition.DataType != nil {
				if internalTypeName, ok := definition.DataType.GetInternalTypeName(); ok {
					if internalTypeName == ard.TypeInteger {
						// In JSON, everything is a float
						// But we want to support inputs coming from JSON
						// So we'll auto-convert
						switch data := childContext.Data.(type) {
						case float64:
							childContext.Data = int64(data)
						case float32:
							childContext.Data = int64(data)
						}
					}
				}
			}

			definition.Value = ReadValue(childContext).(*Value)
		} else {
			childContext.ReportUndeclared("input")
		}
	}
}

// parser.Renderable interface
func (self *TopologyTemplate) Render() {
	logRender.Debug("topology template")

	var mappedInputs []string

	if self.SubstitutionMappings != nil {
		// Substitution mapping rendering has to happen before input rendering
		// in order to avoid rendering of mapped inputs
		self.SubstitutionMappings.Render(self.InputDefinitions)

		for _, mapping := range self.SubstitutionMappings.PropertyMappings {
			if mapping.InputDefinition != nil {
				mappedInputs = append(mappedInputs, mapping.InputDefinition.Name)
			}
		}
	}

	self.InputDefinitions.Render("input definition", mappedInputs, self.Context.FieldChild("inputs", nil))
	self.OutputDefinitions.Render("output definition", nil, self.Context.FieldChild("outputs", nil))
}

func (self *TopologyTemplate) Normalize(normalServiceTemplate *normal.ServiceTemplate) {
	logNormalize.Debug("topology template")

	if self.Description != nil {
		// Append to description in service template
		if normalServiceTemplate.Description != "" {
			normalServiceTemplate.Description += "\n\n"
		}
		normalServiceTemplate.Description += *self.Description
	}

	self.InputDefinitions.Normalize(normalServiceTemplate.Inputs, self.Context.FieldChild("inputs", nil))
	self.OutputDefinitions.Normalize(normalServiceTemplate.Outputs, self.Context.FieldChild("outputs", nil))

	self.NodeTemplates.Normalize(normalServiceTemplate)
	self.Groups.Normalize(normalServiceTemplate)

	// Workflows must be normalized after node templates and groups
	// (because step activities might call operations on them)
	self.WorkflowDefinitions.Normalize(normalServiceTemplate)

	// Policies must be normalized after workflows
	// (because policy triggers might call them)
	self.Policies.Normalize(normalServiceTemplate)

	if self.SubstitutionMappings != nil {
		self.SubstitutionMappings.Normalize(normalServiceTemplate)
	}
}
