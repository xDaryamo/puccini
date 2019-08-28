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

func GetFunctionSoureCode(name string, c *clout.Clout) (string, error) {
	metadata, err := GetMetadata(c)
	if err != nil {
		return "", err
	}

	functions, ok := metadata["functions"]
	if !ok {
		return "", errors.New("functions not found")
	}
	m, ok := functions.(ard.Map)
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

func GetScriptSourceCode(name string, c *clout.Clout) (string, error) {
	section, err := GetScriptSection(name, c)
	if err != nil {
		return "", err
	}

	sourceCode, ok := section.(string)
	if !ok {
		return "", fmt.Errorf("source code found but not a string: %s", name)
	}

	return sourceCode, nil
}

func SetScriptSourceCode(name string, sourceCode string, c *clout.Clout) error {
	metadata, err := GetMetadata(c)
	if err != nil {
		return err
	}

	return SetMapNested(metadata, name, sourceCode)
}

func GetScripts(name string, c *clout.Clout) (ard.List, error) {
	section, err := GetScriptSection(name, c)
	if err != nil {
		return nil, err
	}

	sourceCodes, ok := section.(ard.Map)
	if !ok {
		return nil, fmt.Errorf("source code found but not a string: %s", name)
	}

	list := make(ard.List, 0, len(sourceCodes))
	for _, sourceCode := range sourceCodes {
		list = append(list, sourceCode)
	}

	return list, nil
}

func GetScriptSection(name string, c *clout.Clout) (interface{}, error) {
	segments, final, err := ParseScriptName(name)
	if err != nil {
		return nil, err
	}

	metadata, err := GetMetadata(c)
	if err != nil {
		return nil, err
	}

	m := metadata
	for _, s := range segments {
		o := m[s]
		var ok bool
		if m, ok = o.(ard.Map); !ok {
			return nil, fmt.Errorf("script not found: %s", name)
		}
	}

	o, ok := m[final]
	if !ok {
		return nil, fmt.Errorf("script not found: %s", name)
	}

	return o, nil
}

func GetMetadata(c *clout.Clout) (ard.Map, error) {
	metadata, ok := c.Metadata["puccini-js"]
	if !ok {
		return nil, errors.New("no scripts in Clout")
	}

	m, ok := metadata.(ard.Map)
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
