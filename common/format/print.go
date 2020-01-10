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
		return PrintXMLDocument(xmlDocument, pretty)
	}

	switch format {
	case "yaml", "":
		return PrintYAML(data, pretty)
	case "json":
		return PrintJSON(data, pretty)
	case "xml":
		return PrintXML(data, pretty)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

func PrintYAML(data interface{}, pretty bool) error {
	indent := "          "
	if pretty {
		indent = terminal.Indent
	}
	return WriteYAML(data, terminal.Stdout, indent)
}

func PrintJSON(data interface{}, pretty bool) error {
	if pretty {
		bytes, err := prettyjsonFormatter.Marshal(data)
		if err != nil {
			return err
		}
		fmt.Fprintf(terminal.Stdout, "%s\n", bytes)
	} else {
		return WriteJSON(data, terminal.Stdout, "")
	}
	return nil
}

func PrintXML(data interface{}, pretty bool) error {
	indent := ""
	if pretty {
		indent = terminal.Indent
	}
	if err := WriteXML(data, terminal.Stdout, indent); err != nil {
		return err
	}
	if pretty {
		fmt.Fprintln(terminal.Stdout)
	}
	return nil
}

func PrintXMLDocument(xmlDocument *etree.Document, pretty bool) error {
	if pretty {
		xmlDocument.Indent(terminal.IndentSpaces)
	} else {
		xmlDocument.Indent(0)
	}
	_, err := xmlDocument.WriteTo(terminal.Stdout)
	return err
}
