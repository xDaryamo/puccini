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
		return ValidateYAML(code)
	case "json":
		return ValidateJSON(code)
	case "xml":
		return ValidateJSON(code)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}
func ValidateYAML(code string) error {
	_, err := DecodeYAML(code)
	return err
}

func ValidateJSON(code string) error {
	return json.NewDecoder(strings.NewReader(code)).Decode(new(interface{}))
}

func ValidateXML(code string) error {
	return xml.NewDecoder(strings.NewReader(code)).Decode(new(interface{}))
}
