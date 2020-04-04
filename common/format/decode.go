package format

import (
	"fmt"
	"strings"

	"github.com/tliron/puccini/ard"
)

func Decode(code string, format string) (ard.Value, error) {
	switch format {
	case "yaml", "":
		return DecodeYAML(code)
	case "json":
		return DecodeJSON(code)
	default:
		return "", fmt.Errorf("unsupported format: %s", format)
	}
}

func DecodeYAML(code string) (ard.Value, error) {
	return ReadYAML(strings.NewReader(code))
}

func DecodeJSON(code string) (ard.Value, error) {
	return ReadJSON(strings.NewReader(code))
}
