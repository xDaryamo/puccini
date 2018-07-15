package tosca

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/tliron/puccini/ard"
	"github.com/tliron/puccini/tosca/reflection"
)

type Reader func(*Context) interface{}

const (
	ReadFieldModeDefault       = 0
	ReadFieldModeList          = 1
	ReadFieldModeSequencedList = 2
	ReadFieldModeItem          = 3
)

// From "read" tags
func (self *Context) ReadFields(entityPtr interface{}, readers map[string]Reader) []string {
	if !self.ValidateType("map") {
		return nil
	}

	var keys []string

	data := self.Data.(ard.Map)
	entity := reflect.ValueOf(entityPtr).Elem()
	tags := reflection.GetFieldTagsForValue(entity, "read")

	// Gather all tagged keys
	for _, tag := range tags {
		t := strings.Split(tag.Tag, ",")
		keys = append(keys, t[0])
	}

	for _, tag := range tags {
		key, mode, read := parseReadTag(tag.Tag, readers)

		keys = append(keys, key)

		if key != "?" {
			self.readField(entity, data, key, tag.FieldName, mode, read)
		} else {
			// Iterate all keys that aren't tagged
			for key = range data {
				tagged := false
				for _, k := range keys {
					if key == k {
						tagged = true
						break
					}
				}
				if !tagged {
					self.readField(entity, data, key, tag.FieldName, ReadFieldModeItem, read)
				}
			}
		}
	}

	return keys
}

func (self *Context) readField(entity reflect.Value, data ard.Map, key string, fieldName string, mode int, read Reader) {
	childData, ok := data[key]
	if !ok {
		return
	}

	context := self.FieldChild(key, childData)
	field := entity.FieldByName(fieldName)

	if read != nil {
		fieldType := field.Type()
		if reflection.IsSliceOfPtrToStruct(fieldType) {
			// Field is compatible with []*interface{}
			if read == nil {
				panicMissingReader(entity)
			}
			slice := field
			switch mode {
			case ReadFieldModeList:
				context.ReadListItems(read, func(item interface{}) {
					slice = reflect.Append(slice, reflect.ValueOf(item))
				})
			case ReadFieldModeSequencedList:
				context.ReadSequencedListItems(read, func(item interface{}) {
					slice = reflect.Append(slice, reflect.ValueOf(item))
				})
			case ReadFieldModeItem:
				length := slice.Len()
				slice = reflect.Append(slice, reflect.ValueOf(read(self.ListChild(length, childData))))
			default:
				context.ReadMapItems(read, func(item interface{}) {
					slice = reflect.Append(slice, reflect.ValueOf(item))
				})
			}
			if slice.IsNil() {
				// If we have no items, at least have an empty slice
				// so that "require" will not see a nil here
				slice = reflect.MakeSlice(slice.Type(), 0, 0)
			}
			field.Set(slice)
		} else if reflection.IsMapOfStringToPtrToStruct(fieldType) {
			// Field is compatible with map[string]*interface{}
			if read == nil {
				panicMissingReader(entity)
			}
			switch mode {
			case ReadFieldModeList:
				context.ReadListItems(read, func(item interface{}) {
					context.setMapItem(field, item)
				})
			case ReadFieldModeSequencedList:
				context.ReadSequencedListItems(read, func(item interface{}) {
					context.setMapItem(field, item)
				})
			case ReadFieldModeItem:
				context.setMapItem(field, read(self.MapChild(key, childData)))
			default:
				context.ReadMapItems(read, func(item interface{}) {
					context.setMapItem(field, item)
				})
			}
		} else {
			if read == nil {
				panicMissingReader(entity)
			}
			item := read(context)
			if item != nil {
				field.Set(reflect.ValueOf(item))
			}
		}
	} else {
		fieldEntityPtr := field.Interface()
		if reflection.IsPtrToString(fieldEntityPtr) {
			// Field is *string
			item := context.ReadString()
			if item != nil {
				field.Set(reflect.ValueOf(item))
			}
		} else if reflection.IsPtrToBool(fieldEntityPtr) {
			// Field is *bool
			item := context.ReadBoolean()
			if item != nil {
				field.Set(reflect.ValueOf(item))
			}
		} else if reflection.IsPtrToSliceOfString(fieldEntityPtr) {
			// Field is *[]string
			item := context.ReadStringList()
			if item != nil {
				field.Set(reflect.ValueOf(item))
			}
		} else if reflection.IsPtrToMapOfStringToString(fieldEntityPtr) {
			// Field is *map[string]string
			item := context.ReadStringMap()
			if item != nil {
				field.Set(reflect.ValueOf(item))
			}
		} else {
			panic(fmt.Sprintf("\"read\" tag's field type \"%T\" is not supported in struct: %T.%s", fieldEntityPtr, entity.Interface(), fieldName))
		}
	}
}

