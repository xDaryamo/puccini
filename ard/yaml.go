package ard

import (
	"fmt"
	"io"
	"strings"

	"gopkg.in/yaml.v3"
)

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

// Utils

func FindYamlNode(node *yaml.Node, path ...PathElement) *yaml.Node {
	if len(path) == 0 {
		return node
	}

	switch node.Kind {
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
					valueNode := node.Content[i+1]
					foundNode := FindYamlNode(valueNode, path[1:]...)
					if foundNode == valueNode {
						// We want the location of the key node, not the value node
						return keyNode
					}
					return foundNode
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

func PrintYamlNodes(writer io.Writer, node *yaml.Node) {
	PrintYamlNode(writer, node, 0)
}

var YamlNodeKinds = map[yaml.Kind]string{
	yaml.DocumentNode: "Document",
	yaml.SequenceNode: "Sequence",
	yaml.MappingNode:  "Mapping",
	yaml.ScalarNode:   "Scalar",
	yaml.AliasNode:    "Alias",
}

func PrintYamlNode(writer io.Writer, node *yaml.Node, indent int) {
	s := ""

	s += strings.Repeat(" ", indent)

	s += YamlNodeKinds[node.Kind]

	switch node.Kind {
	// Document tag is always "", nothing to print
	// Sequence tag is always "!!seq", no need to print
	// Mapping tag is always "!!map", no need to print

	case yaml.ScalarNode:
		s += " "
		s += node.Tag
	}

	if node.Value != "" {
		s += " "
		s += node.Value
	}

	fmt.Fprintln(writer, s)

	indent += 1
	for _, child := range node.Content {
		PrintYamlNode(writer, child, indent)
	}
}
