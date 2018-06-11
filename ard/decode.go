package ard

import (
	"encoding/json"
	"io"

	"gopkg.in/yaml.v2"
)

// TODO: row/column numbers for all parsed objects

func DecodeJson(reader io.Reader) (Map, error) {
	data := make(Map)
	decoder := json.NewDecoder(reader)
	err := decoder.Decode(&data)
	if err != nil {
		return nil, err
	}
	return data, err
}

func DecodeYaml(reader io.Reader) (Map, error) {
	data := make(map[string]interface{}) // *not* Map
	decoder := yaml.NewDecoder(reader)
	decoder.SetStrict(true)
	err := decoder.Decode(&data)
	if err != nil {
		return nil, err
	}
	return EnsureMap(data), nil
}