func (self *Context) setMapItem(field reflect.Value, item interface{}) {
	key := GetKey(item)
	keyValue := reflect.ValueOf(key)
	itemValue := reflect.ValueOf(item)

	existing := field.MapIndex(keyValue)
	if existing.IsValid() {
		self.ReportMapKeyReused(key)
	}

	field.SetMapIndex(keyValue, itemValue)
}

func parseReadTag(tag string, readers map[string]Reader) (string, int, Reader) {
	t := strings.Split(tag, ",")

	key := t[0]
	mode := ReadFieldModeDefault

	var readerName string
	var read Reader
	if len(t) > 1 {
		readerName = t[1]

		if strings.HasPrefix(readerName, "[]") {
			// List
			readerName = readerName[2:]
			mode = ReadFieldModeList
		} else if strings.HasPrefix(readerName, "{}") {
			// Sequenced list
			readerName = readerName[2:]
			mode = ReadFieldModeSequencedList
		}

		var ok bool
		read, ok = readers[readerName]
		if !ok {
			panic(fmt.Sprintf("reader not found: %s", readerName))
		}
	}

	return key, mode, read
}

func panicMissingReader(entity reflect.Value) {
	panic(fmt.Sprintf("\"read\" tag's field type must be specified in struct %T", entity.Interface()))
}

//
// Read helpers
//

func (self *Context) ReadString() *string {
	if self.ValidateType("string") {
		value := self.Data.(string)
		return &value
	}
	return nil
}

func (self *Context) ReadStringList() *[]string {
	if self.ValidateType("list") {
		var strings []string
		for index, data := range self.Data.(ard.List) {
			string, ok := data.(string)
			if !ok {
				self.ListChild(index, data).ReportValueWrongType("string")
				continue
			}
			strings = append(strings, string)
		}
		return &strings
	}
	return nil
}

func (self *Context) ReadStringListFixed(length int) *[]string {
	strings := self.ReadStringList()
	if (strings != nil) && (len(*strings) != length) {
		self.ReportValueWrongLength("list", length)
		return nil
	}
	return strings
}

func (self *Context) ReadStringMap() *map[string]string {
	if self.ValidateType("map") {
		strings := make(map[string]string)
		for key, data := range self.Data.(ard.Map) {
			string, ok := data.(string)
			if !ok {
				self.MapChild(key, data).ReportValueWrongType("string")
				continue
			}
			strings[key] = string
		}
		return &strings
	}
	return nil
}

func (self *Context) ReadInteger() *int64 {
	if self.ValidateType("integer") {
		var value int64
		switch d := self.Data.(type) {
		case int:
			value = int64(d)
		case int64:
			value = d
		case int32:
			value = int64(d)
		}
		return &value
	}
	return nil
}

func (self *Context) ReadFloat() *float64 {
	if self.ValidateType("float") {
		var value float64
		switch d := self.Data.(type) {
		case float64:
			value = d
		case float32:
			value = float64(d)
		}
		return &value
	}
	return nil
}

func (self *Context) ReadBoolean() *bool {
	if self.ValidateType("boolean") {
		value := self.Data.(bool)
		return &value
	}
	return nil
}

type Processor func(interface{})

func (self *Context) ReadMapItems(read Reader, process Processor) bool {
	if self.ValidateType("map") {
		for itemName, data := range self.Data.(ard.Map) {
			process(read(self.MapChild(itemName, data)))
		}
		return true
	}
	return false
}

func (self *Context) ReadListItems(read Reader, process Processor) bool {
	if self.ValidateType("list") {
		for index, data := range self.Data.(ard.List) {
			process(read(self.ListChild(index, data)))
		}
		return true
	}
	return false
}

func (self *Context) ReadSequencedListItems(read Reader, process Processor) bool {
	if self.ValidateType("list") {
		for index, data := range self.Data.(ard.List) {
			if !reflection.IsMap(data) {
				self.ReportFieldMalformedSequencedList()
				return false
			}
			item := data.(ard.Map)
			if len(item) != 1 {
				self.ReportFieldMalformedSequencedList()
				return false
			}
			for itemName, data := range item {
				process(read(self.SequencedListChild(index, itemName, data)))
			}
		}
		return true
	}
	return false
}
