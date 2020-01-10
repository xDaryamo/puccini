package format

import (
	"encoding/xml"
	"reflect"

	"github.com/tliron/yamlkeys"
)

func EnsureXML(data interface{}) interface{} {
	if data == nil {
		return nil
	}

	value := reflect.ValueOf(data)
	type_ := value.Type()

	if type_.Kind() == reflect.Slice {
		length := value.Len()
		slice := make([]interface{}, length)
		for index := 0; index < length; index++ {
			v := value.Index(index).Interface()
			slice[index] = EnsureXML(v)
		}
		return slice
	} else if type_.Kind() == reflect.Map {
		// Convert to slice of XMLMapEntry
		slice := make([]XMLMapEntry, value.Len())
		for index, key := range value.MapKeys() {
			k := yamlkeys.KeyData(key.Interface())
			v := value.MapIndex(key).Interface()
			slice[index] = XMLMapEntry{EnsureXML(k), EnsureXML(v)}
		}
		return XMLMap{slice}
	}

	return data
}

//
// XMLMap
//

type XMLMap struct {
	Entries []XMLMapEntry
}

var XMLMapStartElement = xml.StartElement{Name: xml.Name{Local: "map"}}

// xml.Marshaler interface
func (self XMLMap) MarshalXML(encoder *xml.Encoder, start xml.StartElement) error {
	var err error
	if err = encoder.EncodeToken(XMLMapStartElement); err != nil {
		return err
	}
	if err = encoder.Encode(self.Entries); err != nil {
		return err
	}
	return encoder.EncodeToken(XMLMapStartElement.End())
}

//
// XMLMapEntry
//

type XMLMapEntry struct {
	Key   interface{}
	Value interface{}
}

var XMLMapEntryStart = xml.StartElement{Name: xml.Name{Local: "entry"}}
var XMLKeyStart = xml.StartElement{Name: xml.Name{Local: "key"}}
var XMLValueStart = xml.StartElement{Name: xml.Name{Local: "value"}}

// xml.Marshaler interface
func (self XMLMapEntry) MarshalXML(encoder *xml.Encoder, start xml.StartElement) error {
	var err error
	if err := encoder.EncodeToken(XMLMapEntryStart); err != nil {
		return err
	}
	if err := encoder.EncodeElement(self.Key, XMLKeyStart); err != nil {
		return err
	}
	if err := encoder.EncodeToken(XMLValueStart); err != nil {
		return err
	}
	if err = encoder.Encode(self.Value); err != nil {
		return err
	}
	if err = encoder.EncodeToken(XMLValueStart.End()); err != nil {
		return err
	}
	return encoder.EncodeToken(XMLMapEntryStart.End())
}
