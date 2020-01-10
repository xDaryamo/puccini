package format

import (
	"fmt"

	"github.com/beevik/etree"
	"github.com/hokaccha/go-prettyjson"
	"github.com/tliron/puccini/common/terminal"
)

var prettyjsonFormatter = prettyjson.NewFormatter()

func init() {
	prettyjsonFormatter.Indent = terminal.IndentSpaces
}

func Print(data interface{}, format string, pretty bool) error {
	// Special handling for strings
	if s, ok := data.(string); ok {
		if pretty {
			s += "\n"
		}
		_, err := fmt.Fprint(terminal.Stdout, s)
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
		indent = terminal.Indent
	}
	return WriteYaml(data, terminal.Stdout, indent)
}

func PrintJson(data interface{}, pretty bool) error {
	if pretty {
		bytes, err := prettyjsonFormatter.Marshal(data)
		if err != nil {
			return err
		}
		fmt.Fprintf(terminal.Stdout, "%s\n", bytes)
	} else {
		return WriteJson(data, terminal.Stdout, "")
	}
	return nil
}

func PrintXml(data interface{}, pretty bool) error {
	indent := ""
	if pretty {
		indent = terminal.Indent
	}
	if err := WriteXml(data, terminal.Stdout, indent); err != nil {
		return err
	}
	if pretty {
		fmt.Fprintln(terminal.Stdout)
	}
	return nil
}

func PrintXmlDocument(xmlDocument *etree.Document, pretty bool) error {
	if pretty {
		xmlDocument.Indent(terminal.IndentSpaces)
	} else {
		xmlDocument.Indent(0)
	}
	_, err := xmlDocument.WriteTo(terminal.Stdout)
	return err
}
