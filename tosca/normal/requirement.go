package normal

import (
	"encoding/json"
)

//
// Requirement
//

type Requirement struct {
	SourceNodeTemplate              *NodeTemplate
	Name                            string
	CapabilityTypeName              *string
	CapabilityName                  *string
	NodeTypeName                    *string
	NodeTemplate                    *NodeTemplate
	NodeTemplatePropertyConstraints FunctionCallMap
	CapabilityPropertyConstraints   FunctionCallMapMap
	Relationship                    *Relationship
	Path                            string
}

func (self *NodeTemplate) NewRequirement(name string, path string) *Requirement {
	requirement := &Requirement{
		SourceNodeTemplate:              self,
		Name:                            name,
		NodeTemplatePropertyConstraints: make(FunctionCallMap),
		CapabilityPropertyConstraints:   make(FunctionCallMapMap),
		Path:                            path,
	}
	self.Requirements = append(self.Requirements, requirement)
	return requirement
}

func (self *Requirement) Marshalable() interface{} {
	var capabilityTypeName string
	var capabilityName string
	var nodeTypeName string
	var nodeTemplateName string
	if self.CapabilityTypeName != nil {
		capabilityTypeName = *self.CapabilityTypeName
	}
	if self.CapabilityName != nil {
		capabilityName = *self.CapabilityName
	}
	if self.NodeTypeName != nil {
		nodeTypeName = *self.NodeTypeName
	}
	if self.NodeTemplate != nil {
		nodeTemplateName = self.NodeTemplate.Name
	}

	return &struct {
		Name                            string             `json:"name" yaml:"name"`
		CapabilityTypeName              string             `json:"capabilityTypeName" yaml:"capabilityTypeName"`
		CapabilityName                  string             `json:"capabilityName" yaml:"capabilityName"`
		NodeTypeName                    string             `json:"nodeTypeName" yaml:"nodeTypeName" `
		NodeTemplateName                string             `json:"nodeTemplateName" yaml:"nodeTemplateName"`
		NodeTemplatePropertyConstraints FunctionCallMap    `json:"nodeTemplatePropertyConstraints" yaml:"nodeTemplatePropertyConstraints"`
		CapabilityPropertyConstraints   FunctionCallMapMap `json:"capabilityPropertyConstraints" yaml:"capabilityPropertyConstraints"`
		Relationship                    *Relationship      `json:"relationship" yaml:"relationship"`
		Path                            string             `json:"path" yaml:"path"`
	}{
		Name:                            self.Name,
		CapabilityTypeName:              capabilityTypeName,
		CapabilityName:                  capabilityName,
		NodeTypeName:                    nodeTypeName,
		NodeTemplateName:                nodeTemplateName,
		NodeTemplatePropertyConstraints: self.NodeTemplatePropertyConstraints,
		CapabilityPropertyConstraints:   self.CapabilityPropertyConstraints,
		Relationship:                    self.Relationship,
		Path:                            self.Path,
	}
}

// json.Marshaler interface
func (self *Requirement) MarshalJSON() ([]byte, error) {
	return json.Marshal(self.Marshalable())
}

// yaml.Marshaler interface
func (self *Requirement) MarshalYAML() (interface{}, error) {
	return self.Marshalable(), nil
}

//
// Requirements
//

type Requirements []*Requirement
