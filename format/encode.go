package format

import (
	"fmt"
	"strings"
)

func Encode(data interface{}, format string) (string, error) {
	switch format {
	case "yaml", "":
		return EncodeYaml(data, Indent)
	case "json":
		return EncodeJson(data, Indent)
	case "xml":
		return EncodeXml(data, Indent)
	default:
		return "", fmt.Errorf("unsupported format: %s", format)
	}
}

func EncodeYaml(data interface{}, indent string) (string, error) {
	var writer strings.Builder
	if err := WriteYaml(data, &writer, indent); err == nil {
		return writer.String(), nil
	} else {
		return "", err
	}
}

func EncodeJson(data interface{}, indent string) (string, error) {
	var writer strings.Builder
	if err := WriteJson(data, &writer, indent); err == nil {
		s := writer.String()
		if indent == "" {
			// json.Encoder adds a "\n", unlike json.Marshal
			s = strings.Trim(s, "\n")
		}
		return s, nil
	} else {
		return "", err
	}
}

func EncodeXml(data interface{}, indent string) (string, error) {
	var writer strings.Builder
	if err := WriteXml(data, &writer, indent); err == nil {
		return writer.String(), nil
	} else {
		return "", err
	}
}
