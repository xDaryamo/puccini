package ard

import (
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
		if len(node.Content) > 0 {
			// Length *should* be 1
			return FindYAMLNode(node.Content[0], path...)
		}

	case yaml.MappingNode:
		pathElement := path[0]
		switch pathElement.Type {
		case FieldPathType, MapPathType:
			value := pathElement.Value.(string)

			// Content is a slice of pairs of key-followed-by-value
			length := len(node.Content)
			for i := 0; i < length; i += 2 {
				keyNode := node.Content[i]

				if i+1 >= length {
					// Length *should* be an even number
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
				if (keyNode.Kind == yaml.ScalarNode) && (keyNode.Tag == "!!str") && (keyNode.Value == value) {
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

	return node
}
