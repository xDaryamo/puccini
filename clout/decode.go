package clout

import (
	"encoding/json"
	"io"

	"gopkg.in/yaml.v2"
)

func DecodeJson(reader io.Reader) (*Clout, error) {
	var c Clout

	decoder := json.NewDecoder(reader)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&c)
	if err != nil {
		return nil, err
	}

	err = c.Resolve()
	if err != nil {
		return nil, err
	}

	return &c, nil
}

func DecodeYaml(reader io.Reader) (*Clout, error) {
	var c Clout

	decoder := yaml.NewDecoder(reader)
	decoder.SetStrict(true)

	err := decoder.Decode(&c)
	if err != nil {
		return nil, err
	}

	err = c.Resolve()
	if err != nil {
		return nil, err
	}

	return &c, nil
}
