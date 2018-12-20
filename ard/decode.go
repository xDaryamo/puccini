package ard

import (
	"encoding/json"
	"encoding/xml"
	"io"

	"gopkg.in/yaml.v2"
)

// TODO: row/column numbers for all parsed objects

func DecodeJson(reader io.Reader) (Map, error) {
	data := make(Map)
	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(&data); err != nil {
		return nil, err
	}
	return data, nil
}

func DecodeYaml(reader io.Reader) (Map, error) {
	data := make(map[string]interface{}) // *not* Map
	decoder := yaml.NewDecoder(reader)
	decoder.SetStrict(true)
	if err := decoder.Decode(&data); err != nil {
		return nil, err
	}
	return EnsureMap(data), nil
}

func DecodeXml(reader io.Reader) (Map, error) {
	data := make(Map)
	decoder := xml.NewDecoder(reader)
	if err := decoder.Decode(&data); err != nil {
		return nil, err
	}
	return data, nil
}
