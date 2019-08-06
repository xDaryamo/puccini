package tosca_v1_3

import (
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
)

//
// TopologyTemplate
//
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.9
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.8
//

type TopologyTemplate struct {
	*Entity `name:"topology template"`

	Description                *string               `read:"description"`
	NodeTemplates              NodeTemplates         `read:"node_templates,NodeTemplate"`
	RelationshipTemplates      RelationshipTemplates `read:"relationship_templates,RelationshipTemplate"`
	Groups                     Groups                `read:"groups,Group"`
	Policies                   Policies              `read:"policies,Policy"`
	InputParameterDefinitions  ParameterDefinitions  `read:"inputs,ParameterDefinition"`
	OutputParameterDefinitions ParameterDefinitions  `read:"outputs,ParameterDefinition"`
	WorkflowDefinitions        WorkflowDefinitions   `read:"workflows,WorkflowDefinition"`
	SubstitutionMappings       *SubstitutionMappings `read:"substitution_mappings,SubstitutionMappings"`
}

func NewTopologyTemplate(context *tosca.Context) *TopologyTemplate {
	return &TopologyTemplate{
		Entity:                     NewEntity(context),
		InputParameterDefinitions:  make(ParameterDefinitions),
		OutputParameterDefinitions: make(ParameterDefinitions),
		WorkflowDefinitions:        make(WorkflowDefinitions),
	}
}

// tosca.Reader signature
func ReadTopologyTemplate(context *tosca.Context) interface{} {
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
func (self *TopologyTemplate) SetInputs(inputs map[string]interface{}) {
	context := self.Context.FieldChild("inputs", nil)
	for name, data := range inputs {
		childContext := context.MapChild(name, data)
		if definition, ok := self.InputParameterDefinitions[name]; ok {
			if definition.DataType != nil {
				if typeName, ok := definition.DataType.GetInternalTypeName(); ok {
					if typeName == "integer" {
						// In JSON, everything is a float
						// But we want to support inputs coming from JSON
						// So we'll auto-convert
						switch v := childContext.Data.(type) {
						case float64:
							childContext.Data = int64(v)
						case float32:
							childContext.Data = int64(v)
						}
					}
				}
			}

			definition.Value = ReadValue(childContext).(*Value)
		} else {
			childContext.ReportUndefined("input")
		}
	}
}

// tosca.Renderable interface
func (self *TopologyTemplate) Render() {
	log.Info("{render} topology template")

	self.InputParameterDefinitions.Render("input definition", self.Context.FieldChild("inputs", nil))
	self.OutputParameterDefinitions.Render("output definition", self.Context.FieldChild("outputs", nil))
}

func (self *TopologyTemplate) Normalize(s *normal.ServiceTemplate) {
	log.Info("{normalize} topology template")

	if self.Description != nil {
		// Append to description in service template
		if s.Description != "" {
			s.Description += "\n\n"
		}
		s.Description += *self.Description
	}

	self.InputParameterDefinitions.Normalize(s.Inputs, self.Context.FieldChild("inputs", nil))
	self.OutputParameterDefinitions.Normalize(s.Outputs, self.Context.FieldChild("outputs", nil))

	self.NodeTemplates.Normalize(s)
	self.Groups.Normalize(s)

	// Workflows must be normalized after node templates and groups
	// (because step activities might call operations on them)
	self.WorkflowDefinitions.Normalize(s)

	// Policies must be normalized after workflows
	// (because policy triggers might call them)
	self.Policies.Normalize(s)

	if self.SubstitutionMappings != nil {
		self.SubstitutionMappings.Normalize(s)
	}
}
