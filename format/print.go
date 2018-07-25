package format

import (
	"fmt"
	"os"

	"github.com/beevik/etree"
	"github.com/hokaccha/go-prettyjson"
)

var prettyjsonFormatter = prettyjson.NewFormatter()

func init() {
	prettyjsonFormatter.Indent = IndentSpaces
}

func Print(data interface{}, format string, pretty bool) error {
	// Special handling for etree
	if xmlDocument, ok := data.(*etree.Document); ok {
		xmlDocument.Indent(IndentSpaces)
		_, err := xmlDocument.WriteTo(os.Stdout)
		return err
	}

	switch format {
	case "yaml", "":
		return PrintYaml(data, pretty)
	case "json":
		return PrintJson(data, pretty)
	case "xml":
		return PrintXml(data, pretty)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

func PrintYaml(data interface{}, pretty bool) error {
	return WriteYaml(data, os.Stdout)
}

func PrintJson(data interface{}, pretty bool) error {
	if pretty {
		bytes, err := prettyjsonFormatter.Marshal(data)
		if err != nil {
			return err
		}
		fmt.Printf("%s\n", bytes)
	} else {
		return WriteJson(data, os.Stdout, "")
	}
	return nil
}

func PrintXml(data interface{}, pretty bool) error {
	indent := ""
	if pretty {
		indent = Indent
	}
	err := WriteXml(data, os.Stdout, indent)
	if err != nil {
		return err
	}
	if pretty {
		fmt.Println()
	}
	return nil
}
