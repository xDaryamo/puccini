package ard

import (
	"gopkg.in/yaml.v3"
)

//
// Locator
//

type Locator interface {
	Locate(path ...PathElement) (int, int, bool)
}

//
// YamlLocator
//

type YamlLocator struct {
	RootNode *yaml.Node
}

func NewYamlLocator(rootNode *yaml.Node) *YamlLocator {
	return &YamlLocator{rootNode}
}

// Locator interface
func (self *YamlLocator) Locate(path ...PathElement) (int, int, bool) {
	if node := FindYamlNode(self.RootNode, path...); node != nil {
		return node.Line, node.Column, true
	} else {
		return 0, 0, false
	}
}
