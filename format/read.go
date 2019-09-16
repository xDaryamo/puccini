package format

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"

	"github.com/tliron/puccini/ard"
	"github.com/tliron/yamlkeys"
	"gopkg.in/yaml.v3"
)

func Read(reader io.Reader, format string) (interface{}, error) {
	switch format {
	case "yaml", "":
		return ReadYaml(reader)
	case "json":
		return ReadJson(reader)
	case "xml":
		return ReadXml(reader)
	default:
		return "", fmt.Errorf("unsupported format: %s", format)
	}
}

func ReadYaml(reader io.Reader) (interface{}, error) {
	var node yaml.Node
	decoder := yaml.NewDecoder(reader)
	if err := decoder.Decode(&node); err == nil {
		if data, err := yamlkeys.DecodeNode(&node); err == nil {
			return data, nil
		} else {
			return nil, err
		}
	} else {
		return nil, err
	}
}

func ReadJson(reader io.Reader) (interface{}, error) {
	var data interface{}
	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(&data); err == nil {
		data, _ = ard.ToMaps(data)
		return data, nil
	} else {
		return nil, err
	}
}

func ReadXml(reader io.Reader) (interface{}, error) {
	var data interface{}
	decoder := xml.NewDecoder(reader)
	if err := decoder.Decode(&data); err == nil {
		data, _ = ard.ToMaps(data)
		return data, nil
	} else {
		return nil, err
	}
}
