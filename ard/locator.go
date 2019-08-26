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
	}
	return 0, 0, false
}

// Utils

func FindYamlNode(node *yaml.Node, path ...PathElement) *yaml.Node {
	if len(path) == 0 {
		return node
	}

	switch node.Kind {
	case yaml.AliasNode:
		return FindYamlNode(node.Alias, path...)

	case yaml.DocumentNode:
		for _, childNode := range node.Content {
			if foundNode := FindYamlNode(childNode, path...); foundNode != nil {
				return foundNode
			}
		}

	case yaml.MappingNode:
		pathElement := path[0]
		switch pathElement.Type {
		case FieldPathType, MapPathType:
			v := pathElement.Value.(string)

			// Content is a slice of pairs of key followed by value
			length := len(node.Content)
			for i := 0; i < length; i += 2 {
				keyNode := node.Content[i]
				if (keyNode.Kind == yaml.ScalarNode) && (keyNode.Tag == "!!str") && (keyNode.Value == v) {
					if i+1 < length {
						valueNode := node.Content[i+1]
						foundNode := FindYamlNode(valueNode, path[1:]...)
						if foundNode == valueNode {
							// We will use the key node for the location instead of the value node
							return keyNode
						}
						return foundNode
					} else {
						// Missing value - content is malformed?
						return keyNode
					}
				}
			}

			return node
		}

	case yaml.SequenceNode:
		pathElement := path[0]
		switch pathElement.Type {
		case ListPathType:
			index := pathElement.Value.(int)
			if index < len(node.Content) {
				return FindYamlNode(node.Content[index], path[1:]...)
			}
			return node
		}
	}

	return nil
}
