package format

import (
	"fmt"

	"github.com/hokaccha/go-prettyjson"
)

var prettyjsonFormatter = prettyjson.NewFormatter()

func init() {
	prettyjsonFormatter.Indent = IndentSpaces
}

func Print(data interface{}, format string, pretty bool) error {
	switch format {
	case "json":
		return PrintJson(data, pretty)
	case "yaml", "":
		return PrintYaml(data, pretty)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

func PrintJson(data interface{}, pretty bool) error {
	if pretty {
		bytes, err := prettyjsonFormatter.Marshal(data)
		if err != nil {
			return err
		}
		fmt.Printf("%s\n", bytes)
	} else {
		s, err := EncodeJson(data, Indent)
		if err != nil {
			return err
		}
		fmt.Print(s)
	}
	return nil
}

func PrintYaml(data interface{}, pretty bool) error {
	s, err := EncodeYaml(data)
	if err != nil {
		return err
	}
	fmt.Print(s)
	return nil
}
