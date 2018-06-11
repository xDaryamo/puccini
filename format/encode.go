package format

import (
	"fmt"
	"strings"
)

func Encode(data interface{}, format string) (string, error) {
	switch format {
	case "json":
		return EncodeJson(data, Indent)
	case "yaml", "":
		return EncodeYaml(data)
	default:
		return "", fmt.Errorf("unsupported format: %s", format)
	}
}

func EncodeJson(data interface{}, indent string) (string, error) {
	var writer strings.Builder
	err := WriteJson(data, &writer, indent)
	if err != nil {
		return "", err
	}
	s := writer.String()
	if indent == "" {
		// json.Encoder adds a "\n", unlike json.Marshal
		s = strings.Trim(s, "\n")
	}
	return s, nil
}

func EncodeYaml(data interface{}) (string, error) {
	var writer strings.Builder
	err := WriteYaml(data, &writer)
	if err != nil {
		return "", err
	}
	return writer.String(), nil
}
