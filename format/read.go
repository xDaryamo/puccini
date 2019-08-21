package format

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"

	"github.com/tliron/puccini/ard"
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
	var data interface{}
	decoder := yaml.NewDecoder(reader)
	if err := decoder.Decode(&data); err != nil {
		return nil, err
	}
	data, _ = ard.EnsureValue(data)
	return data, nil
}

func ReadJson(reader io.Reader) (interface{}, error) {
	var data interface{}
	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(&data); err != nil {
		return nil, err
	}
	return data, nil
}

func ReadXml(reader io.Reader) (interface{}, error) {
	var data interface{}
	decoder := xml.NewDecoder(reader)
	if err := decoder.Decode(&data); err != nil {
		return nil, err
	}
	return data, nil
}
