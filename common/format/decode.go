package format

import (
	"fmt"
	"strings"
)

func Decode(code string, format string) (interface{}, error) {
	switch format {
	case "yaml", "":
		return DecodeYAML(code)
	case "json":
		return DecodeJSON(code)
	default:
		return "", fmt.Errorf("unsupported format: %s", format)
	}
}

func DecodeYAML(code string) (interface{}, error) {
	return ReadYAML(strings.NewReader(code))
}

func DecodeJSON(code string) (interface{}, error) {
	return ReadJSON(strings.NewReader(code))
}
