package js

import (
	"errors"
	"fmt"
	"strings"
	"unicode"

	"github.com/tliron/puccini/ard"
	"github.com/tliron/puccini/clout"
)

func Cleanup(script string) string {
	script = strings.Trim(script, " \n")
	script = strings.Replace(script, "\t", "  ", -1)
	return script
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

func GetFunctionSoureCode(name string, clout_ *clout.Clout) (string, error) {
	metadata, err := GetMetadata(clout_)
	if err != nil {
		return "", err
	}

	functions, ok := metadata["functions"]
	if !ok {
		return "", errors.New("functions not found")
	}
	m, ok := functions.(ard.StringMap)
	if !ok {
		return "", errors.New("malformed functions section")
	}

	function, ok := m[name]
	if !ok {
		return "", fmt.Errorf("function \"%s\" not found", name)
	}

	script, ok := function.(string)
	if !ok {
		return "", fmt.Errorf("function \"%s\" found but not a string", name)
	}

	return script, nil
}

func GetScriptSourceCode(name string, clout_ *clout.Clout) (string, error) {
	section, err := GetScriptSection(name, clout_)
	if err != nil {
		return "", err
	}

	sourceCode, ok := section.(string)
	if !ok {
		return "", fmt.Errorf("source code found but not a string: %s", name)
	}

	return sourceCode, nil
}

func SetScriptSourceCode(name string, sourceCode string, clout_ *clout.Clout) error {
	metadata, err := GetMetadata(clout_)
	if err != nil {
		return err
	}

	return SetMapNested(metadata, name, sourceCode)
}

func GetScripts(name string, clout_ *clout.Clout) (ard.List, error) {
	section, err := GetScriptSection(name, clout_)
	if err != nil {
		return nil, err
	}

	sourceCodes, ok := section.(ard.StringMap)
	if !ok {
		return nil, fmt.Errorf("source code found but not a string: %s", name)
	}

	list := make(ard.List, 0, len(sourceCodes))
	for _, sourceCode := range sourceCodes {
		list = append(list, sourceCode)
	}

	return list, nil
}

func GetScriptSection(name string, clout_ *clout.Clout) (interface{}, error) {
	segments, final, err := ParseScriptName(name)
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
			return nil, fmt.Errorf("script not found: %s", name)
		}
	}

	o, ok := m[final]
	if !ok {
		return nil, fmt.Errorf("script not found: %s", name)
	}

	return o, nil
}

func GetMetadata(clout_ *clout.Clout) (ard.StringMap, error) {
	metadata, ok := clout_.Metadata["puccini-js"]
	if !ok {
		return nil, errors.New("no scripts in Clout")
	}

	m, ok := metadata.(ard.StringMap)
	if !ok {
		return nil, errors.New("malformed \"puccini-js\" metadata in Clout")
	}

	return m, nil
}

func ParseScriptName(name string) ([]string, string, error) {
	segments := strings.Split(name, ".")
	last := len(segments) - 1
	if last == -1 {
		return nil, "", fmt.Errorf("malformed script name: %s", name)
	}

	return segments[:last], segments[last], nil
}
