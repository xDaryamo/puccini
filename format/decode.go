package format

import (
	"fmt"
	"strings"
)

func Decode(code string, format string) (interface{}, error) {
	switch format {
	case "yaml", "":
		return DecodeYaml(code)
	case "json":
		return DecodeJson(code)
	default:
		return "", fmt.Errorf("unsupported format: %s", format)
	}
}

func DecodeYaml(code string) (interface{}, error) {
	return ReadYaml(strings.NewReader(code))
}

func DecodeJson(code string) (interface{}, error) {
	return ReadJson(strings.NewReader(code))
}
