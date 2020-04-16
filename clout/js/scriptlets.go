package js

import (
	"fmt"
	"strings"

	"github.com/tliron/puccini/ard"
	cloutpkg "github.com/tliron/puccini/clout"
)

func CleanupScriptlet(scriptlet string) string {
	scriptlet = strings.Trim(scriptlet, " \t\n")
	scriptlet = strings.Replace(scriptlet, "\t", "  ", -1)
	return scriptlet
}

func GetScriptlet(name string, clout *cloutpkg.Clout) (string, error) {
	section, err := GetMetadataSection(name, clout)
	if err != nil {
		return "", err
	}

	scriptlet, ok := section.(string)
	if !ok {
		return "", fmt.Errorf("scriptlet found in metadata but not a string: %s", name)
	}

	return scriptlet, nil
}

func SetScriptlet(name string, scriptlet string, clout *cloutpkg.Clout) error {
	metadata, err := GetMetadata(clout)
	if err != nil {
		return err
	}

	return ard.StringMapPutNested(metadata, name, scriptlet)
}

func GetScriptletNames(baseName string, clout *cloutpkg.Clout) ([]string, error) {
	section, err := GetMetadataSection(baseName, clout)
	if err != nil {
		return nil, err
	}

	scriptlets, ok := section.(ard.StringMap)
	if !ok {
		return nil, fmt.Errorf("scriptlet metadata found but not a map: %s", baseName)
	}

	list := make([]string, 0, len(scriptlets))
	for name, _ := range scriptlets {
		list = append(list, fmt.Sprintf("%s.%s", baseName, name))
	}

	return list, nil
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
