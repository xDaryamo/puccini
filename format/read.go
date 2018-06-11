package format

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/tliron/puccini/ard"
	"gopkg.in/yaml.v2"
)

func Read(reader io.Reader, format string) (interface{}, error) {
	switch format {
	case "json":
		return ReadJson(reader)
	case "yaml", "":
		return ReadYaml(reader)
	default:
		return "", fmt.Errorf("unsupported format: %s", format)
	}
}

func ReadJson(reader io.Reader) (interface{}, error) {
	var data interface{}
	decoder := json.NewDecoder(reader)
	err := decoder.Decode(&data)
	if err != nil {
		return nil, err
	}
	return data, err
}

func ReadYaml(reader io.Reader) (interface{}, error) {
	var data interface{}
	decoder := yaml.NewDecoder(reader)
	decoder.SetStrict(true)
	err := decoder.Decode(&data)
	if err != nil {
		return nil, err
	}
	data, _ = ard.EnsureValue(data)
	return data, nil
}
