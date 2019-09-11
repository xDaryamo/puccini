package format

import (
	"fmt"

	"github.com/beevik/etree"
	"github.com/hokaccha/go-prettyjson"
)

var prettyjsonFormatter = prettyjson.NewFormatter()

func init() {
	prettyjsonFormatter.Indent = IndentSpaces
}

func Print(data interface{}, format string, pretty bool) error {
	// Special handling for strings
	if s, ok := data.(string); ok {
		if pretty {
			s += "\n"
		}
		_, err := fmt.Fprint(Stdout, s)
		return err
	}

	// Special handling for etree
	if xmlDocument, ok := data.(*etree.Document); ok {
		return PrintXmlDocument(xmlDocument, pretty)
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
	indent := "          "
	if pretty {
		indent = Indent
	}
	return WriteYaml(data, Stdout, indent)
}

func PrintJson(data interface{}, pretty bool) error {
	if pretty {
		bytes, err := prettyjsonFormatter.Marshal(data)
		if err != nil {
			return err
		}
		fmt.Fprintf(Stdout, "%s\n", bytes)
	} else {
		return WriteJson(data, Stdout, "")
	}
	return nil
}

func PrintXml(data interface{}, pretty bool) error {
	indent := ""
	if pretty {
		indent = Indent
	}
	if err := WriteXml(data, Stdout, indent); err != nil {
		return err
	}
	if pretty {
		fmt.Fprintln(Stdout)
	}
	return nil
}

func PrintXmlDocument(xmlDocument *etree.Document, pretty bool) error {
	if pretty {
		xmlDocument.Indent(IndentSpaces)
	} else {
		xmlDocument.Indent(0)
	}
	_, err := xmlDocument.WriteTo(Stdout)
	return err
}
