package tosca_v2_0

import (
	"github.com/tliron/go-ard"
	"github.com/tliron/puccini/tosca/parsing"
)

//
// InterfaceMapping
//
// [TOSCA-v2.0] @ 15.6 Interface Mapping
// Interface mapping allows an interface operation on the substituted node
// to be mapped to workflow in the substituting service template.
//
// Grammar:
// <interface_name>:
//   <operation_name>: <workflow_name>
//

type InterfaceMapping struct {
	*Entity `name:"interface mapping"`
	Name    string

	// TOSCA 2.0: Maps operation names to workflow names
	OperationMappings map[string]string `json:"operationMappings" yaml:"operationMappings"`
}

func NewInterfaceMapping(context *parsing.Context) *InterfaceMapping {
	return &InterfaceMapping{
		Entity:            NewEntity(context),
		Name:              context.Name,
		OperationMappings: make(map[string]string),
	}
}

// ([parsing.Reader] signature)
func ReadInterfaceMapping(context *parsing.Context) parsing.EntityPtr {
	self := NewInterfaceMapping(context)

	// TOSCA 2.0: Read operation mappings as a map
	if context.ValidateType(ard.TypeMap) {
		if operationMappings := context.ReadStringMap(); operationMappings != nil {
			// Convert map[string]any to map[string]string
			for operationName, workflowName := range *operationMappings {
				if workflowNameStr, ok := workflowName.(string); ok {
					self.OperationMappings[operationName] = workflowNameStr
				} else {
					context.MapChild(operationName, workflowName).ReportValueMalformed("interface mapping", "workflow name must be a string")
				}
			}
		}
	}

	return self
}

// ([parsing.Mappable] interface)
func (self *InterfaceMapping) GetKey() string {
	return self.Name
}

// ([parsing.Renderable] interface)
func (self *InterfaceMapping) Render() {
	// Avoid rendering more than once
	self.renderOnce.Do(self.render)
}

func (self *InterfaceMapping) render() {
	logRender.Debug("interface mapping")

	// Check if this is a TOSCA 1.3 interface mapping
	if _, isTosca13 := self.OperationMappings["__TOSCA_1_3__"]; isTosca13 {
		// Skip workflow validation for TOSCA 1.3 format
		return
	}

	// TOSCA 2.0: Validate that workflows exist in the service template
	if self.OperationMappings != nil {
		for operationName, workflowName := range self.OperationMappings {
			if !self.validateWorkflowExists(workflowName) {
				self.Context.MapChild(operationName, workflowName).ReportUnknown("workflow")
			}
		}
	}
}

// validateWorkflowExists checks if a workflow exists in the service template
func (self *InterfaceMapping) validateWorkflowExists(workflowName string) bool {
	// Find the service template context
	serviceContext := self.findServiceTemplateContext()
	if serviceContext == nil {
		return false
	}

	// Check if workflows section exists
	workflowsContext, ok := serviceContext.GetFieldChild("workflows")
	if !ok {
		return false
	}

	// Check if the specific workflow exists
	_, ok = workflowsContext.GetFieldChild(workflowName)
	return ok
}

// findServiceTemplateContext finds the parent service template context
func (self *InterfaceMapping) findServiceTemplateContext() *parsing.Context {
	context := self.Context
	for context.Parent != nil && context.Name != "service_template" {
		context = context.Parent
	}
	if context.Name == "service_template" {
		return context
	}
	return nil
}

//
// InterfaceMappings
//

type InterfaceMappings map[string]*InterfaceMapping
