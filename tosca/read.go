package tosca

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/tliron/kutil/ard"
	"github.com/tliron/kutil/reflection"
	"github.com/tliron/yamlkeys"
)

type Reader = func(*Context) EntityPtr

type Readers map[string]Reader

type PreReadable interface {
	PreRead()
}

const (
	ReadFieldModeDefault             = 0
	ReadFieldModeList                = 1
	ReadFieldModeSequencedList       = 2
	ReadFieldModeUniqueSequencedList = 3
	ReadFieldModeItem                = 4
)

// From "read" tags
func (self *Context) ReadFields(entityPtr EntityPtr) []string {
	if preReadable, ok := entityPtr.(PreReadable); ok {
		preReadable.PreRead()
	}

	if !self.ValidateType(ard.TypeMap) {
		return nil
	}

	var keys []string

	entity := reflect.ValueOf(entityPtr).Elem()
	tags := reflection.GetFieldTagsForValue(entity, "read")

	// Read tag overrides
	if self.ReadTagOverrides != nil {
		for fieldName, tag := range self.ReadTagOverrides {
			if tag != "" {
				tags[fieldName] = tag
			} else {
				// Empty tag means delete
				if _, ok := tags[fieldName]; !ok {
					panic(fmt.Sprintf("unknown read field: %q", fieldName))
				}
				delete(tags, fieldName)
			}
		}
	}

	// Gather all tagged keys
	for _, tag := range tags {
		t := strings.Split(tag, ",")
		keys = append(keys, t[0])
	}

	// Parse tags
	var readFields []*ReadField
	for fieldName, tag := range tags {
		readField := NewReadField(fieldName, tag, self, entity)
		if readField.Important {
			// Important fields come first
			readFields = append([]*ReadField{readField}, readFields...)
		} else {
			readFields = append(readFields, readField)
		}
	}

	for _, readField := range readFields {
		if readField.Wildcard {
			// Iterate all keys that aren't tagged
			for key := range self.Data.(ard.Map) {
				tagged := false
				for _, k := range keys {
					if key == k {
						tagged = true
						break
					}
				}
				if !tagged {
					readField.Key = yamlkeys.KeyString(key)
					readField.Read()
				}
			}
		} else {
			readField.Read()
		}
	}

	return keys
}

func (self *Context) setMapItem(field reflect.Value, item interface{}) {
	key := GetKey(item)
	keyValue := reflect.ValueOf(key)
	itemValue := reflect.ValueOf(item)

	existing := field.MapIndex(keyValue)
	if existing.IsValid() {
		self.ReportDuplicateMapKey(key)
		return
	}

	field.SetMapIndex(keyValue, itemValue)
}

func (self *Context) appendUnique(field reflect.Value, item interface{}) reflect.Value {
	length := field.Len()
	if length > 0 {
		key := GetKey(item)
		for index := 0; index < length; index++ {
			element := field.Index(index).Interface()
			if key == GetKey(element) {
				self.ReportDuplicateMapKey(key)
				return field
			}
		}
	}

	return reflect.Append(field, reflect.ValueOf(item))
}

//
// ReadField
//

type ReadField struct {
	FieldName string
	Key       string
	Context   *Context
	Entity    reflect.Value
	Reader    Reader
	Mode      int
	Important bool
	Wildcard  bool
}

func NewReadField(fieldName string, tag string, context *Context, entity reflect.Value) *ReadField {
	// TODO: is it worth caching some of this?

	t := strings.Split(tag, ",")

	var self = ReadField{
		FieldName: fieldName,
		Key:       t[0],
		Context:   context,
		Entity:    entity,
	}

	if self.Key == "?" {
		self.Wildcard = true
		self.Mode = ReadFieldModeItem
	} else {
		self.Mode = ReadFieldModeDefault
	}

	var readerName string
	if len(t) > 1 {
		readerName = t[1]

		if strings.HasPrefix(readerName, "!") {
			// Important
			readerName = readerName[1:]
			self.Important = true
		}

		if strings.HasPrefix(readerName, "[]") {
			// List
			readerName = readerName[2:]
			self.Mode = ReadFieldModeList
		} else if strings.HasPrefix(readerName, "{}") {
			// Sequenced list
			readerName = readerName[2:]
			self.Mode = ReadFieldModeSequencedList
		} else if strings.HasPrefix(readerName, "<>") {
			// Unique sequenced list
			readerName = readerName[2:]
			self.Mode = ReadFieldModeUniqueSequencedList
		}

		var ok bool
		self.Reader, ok = context.Grammar.Readers[readerName]
		if !ok {
			panic(fmt.Sprintf("reader not found: %s", readerName))
		}
	}

	return &self
}

