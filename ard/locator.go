package ard

import (
	"fmt"

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

func (self *YamlLocator) Locate(path ...PathElement) (int, int, bool) {
	if node := FindYamlNode(self.RootNode, path...); node != nil {
		return node.Line, node.Column, true
	}
	return 0, 0, false
}

func FindYamlNode(node *yaml.Node, path ...PathElement) *yaml.Node {
	switch node.Kind {
	case yaml.DocumentNode:
		for _, child := range node.Content {
			if found := FindYamlNode(child, path...); found != nil {
				return found
			}
		}

	case yaml.MappingNode:
		if len(path) > 0 {
			element := path[0]
			switch element.Type {
			case FieldPathType, MapPathType:
				v := element.Value.(string)

				// Content is a slice of pairs of key followed by value
				length := len(node.Content)
				for i := 0; i < length; i += 2 {
					key := node.Content[i]
					if (key.Kind == yaml.ScalarNode) && (key.Tag == "!!str") && (key.Value == v) {
						value := node.Content[i+1]
						found := FindYamlNode(value, path[1:]...)
						if found == value {
							// We want the location of the key, not the value
							return key
						}
						return found
					}
				}
			}
		}

	case yaml.SequenceNode:
		if len(path) > 0 {
			element := path[0]
			switch element.Type {
			case ListPathType:
				index := element.Value.(int)
				if index < len(node.Content) {
					return FindYamlNode(node.Content[index], path[1:]...)
				}
			}
		}
	}

	return node
}

func PrintYamlNodes(node *yaml.Node) {
	fmt.Printf("%d %s %s %d\n", node.Kind, node.Tag, node.Value, len(node.Content))
	for _, child := range node.Content {
		PrintYamlNodes(child)
	}
}
