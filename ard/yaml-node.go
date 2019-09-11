package ard

import (
	"fmt"
	"io"
	"strings"

	"gopkg.in/yaml.v3"
)

var YamlNodeKinds = map[yaml.Kind]string{
	yaml.DocumentNode: "Document",
	yaml.SequenceNode: "Sequence",
	yaml.MappingNode:  "Mapping",
	yaml.ScalarNode:   "Scalar",
	yaml.AliasNode:    "Alias",
}

func DecodeYamlNode(node *yaml.Node) (interface{}, error) {
	switch node.Kind {
	case yaml.AliasNode:
		return DecodeYamlNode(node.Alias)

	case yaml.DocumentNode:
		var slice []interface{}

		for _, childNode := range node.Content {
			if value, err := DecodeYamlNode(childNode); err == nil {
				slice = append(slice, value)
			} else {
				return nil, err
			}
		}

		switch len(slice) {
		case 1:
			// Single document
			return slice[0], nil
		case 0:
			// Empty
			return make(Map), nil
		default:
			// Multiple documents
			return slice, nil
		}

	case yaml.MappingNode:
		map_ := make(Map)

		// Content is a slice of pairs of key followed by value
		length := len(node.Content)
		if length%2 != 0 {
			panic("malformed YAML map")
		}

		for i := 0; i < length; i += 2 {
			keyNode := node.Content[i]
			valueNode := node.Content[i+1]

			if value, err := DecodeYamlNode(valueNode); err == nil {
				if (keyNode.Kind == yaml.ScalarNode) && (keyNode.Tag == "!!merge") {
					// See: https://yaml.org/type/merge.html
					switch value.(type) {
					case Map:
						MapMerge(map_, value.(Map), false)

					case List:
						for _, v := range value.(List) {
							if m, ok := v.(Map); ok {
								MapMerge(map_, m, false)
							} else {
								PanicMalformedMerge(keyNode)
							}
						}

					default:
						PanicMalformedMerge(keyNode)
					}
				} else {
					if key, keyData, err := DecodeYamlKeyNode(keyNode); err == nil {
						if keyData == nil {
							if _, ok := map_[key]; ok {
								return nil, ErrorDuplicateKey(keyNode, key)
							}
						} else {
							for k, _ := range map_ {
								if Equals(keyData, KeyData(k)) {
									return nil, ErrorDuplicateKey(keyNode, key)
								}
							}
						}
						map_[key] = value
					} else {
						return nil, err
					}
				}
			} else {
				return nil, err
			}
		}

		return map_, nil

	case yaml.SequenceNode:
		var slice []interface{}

		for _, childNode := range node.Content {
			if value, err := DecodeYamlNode(childNode); err == nil {
				slice = append(slice, value)
			} else {
				return nil, err
			}
		}

		return slice, nil

	case yaml.ScalarNode:
		var value interface{}
		if err := node.Decode(&value); err == nil {
			return value, nil
		} else {
			return nil, err
		}
	}

	panic("malformed YAML node")
}

func DecodeYamlKeyNode(node *yaml.Node) (interface{}, interface{}, error) {
	// Workaround for gopkg.in/yaml.v3 not supporting the decoding of complex keys
	// See: https://github.com/go-yaml/yaml/issues/502

	if data, err := DecodeYamlNode(node); err == nil {
		if IsBasicType(data) {
			return data, nil, nil
		} else {
			// A pointer is hashable but its *contents* are not compared, so we will be losing
			// the ability to test for uniqueness at this stage
			key, err := NewYamlKey(data)
			return key, data, err
		}
	} else {
		return nil, nil, err
	}
}

func FindYamlNode(node *yaml.Node, path ...PathElement) *yaml.Node {
	if len(path) == 0 {
		return node
	}

	switch node.Kind {
	case yaml.AliasNode:
		return FindYamlNode(node.Alias, path...)

	case yaml.DocumentNode:
		for _, childNode := range node.Content {
			// We assume it's a single document
			return FindYamlNode(childNode, path...)
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
					foundNode := FindYamlNode(valueNode, path...)
					if foundNode != valueNode {
						return foundNode
					}
				}

				// We only support comparisons with string keys
				if (keyNode.Kind == yaml.ScalarNode) && (keyNode.Tag == "!!str") && (keyNode.Value == v) {
					valueNode := node.Content[i+1]
					foundNode := FindYamlNode(valueNode, path[1:]...)
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
				return FindYamlNode(node.Content[index], path[1:]...)
			}
		}
	}

	return node
}

// Write

func WriteYamlNodes(writer io.Writer, node *yaml.Node) {
	WriteYamlNode(writer, node, 0)
}

func WriteYamlNode(writer io.Writer, node *yaml.Node, indent int) {
	s := ""

	s += strings.Repeat(" ", indent)

	s += YamlNodeKinds[node.Kind]

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
		WriteYamlNode(writer, child, indent)
	}
}

// Utils

func IsBasicType(data interface{}) bool {
	switch data.(type) {
	case bool, string, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, uintptr, float32, float64, complex64, complex128:
		return true
	}
	return false
}

func ErrorDuplicateKey(node *yaml.Node, key interface{}) error {
	return fmt.Errorf("malformed YAML @%d,%d: duplicate map key: %s", node.Line, node.Column, key)
}

func PanicMalformedMerge(node *yaml.Node) {
	panic(fmt.Sprintf("malformed YAML @%d,%d: merge", node.Line, node.Column))
}
