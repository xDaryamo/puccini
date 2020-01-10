package format

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/tliron/puccini/ard"
	"github.com/tliron/yamlkeys"
	"gopkg.in/yaml.v3"
)

func Read(reader io.Reader, format string) (interface{}, error) {
	switch format {
	case "yaml", "":
		return ReadYAML(reader)
	case "json":
		return ReadJSON(reader)
	default:
		return "", fmt.Errorf("unsupported format: %s", format)
	}
}

func ReadYAML(reader io.Reader) (interface{}, error) {
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

func ReadJSON(reader io.Reader) (interface{}, error) {
	var data interface{}
	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(&data); err == nil {
		data, _ = ard.ToMaps(data)
		return data, nil
	} else {
		return nil, err
	}
}
