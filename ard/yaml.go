package ard

import (
	"fmt"
	"io"
	"strings"

	"gopkg.in/yaml.v3"
)

func FindYAMLNode(node *yaml.Node, path ...PathElement) *yaml.Node {
	if len(path) == 0 {
		return node
	}

	switch node.Kind {
	case yaml.AliasNode:
		return FindYAMLNode(node.Alias, path...)

	case yaml.DocumentNode:
		for _, childNode := range node.Content {
			// We assume it's a single document
			return FindYAMLNode(childNode, path...)
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

				if i+1 >= length {
					return keyNode
				}

				// Is it in one of the merged values?
				if (keyNode.Kind == yaml.ScalarNode) && (keyNode.Tag == "!!merge") {
					valueNode := node.Content[i+1]
					foundNode := FindYAMLNode(valueNode, path...)
					if foundNode != valueNode {
						return foundNode
					}
				}

				// We only support comparisons with string keys
				if (keyNode.Kind == yaml.ScalarNode) && (keyNode.Tag == "!!str") && (keyNode.Value == v) {
					valueNode := node.Content[i+1]
					foundNode := FindYAMLNode(valueNode, path[1:]...)
					if foundNode == valueNode {
						// We will use the key node for the location instead of the value node
						return keyNode
					}
					return foundNode
				}
			}
		}

	case yaml.SequenceNode:
		pathElement := path[0]
		switch pathElement.Type {
		case ListPathType:
			index := pathElement.Value.(int)
			if index < len(node.Content) {
				return FindYAMLNode(node.Content[index], path[1:]...)
			}
		}
	}

	// case yaml.ScalarNode:
	return node
}

// Write

var YAMLNodeKinds = map[yaml.Kind]string{
	yaml.DocumentNode: "Document",
	yaml.SequenceNode: "Sequence",
	yaml.MappingNode:  "Mapping",
	yaml.ScalarNode:   "Scalar",
	yaml.AliasNode:    "Alias",
}

func WriteYAMLNodes(writer io.Writer, node *yaml.Node) {
	WriteYAMLNode(writer, node, 0)
}

func WriteYAMLNode(writer io.Writer, node *yaml.Node, indent int) {
	s := ""

	s += strings.Repeat(" ", indent)

	s += YAMLNodeKinds[node.Kind]

	switch node.Kind {
	// Document and alias tag is always "", nothing to print
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
		WriteYAMLNode(writer, child, indent)
	}
}
