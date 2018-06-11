package normal

import (
	"fmt"

	"github.com/tliron/puccini/format"
)

//
// NodeTemplate
//

type NodeTemplate struct {
	ServiceTemplate *ServiceTemplate `json:"-" yaml:"-"`
	Name            string           `json:"-" yaml:"-"`
	Description     string           `json:"description" yaml:"description"`
	Types           Types            `json:"types" yaml:"types"`
	Directives      []string         `json:"directives" yaml:"directives"`
	Properties      Constrainables   `json:"properties" yaml:"properties"`
	Attributes      Constrainables   `json:"attributes" yaml:"attributes"`
	Capabilities    Capabilities     `json:"capabilities" yaml:"capabilities"`
	Relationships   Relationships    `json:"relationships" yaml:"relationships"`
	Interfaces      Interfaces       `json:"interfaces" yaml:"interfaces"`
	Artifacts       Artifacts        `json:"artifacts" yaml:"artifacts"`
	Policies        []*Policy        `json:"-" yaml:"-"`
	Groups          []*Group         `json:"-" yaml:"-"`
}

func (self *ServiceTemplate) NewNodeTemplate(name string) *NodeTemplate {
	nodeTemplate := &NodeTemplate{
		ServiceTemplate: self,
		Name:            name,
		Types:           make(Types),
		Directives:      make([]string, 0),
		Properties:      make(Constrainables),
		Attributes:      make(Constrainables),
		Capabilities:    make(Capabilities),
		Relationships:   make(Relationships, 0),
		Interfaces:      make(Interfaces),
		Artifacts:       make(Artifacts),
		Policies:        make([]*Policy, 0),
		Groups:          make([]*Group, 0),
	}
	self.NodeTemplates[name] = nodeTemplate
	return nodeTemplate
}

// Print

func (self *NodeTemplate) Print(indent int) {
	format.PrintIndent(indent)
	fmt.Printf("%s\n", format.ColorTypeName(self.Name))

	length := len(self.Relationships)
	last := length - 1

	var treePrefix format.TreePrefix
	for i, relationship := range self.Relationships {
		isLast := i == last
		relationship.Print(indent, treePrefix, isLast)
	}
}

//
// NodeTemplates
//

type NodeTemplates map[string]*NodeTemplate

// For access in JavaScript
func (self NodeTemplates) Object(name string) map[string]interface{} {
	o := make(map[string]interface{})
	for key, nodeTemplate := range self {
		o[key] = nodeTemplate
	}
	return o
}
