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
		return WriteYaml(data, writer)
	case "json":
		return WriteJson(data, writer, indent)
	case "xml":
		return WriteXml(data, writer, indent)
	default:
		return fmt.Errorf("unsupported format: %s", format)
	}
}

func WriteYaml(data interface{}, writer io.Writer) error {
	encoder := yaml.NewEncoder(writer)
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
	return encoder.Encode(data)
}

func WriteXmlDocument(xmlDocument *etree.Document, writer io.Writer, indent string) error {
	xmlDocument.Indent(len(indent))
	_, err := xmlDocument.WriteTo(writer)
	return err
}
