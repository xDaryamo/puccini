package normal

import (
	"github.com/tliron/kutil/ard"
)

//
// NodeTemplate
//

type NodeTemplate struct {
	ServiceTemplate *ServiceTemplate `json:"-" yaml:"-"`
	Name            string           `json:"-" yaml:"-"`

	Metadata     map[string]string `json:"metadata" yaml:"metadata"`
	Description  string            `json:"description" yaml:"description"`
	Types        Types             `json:"types" yaml:"types"`
	Directives   []string          `json:"directives" yaml:"directives"`
	Properties   Constrainables    `json:"properties" yaml:"properties"`
	Attributes   Constrainables    `json:"attributes" yaml:"attributes"`
	Requirements Requirements      `json:"requirements" yaml:"requirements"`
	Capabilities Capabilities      `json:"capabilities" yaml:"capabilities"`
	Interfaces   Interfaces        `json:"interfaces" yaml:"interfaces"`
	Artifacts    Artifacts         `json:"artifacts" yaml:"artifacts"`
}

func (self *ServiceTemplate) NewNodeTemplate(name string) *NodeTemplate {
	nodeTemplate := &NodeTemplate{
		ServiceTemplate: self,
		Name:            name,
		Metadata:        make(map[string]string),
		Types:           make(Types),
		Directives:      make([]string, 0),
		Properties:      make(Constrainables),
		Attributes:      make(Constrainables),
		Requirements:    make(Requirements, 0),
		Capabilities:    make(Capabilities),
		Interfaces:      make(Interfaces),
		Artifacts:       make(Artifacts),
	}
	self.NodeTemplates[name] = nodeTemplate
	return nodeTemplate
}

//
// NodeTemplates
//

type NodeTemplates map[string]*NodeTemplate

// For access in JavaScript
func (self NodeTemplates) Object(name string) map[string]interface{} {
	// JavaScript requires keys to be strings, so we would lose complex keys
	o := make(ard.StringMap)
	for key, nodeTemplate := range self {
		o[key] = nodeTemplate
	}
	return o
}
