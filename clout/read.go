package clout

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"

	"github.com/tliron/puccini/ard"
	"github.com/tliron/yamlkeys"
)

func Read(reader io.Reader, format string) (*Clout, error) {
	switch format {
	case "yaml", "":
		return ReadYAML(reader)
	case "json":
		return ReadJSON(reader)
	case "xml":
		return ReadXML(reader)
	default:
		return nil, fmt.Errorf("unsupported format: %s", format)
	}
}

func ReadYAML(reader io.Reader) (*Clout, error) {
	var err error
	var ok bool

	var data interface{}
	if data, err = yamlkeys.Decode(reader); err != nil {
		return nil, err
	}

	var map_ ard.Map
	if map_, ok = data.(ard.Map); !ok {
		return nil, fmt.Errorf("not a map: %T", data)
	}

	var clout *Clout
	if clout, err = Decode(map_); err != nil {
		return nil, err
	}

	if err = clout.Resolve(); err != nil {
		return nil, err
	}

	return clout, nil
}

func ReadJSON(reader io.Reader) (*Clout, error) {
	var clout Clout
	var err error

	decoder := json.NewDecoder(reader)
	decoder.DisallowUnknownFields()

	if err = decoder.Decode(&clout); err != nil {
		return nil, err
	}

	if err = clout.Resolve(); err != nil {
		return nil, err
	}

	return &clout, nil
}

func ReadXML(reader io.Reader) (*Clout, error) {
	var clout Clout
	var err error

	decoder := xml.NewDecoder(reader)

	if err = decoder.Decode(&clout); err != nil {
		return nil, err
	}

	if err = clout.Resolve(); err != nil {
		return nil, err
	}

	return &clout, nil
}
