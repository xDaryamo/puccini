package ard

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"io"
	"io/ioutil"

	"gopkg.in/yaml.v3"
)

func DecodeJson(reader io.Reader) (Map, Locator, error) {
	data := make(Map)
	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(&data); err != nil {
		return nil, nil, err
	}
	return data, nil, nil
}

func DecodeYaml(reader io.Reader) (Map, Locator, error) {
	data := make(map[string]interface{}) // *not* Map
	var locator Locator

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
		decoder.KnownFields(true)
		if err := decoder.Decode(&data); err != nil {
			return nil, nil, err
		}
	} else {
		return nil, nil, err
	}

	return EnsureMap(data), locator, nil
}

func DecodeXml(reader io.Reader) (Map, Locator, error) {
	data := make(Map)
	decoder := xml.NewDecoder(reader)
	if err := decoder.Decode(&data); err != nil {
		return nil, nil, err
	}
	return data, nil, nil
}
