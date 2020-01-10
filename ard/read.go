package ard

import (
	"encoding/json"
	"fmt"
	"io"

	"github.com/tliron/puccini/url"
	"github.com/tliron/yamlkeys"
	"gopkg.in/yaml.v3"
)

func ReadFromURL(url_ url.URL, locate bool) (Map, Locator, error) {
	reader, err := url_.Open()
	if err != nil {
		return nil, nil, err
	}

	if readerCloser, ok := reader.(io.ReadCloser); ok {
		defer readerCloser.Close()
	}

	return Read(reader, url_.Format(), locate)
}

func Read(reader io.Reader, format string, locate bool) (Map, Locator, error) {
	switch format {
	case "yaml", "":
		return ReadYAML(reader, locate)
	case "json":
		return ReadJSON(reader, locate)
	default:
		return nil, nil, fmt.Errorf("unsupported format: \"%s\"", format)
	}
}

func ReadYAML(reader io.Reader, locate bool) (Map, Locator, error) {
	var data Map
	var locator Locator
	var node yaml.Node

	decoder := yaml.NewDecoder(reader)
	if err := decoder.Decode(&node); err == nil {
		if decoded, err := yamlkeys.DecodeNode(&node); err == nil {
			var ok bool
			if data, ok = decoded.(Map); ok {
				if locate {
					locator = NewYAMLLocator(&node)
				}
			} else {
				return nil, nil, fmt.Errorf("YAML content is a \"%T\" instead of a map", decoded)
			}
		} else {
			return nil, nil, err
		}
	} else {
		return nil, nil, err
	}

	// We do not need to call EnsureMaps because yamlkeys takes care of it
	return data, locator, nil
}

func ReadJSON(reader io.Reader, locate bool) (Map, Locator, error) {
	data := make(Map)
	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(&data); err == nil {
		return EnsureMaps(data), nil, nil
	} else {
		return nil, nil, err
	}
}
