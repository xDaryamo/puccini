package js

import (
	"errors"
	"fmt"
	"strings"

	"github.com/tliron/puccini/ard"
	"github.com/tliron/puccini/clout"
)

func CleanupScriptlet(scriptlet string) string {
	scriptlet = strings.Trim(scriptlet, " \t\n")
	scriptlet = strings.Replace(scriptlet, "\t", "  ", -1)
	return scriptlet
}

func GetScriptlet(name string, clout_ *clout.Clout) (string, error) {
	section, err := GetMetadataSection(name, clout_)
	if err != nil {
		return "", err
	}

	scriptlet, ok := section.(string)
	if !ok {
		return "", fmt.Errorf("scriptlet found in metadata but not a string: %s", name)
	}

	return scriptlet, nil
}

func GetScriptlets(name string, clout_ *clout.Clout) (ard.List, error) {
	section, err := GetMetadataSection(name, clout_)
	if err != nil {
		return nil, err
	}

	scriptlets, ok := section.(ard.StringMap)
	if !ok {
		return nil, fmt.Errorf("scriptlet metadata found but not a map: %s", name)
	}

	list := make(ard.List, 0, len(scriptlets))
	for _, scriptlet := range scriptlets {
		list = append(list, scriptlet)
	}

	return list, nil
}

func GetFunctionScriptlet(name string, clout_ *clout.Clout) (string, error) {
	metadata, err := GetMetadata(clout_)
	if err != nil {
		return "", err
	}

	functions, ok := metadata["functions"]
	if !ok {
		return "", errors.New("\"functions\" section not found in metadata")
	}

	m, ok := functions.(ard.StringMap)
	if !ok {
		return "", errors.New("malformed \"functions\" section in metadata")
	}

	function, ok := m[name]
	if !ok {
		return "", fmt.Errorf("function \"%s\" not found in metadata", name)
	}

	scriptlet, ok := function.(string)
	if !ok {
		return "", fmt.Errorf("function \"%s\" found in metadata but not a string", name)
	}

	return scriptlet, nil
}

func SetScriptlet(name string, scriptlet string, clout_ *clout.Clout) error {
	metadata, err := GetMetadata(clout_)
	if err != nil {
		return err
	}

	return ard.StringMapPutNested(metadata, name, scriptlet)
}

// Utils

func parseScriptletName(name string) ([]string, string, error) {
	segments := strings.Split(name, ".")
	last := len(segments) - 1
	if last == -1 {
		return nil, "", fmt.Errorf("malformed scriptlet name: %s", name)
	}

	return segments[:last], segments[last], nil
}