func (self *ReadField) Read() {
	childData, ok := self.Context.Data.(ard.Map)[self.Key]
	if !ok {
		return
	}

	field := self.Entity.FieldByName(self.FieldName)

	if self.Reader != nil {
		fieldType := field.Type()
		if reflection.IsSliceOfPtrToStruct(fieldType) {
			// Field is compatible with []*interface{}
			slice := field
			switch self.Mode {
			case ReadFieldModeList:
				self.Context.FieldChild(self.Key, childData).ReadListItems(self.Reader, func(item interface{}) {
					slice = reflect.Append(slice, reflect.ValueOf(item))
				})
			case ReadFieldModeSequencedList:
				self.Context.FieldChild(self.Key, childData).ReadSequencedListItems(self.Reader, func(item interface{}) {
					slice = reflect.Append(slice, reflect.ValueOf(item))
				})
			case ReadFieldModeUniqueSequencedList:
				context := self.Context.FieldChild(self.Key, childData)
				context.ReadSequencedListItems(self.Reader, func(item interface{}) {
					slice = context.appendUnique(slice, item)
				})
			case ReadFieldModeItem:
				length := slice.Len()
				slice = reflect.Append(slice, reflect.ValueOf(self.Reader(self.Context.ListChild(length, childData))))
			default:
				self.Context.FieldChild(self.Key, childData).ReadMapItems(self.Reader, func(item interface{}) {
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
			switch self.Mode {
			case ReadFieldModeList:
				context := self.Context.FieldChild(self.Key, childData)
				context.ReadListItems(self.Reader, func(item interface{}) {
					context.setMapItem(field, item)
				})
			case ReadFieldModeSequencedList, ReadFieldModeUniqueSequencedList:
				context := self.Context.FieldChild(self.Key, childData)
				context.ReadSequencedListItems(self.Reader, func(item interface{}) {
					context.setMapItem(field, item)
				})
			case ReadFieldModeItem:
				context := self.Context.FieldChild(self.Key, childData)
				context.setMapItem(field, self.Reader(context))
			default:
				context := self.Context.FieldChild(self.Key, childData)
				context.ReadMapItems(self.Reader, func(item interface{}) {
					context.setMapItem(field, item)
				})
			}
		} else {
			// Field is compatible with *interface{}
			item := self.Reader(self.Context.FieldChild(self.Key, childData))
			if item != nil {
				field.Set(reflect.ValueOf(item))
			}
		}
	} else {
		fieldEntityPtr := field.Interface()
		if reflection.IsPtrToString(fieldEntityPtr) {
			// Field is *string
			item := self.Context.FieldChild(self.Key, childData).ReadString()
			if item != nil {
				field.Set(reflect.ValueOf(item))
			}
		} else if reflection.IsPtrToInt64(fieldEntityPtr) {
			// Field is *int64
			item := self.Context.FieldChild(self.Key, childData).ReadInteger()
			if item != nil {
				field.Set(reflect.ValueOf(item))
			}
		} else if reflection.IsPtrToFloat64(fieldEntityPtr) {
			// Field is *float64
			item := self.Context.FieldChild(self.Key, childData).ReadFloat()
			if item != nil {
				field.Set(reflect.ValueOf(item))
			}
		} else if reflection.IsPtrToBool(fieldEntityPtr) {
			// Field is *bool
			item := self.Context.FieldChild(self.Key, childData).ReadBoolean()
			if item != nil {
				field.Set(reflect.ValueOf(item))
			}
		} else if reflection.IsPtrToSliceOfString(fieldEntityPtr) {
			// Field is *[]string
			item := self.Context.FieldChild(self.Key, childData).ReadStringList()
			if item != nil {
				field.Set(reflect.ValueOf(item))
			}
		} else if reflection.IsPtrToMapOfStringToString(fieldEntityPtr) {
			// Field is *map[string]string
			item := self.Context.FieldChild(self.Key, childData).ReadStringStringMap()
			if item != nil {
				field.Set(reflect.ValueOf(item))
			}
		} else {
			panic(fmt.Sprintf("\"read\" tag's field type \"%T\" is not supported in struct: %T.%s", fieldEntityPtr, self.Entity.Interface(), self.FieldName))
		}
	}
}

//
// Read helpers
//

func (self *Context) ReadString() *string {
	if self.ValidateType(ard.TypeString) {
		value := self.Data.(string)
		return &value
	}
	return nil
}

func (self *Context) ReadStringList() *[]string {
	if self.ValidateType(ard.TypeList) {
		var strings []string
		for index, data := range self.Data.(ard.List) {
			string, ok := data.(string)
			if !ok {
				self.ListChild(index, data).ReportValueWrongType(ard.TypeString)
				continue
			}
			strings = append(strings, string)
		}
		return &strings
	}
	return nil
}

func (self *Context) ReadStringOrStringList() *[]string {
	if self.Is(ard.TypeList) {
		return self.ReadStringList()
	} else if self.ValidateType(ard.TypeList, ard.TypeString) {
		return &[]string{*self.ReadString()}
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

func (self *Context) ReadStringMap() *map[string]interface{} {
	if self.ValidateType(ard.TypeMap) {
		map_ := make(map[string]interface{})
		for key, data := range self.Data.(ard.Map) {
			var ok bool
			var key_ string
			if key_, ok = key.(string); ok {
			} else if self.HasQuirk(QuirkDataTypesStringPermissive) {
				key_ = ard.ValueToString(data)
			} else {
				self.MapChild(key, data).ReportAspectWrongType("map key", key, ard.TypeString)
				continue
			}

			map_[key_] = data
		}
		return &map_
	} else {
		return nil
	}
}

func (self *Context) ReadStringStringMap() *map[string]string {
	if self.ValidateType(ard.TypeMap) {
		map_ := make(map[string]string)
		for key, data := range self.Data.(ard.Map) {
			var ok bool
			var key_ string
			if key_, ok = key.(string); ok {
			} else if self.HasQuirk(QuirkDataTypesStringPermissive) {
				key_ = ard.ValueToString(data)
			} else {
				self.MapChild(key, data).ReportAspectWrongType("map key", key, ard.TypeString)
				continue
			}

			var data_ string
			if data_, ok = data.(string); ok {
			} else if self.HasQuirk(QuirkDataTypesStringPermissive) {
				data_ = ard.ValueToString(data)
			} else {
				self.MapChild(key, data).ReportValueWrongType(ard.TypeString)
				continue
			}

			map_[key_] = data_
		}
		return &map_
	} else {
		return nil
	}
}

func (self *Context) ReadInteger() *int64 {
	if self.ValidateType(ard.TypeInteger) {
		var value int64
		switch d := self.Data.(type) {
		case int64:
			value = d
		case int32:
			value = int64(d)
		case int16:
			value = int64(d)
		case int8:
			value = int64(d)
		case int:
			value = int64(d)
		}
		return &value
	}
	return nil
}

func (self *Context) ReadFloat() *float64 {
	if self.ValidateType(ard.TypeFloat) {
		var value float64
		switch data := self.Data.(type) {
		case float64:
			value = data
		case float32:
			value = float64(data)
		}
		return &value
	}
	return nil
}

func (self *Context) ReadBoolean() *bool {
	if self.ValidateType(ard.TypeBoolean) {
		value := self.Data.(bool)
		return &value
	}
	return nil
}

type Processor = func(ard.Value)

func (self *Context) ReadMapItems(read Reader, process Processor) bool {
	if self.ValidateType(ard.TypeMap) {
		for itemName, data := range self.Data.(ard.Map) {
			process(read(self.MapChild(itemName, data)))
		}
		return true
	}
	return false
}

func (self *Context) ReadListItems(read Reader, process Processor) bool {
	if self.ValidateType(ard.TypeList) {
		for index, data := range self.Data.(ard.List) {
			process(read(self.ListChild(index, data)))
		}
		return true
	}
	return false
}

func (self *Context) ReadSequencedListItems(read Reader, process Processor) bool {
	if self.ValidateType(ard.TypeList) {
		for index, data := range self.Data.(ard.List) {
			if !ard.IsMap(data) {
				self.ReportFieldMalformedSequencedList()
				return false
			}
			item := data.(ard.Map)
			if len(item) != 1 {
				self.ReportFieldMalformedSequencedList()
				return false
			}
			for itemName, data := range item {
				process(read(self.SequencedListChild(index, yamlkeys.KeyString(itemName), data)))
			}
		}
		return true
	}
	return false
}

// Utils

func getReadTagKey(entity reflect.Value, fieldName string) (string, bool) {
	if structField, ok := entity.Type().FieldByName(fieldName); ok {
		if tag, ok := structField.Tag.Lookup("read"); ok {
			t := strings.Split(tag, ",")
			return t[0], true
		}
	}
	return "", false
}
