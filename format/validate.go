package format

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"strings"
)

func Validate(code string, format string) error {
	switch format {
	case "yaml", "":
		return ValidateYaml(code)
	case "json":
		return ValidateJson(code)
	case "xml":
		return ValidateJson(code)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}
func ValidateYaml(code string) error {
	_, err := DecodeYaml(code)
	return err
}

func ValidateJson(code string) error {
	return json.NewDecoder(strings.NewReader(code)).Decode(new(interface{}))
}

func ValidateXml(code string) error {
	return xml.NewDecoder(strings.NewReader(code)).Decode(new(interface{}))
}
