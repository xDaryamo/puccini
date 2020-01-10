package js

import (
	"errors"
	"fmt"

	"github.com/tliron/puccini/ard"
	"github.com/tliron/puccini/clout"
)

func GetMetadata(clout_ *clout.Clout) (ard.StringMap, error) {
	metadata, ok := clout_.Metadata["puccini-js"]
	if !ok {
		return nil, errors.New("no \"puccini-js\" metadata in Clout")
	}

	m, ok := metadata.(ard.StringMap)
	if !ok {
		return nil, errors.New("malformed \"puccini-js\" metadata in Clout")
	}

	return m, nil
}

func GetMetadataSection(name string, clout_ *clout.Clout) (interface{}, error) {
	segments, final, err := parseScriptletName(name)
	if err != nil {
		return nil, err
	}

	metadata, err := GetMetadata(clout_)
	if err != nil {
		return nil, err
	}

	m := metadata
	for _, s := range segments {
		o := m[s]
		var ok bool
		if m, ok = o.(ard.StringMap); !ok {
			return nil, fmt.Errorf("scriptlet metadata not found: %s", name)
		}
	}

	section, ok := m[final]
	if !ok {
		return nil, fmt.Errorf("scriptlet metadata not found: %s", name)
	}

	return section, nil
}
