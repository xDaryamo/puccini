package format

import (
	"encoding/xml"
	"reflect"
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
			v := value.MapIndex(key).Interface()
			slice[index] = XmlMapEntry{key.String(), EnsureXml(v)}
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
	err := encoder.EncodeToken(XmlMapStartElement)
	if err != nil {
		return err
	}
	err = encoder.Encode(self.Entries)
	if err != nil {
		return err
	}
	return encoder.EncodeToken(XmlMapStartElement.End())
}

//
// XmlMapEntry
//

type XmlMapEntry struct {
	Key   string
	Value interface{}
}

var XmlMapEntryStart = xml.StartElement{Name: xml.Name{Local: "entry"}}
var XmlKeyStart = xml.StartElement{Name: xml.Name{Local: "key"}}
var XmlValueStart = xml.StartElement{Name: xml.Name{Local: "value"}}

// xml.Marshaler interface
func (self XmlMapEntry) MarshalXML(encoder *xml.Encoder, start xml.StartElement) error {
	err := encoder.EncodeToken(XmlMapEntryStart)
	if err != nil {
		return err
	}
	err = encoder.EncodeElement(self.Key, XmlKeyStart)
	if err != nil {
		return err
	}
	err = encoder.EncodeToken(XmlValueStart)
	if err != nil {
		return err
	}
	err = encoder.Encode(self.Value)
	if err != nil {
		return err
	}
	err = encoder.EncodeToken(XmlValueStart.End())
	if err != nil {
		return err
	}
	return encoder.EncodeToken(XmlMapEntryStart.End())
}
