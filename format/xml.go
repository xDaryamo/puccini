package format

import (
	"encoding/xml"
	"reflect"

	"github.com/tliron/yamlkeys"
)

func EnsureXml(data interface{}) interface{} {
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
			slice[index] = EnsureXml(v)
		}
		return slice
	} else if type_.Kind() == reflect.Map {
		// Convert to slice of XmlMapEntry
		slice := make([]XmlMapEntry, value.Len())
		for index, key := range value.MapKeys() {
			k := yamlkeys.KeyData(key.Interface())
			v := value.MapIndex(key).Interface()
			slice[index] = XmlMapEntry{EnsureXml(k), EnsureXml(v)}
		}
		return XmlMap{slice}
	}

	return data
}

//
// XmlMap
//

type XmlMap struct {
	Entries []XmlMapEntry
}

var XmlMapStartElement = xml.StartElement{Name: xml.Name{Local: "map"}}

// xml.Marshaler interface
func (self XmlMap) MarshalXML(encoder *xml.Encoder, start xml.StartElement) error {
	var err error
	if err = encoder.EncodeToken(XmlMapStartElement); err != nil {
		return err
	}
	if err = encoder.Encode(self.Entries); err != nil {
		return err
	}
	return encoder.EncodeToken(XmlMapStartElement.End())
}

//
// XmlMapEntry
//

type XmlMapEntry struct {
	Key   interface{}
	Value interface{}
}

var XmlMapEntryStart = xml.StartElement{Name: xml.Name{Local: "entry"}}
var XmlKeyStart = xml.StartElement{Name: xml.Name{Local: "key"}}
var XmlValueStart = xml.StartElement{Name: xml.Name{Local: "value"}}

// xml.Marshaler interface
func (self XmlMapEntry) MarshalXML(encoder *xml.Encoder, start xml.StartElement) error {
	var err error
	if err := encoder.EncodeToken(XmlMapEntryStart); err != nil {
		return err
	}
	if err := encoder.EncodeElement(self.Key, XmlKeyStart); err != nil {
		return err
	}
	if err := encoder.EncodeToken(XmlValueStart); err != nil {
		return err
	}
	if err = encoder.Encode(self.Value); err != nil {
		return err
	}
	if err = encoder.EncodeToken(XmlValueStart.End()); err != nil {
		return err
	}
	return encoder.EncodeToken(XmlMapEntryStart.End())
}
