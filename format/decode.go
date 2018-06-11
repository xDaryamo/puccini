package format

import (
	"fmt"
	"strings"
)

func Decode(code string, format string) (interface{}, error) {
	switch format {
	case "json":
		return DecodeJson(code)
	case "yaml", "":
		return DecodeYaml(code)
	default:
		return "", fmt.Errorf("unsupported format: %s", format)
	}
}

func DecodeJson(code string) (interface{}, error) {
	return ReadJson(strings.NewReader(code))
}

func DecodeYaml(code string) (interface{}, error) {
	return ReadYaml(strings.NewReader(code))
}
