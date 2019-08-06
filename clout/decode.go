package clout

import (
	"encoding/json"
	"encoding/xml"
	"io"

	"gopkg.in/yaml.v3"
)

func DecodeJson(reader io.Reader) (*Clout, error) {
	var c Clout
	var err error

	decoder := json.NewDecoder(reader)
	decoder.DisallowUnknownFields()

	if err = decoder.Decode(&c); err != nil {
		return nil, err
	}

	if err = c.Resolve(); err != nil {
		return nil, err
	}

	return &c, nil
}

func DecodeYaml(reader io.Reader) (*Clout, error) {
	var c Clout
	var err error

	decoder := yaml.NewDecoder(reader)
	decoder.KnownFields(true)

	if err = decoder.Decode(&c); err != nil {
		return nil, err
	}

	if err = c.Resolve(); err != nil {
		return nil, err
	}

	return &c, nil
}

func DecodeXml(reader io.Reader) (*Clout, error) {
	var c Clout
	var err error

	decoder := xml.NewDecoder(reader)

	if err = decoder.Decode(&c); err != nil {
		return nil, err
	}

	if err = c.Resolve(); err != nil {
		return nil, err
	}

	return &c, nil
}
