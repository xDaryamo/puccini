package format

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"

	"github.com/beevik/etree"
	"gopkg.in/yaml.v3"
)

func Write(data interface{}, format string, indent string, writer io.Writer) error {
	// Special handling for strings
	if s, ok := data.(string); ok {
		_, err := io.WriteString(writer, s)
		return err
	}

	// Special handling for etree
	if xmlDocument, ok := data.(*etree.Document); ok {
		return WriteXmlDocument(xmlDocument, writer, indent)
	}

	switch format {
	case "yaml", "":
		return WriteYaml(data, writer, indent)
	case "json":
		return WriteJson(data, writer, indent)
	case "xml":
		return WriteXml(data, writer, indent)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

func WriteYaml(data interface{}, writer io.Writer, indent string) error {
	encoder := yaml.NewEncoder(writer)
	// BUG: currently does not allow a value of 1, see: https://github.com/go-yaml/yaml/issues/501
	encoder.SetIndent(len(indent)) // This might not work as expected for tabs!
	if slice, ok := data.([]interface{}); ok {
		// YAML separates each entry with "---"
		// (In JSON the slice would be written as an array)
		for _, d := range slice {
			if err := encoder.Encode(d); err != nil {
				return err
			}
		}
		return nil
	} else {
		return encoder.Encode(data)
	}
}

func WriteJson(data interface{}, writer io.Writer, indent string) error {
	encoder := json.NewEncoder(writer)
	encoder.SetIndent("", indent)
	return encoder.Encode(data)
}

func WriteXml(data interface{}, writer io.Writer, indent string) error {
	// Because we don't provide explicit marshalling for XML in the codebase (as we do for
	// JSON and YAML) then we must normalize the data before encoding it
	data, err := Normalize(data)
	if err != nil {
		return err
	}

	data = EnsureXml(data)

	if _, err := io.WriteString(writer, xml.Header); err != nil {
		return err
	}
	encoder := xml.NewEncoder(writer)
	encoder.Indent("", indent)
	if err := encoder.Encode(data); err != nil {
		return err
	}
	if indent == "" {
		// When there's no indent the XML encoder does not emit a final newline
		// (We want it for consistency with YAML and JSON)
		if _, err := io.WriteString(writer, "\n"); err != nil {
			return err
		}
	}
	return nil
}

func WriteXmlDocument(xmlDocument *etree.Document, writer io.Writer, indent string) error {
	xmlDocument.Indent(len(indent))
	_, err := xmlDocument.WriteTo(writer)
	return err
}
