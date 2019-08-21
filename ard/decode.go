package ard

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"strings"

	"gopkg.in/yaml.v3"
)

func DecodeJson(reader io.Reader, locate bool) (Map, Locator, error) {
	data := make(Map)
	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(&data); err != nil {
		return nil, nil, err
	}
	return data, nil, nil
}

func DecodeYaml(reader io.Reader, locate bool) (Map, Locator, error) {
	data := make(map[string]interface{}) // *not* Map
	var locator Locator

	if locate {
		// We need to read all into a buffer in order to both unmarshal and decode
		if buffer, err := ioutil.ReadAll(reader); err == nil {
			// Unmarshal node
			var node yaml.Node
			if err := yaml.Unmarshal(buffer, &node); err == nil {
				//PrintYamlNodes(os.Stdout, &node)
				locator = NewYamlLocator(&node)
			} else {
				return nil, nil, err
			}

			// Decode
			decoder := yaml.NewDecoder(bytes.NewReader(buffer))
			if err := decoder.Decode(&data); err != nil {
				return nil, nil, err
			}
		} else {
			return nil, nil, err
		}
	} else {
		decoder := yaml.NewDecoder(reader)
		if err := decoder.Decode(&data); err != nil {
			return nil, nil, err
		}
	}

	return EnsureMap(data), locator, nil
}

func DecodeXml(reader io.Reader, locate bool) (Map, Locator, error) {
	data := make(Map)
	decoder := xml.NewDecoder(reader)
	if err := decoder.Decode(&data); err != nil {
		return nil, nil, err
	}
	return data, nil, nil
}

// Utils

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
		PrintYamlNode(writer, child, indent)
	}
}
