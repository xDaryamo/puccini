package js

import (
	"errors"
	"fmt"
	"strings"
	"unicode"

	"github.com/tliron/puccini/ard"
	"github.com/tliron/puccini/clout"
)

func Cleanup(scriptlet string) string {
	scriptlet = strings.Trim(scriptlet, " \n")
	scriptlet = strings.Replace(scriptlet, "\t", "  ", -1)
	return scriptlet
}

func ToJavaScriptStyle(name string) string {
	runes := []rune(name)
	length := len(runes)
	if (length > 0) && unicode.IsUpper(runes[0]) {
		if (length > 1) && unicode.IsUpper(runes[1]) {
			// If the second rune is also uppercase we'll keep the name as is
			return name
		}
		r := make([]rune, 1, length-1)
		r[0] = unicode.ToLower(runes[0])
		return string(append(r, runes[1:]...))
	}
	return name
}

func GetFunctionScriptlet(name string, clout_ *clout.Clout) (string, error) {
	metadata, err := GetJavaScriptMetadata(clout_)
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

func GetScriptlet(name string, clout_ *clout.Clout) (string, error) {
	section, err := GetScriptletMetadata(name, clout_)
	if err != nil {
		return "", err
	}

	scriptlet, ok := section.(string)
	if !ok {
		return "", fmt.Errorf("scriptlet found in metadata but not a string: %s", name)
	}

	return scriptlet, nil
}

func SetScriptlet(name string, scriptlet string, clout_ *clout.Clout) error {
	metadata, err := GetJavaScriptMetadata(clout_)
	if err != nil {
		return err
	}

	return SetMapNested(metadata, name, scriptlet)
}

func GetScriptlets(name string, clout_ *clout.Clout) (ard.List, error) {
	section, err := GetScriptletMetadata(name, clout_)
	if err != nil {
		return nil, err
	}

	scriptlets, ok := section.(ard.StringMap)
	if !ok {
		return nil, fmt.Errorf("scriptlet metadata found but not a map: %s", name)
	}

	list := make(ard.List, 0, len(scriptlets))
	for _, sourceCode := range scriptlets {
		list = append(list, sourceCode)
	}

	return list, nil
}

func GetScriptletMetadata(name string, clout_ *clout.Clout) (interface{}, error) {
	segments, final, err := ParseScriptletName(name)
	if err != nil {
		return nil, err
	}

	metadata, err := GetJavaScriptMetadata(clout_)
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

	o, ok := m[final]
	if !ok {
		return nil, fmt.Errorf("scriptlet metadata not found: %s", name)
	}

	return o, nil
}

func GetJavaScriptMetadata(clout_ *clout.Clout) (ard.StringMap, error) {
	metadata, ok := clout_.Metadata["puccini-js"]
	if !ok {
		return nil, errors.New("no scriptlets in Clout")
	}

	m, ok := metadata.(ard.StringMap)
	if !ok {
		return nil, errors.New("malformed \"puccini-js\" metadata in Clout")
	}

	return m, nil
}

func ParseScriptletName(name string) ([]string, string, error) {
	segments := strings.Split(name, ".")
	last := len(segments) - 1
	if last == -1 {
		return nil, "", fmt.Errorf("malformed scriptlet name: %s", name)
	}

	return segments[:last], segments[last], nil
}
