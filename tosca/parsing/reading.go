package parsing

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/tliron/go-ard"
	"github.com/tliron/kutil/reflection"
	"github.com/tliron/yamlkeys"
)

type Reader = func(*Context) EntityPtr

type Readers map[string]Reader

//
// PreReadable
//

type PreReadable interface {
	PreRead()
}

// From [PreReadable] interface
func PreRead(entityPtr EntityPtr) bool {
	if preReadable, ok := entityPtr.(PreReadable); ok {
		preReadable.PreRead()
		return true
	} else {
		return false
	}
}

// From "read" tags
func (self *Context) ReadFields(entityPtr EntityPtr) []string {
	PreRead(entityPtr)

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

func (self *Context) setMapItem(field reflect.Value, item any) {
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

func (self *Context) appendUnique(field reflect.Value, item any) reflect.Value {
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

type ReadMode int

const (
	ReadFieldModeDefault             ReadMode = 0
	ReadFieldModeList                ReadMode = 1
	ReadFieldModeSequencedList       ReadMode = 2
	ReadFieldModeUniqueSequencedList ReadMode = 3
	ReadFieldModeItem                ReadMode = 4
)

type ReadField struct {
	FieldName string
	Key       string
	Context   *Context
	Entity    reflect.Value
	Reader    Reader
	Mode      ReadMode
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
	field := self.Entity.FieldByName(self.FieldName)
	fieldType := field.Type()

	//----------------------------------------------------------------------
	// 0) Is the key present at all?
	//----------------------------------------------------------------------
	childData, ok := self.Context.Data.(ard.Map)[self.Key]
	if !ok {
		if reflection.IsSliceOfPointerToStruct(fieldType) {
			// Ensure “mandatory” logic sees an empty (not nil) slice
			field.Set(reflect.MakeSlice(fieldType, 0, 0))
		}
		return
	}

	//----------------------------------------------------------------------
	// 1) Paths that use a custom reader
	//----------------------------------------------------------------------
	if self.Reader != nil {

		//--------------------------------------------------------------
		// 1a) Destination is []*Struct
		//--------------------------------------------------------------
		if reflection.IsSliceOfPointerToStruct(fieldType) {
			slice := field
			switch self.Mode {

			case ReadFieldModeList:
				self.Context.FieldChild(self.Key, childData).
					ReadListItems(self.Reader, func(item any) {
						if item != nil { // <-- guard
							slice = reflect.Append(slice, reflect.ValueOf(item))
						}
					})

			case ReadFieldModeSequencedList:
				self.Context.FieldChild(self.Key, childData).
					ReadSequencedListItems(self.Reader, func(item any) {
						if item != nil { // <-- guard
							slice = reflect.Append(slice, reflect.ValueOf(item))
						}
					})

			case ReadFieldModeUniqueSequencedList:
				context := self.Context.FieldChild(self.Key, childData)
				context.ReadSequencedListItems(self.Reader, func(item any) {
					if item != nil { // <-- guard
						slice = context.appendUnique(slice, item)
					}
				})

			case ReadFieldModeItem:
				length := slice.Len()
				if itm := self.Reader(self.Context.ListChild(length, childData)); itm != nil {
					slice = reflect.Append(slice, reflect.ValueOf(itm))
				}

			default: // map-style syntax
				self.Context.FieldChild(self.Key, childData).
					ReadMapItems(self.Reader, func(item any) {
						if item != nil { // <-- guard
							slice = reflect.Append(slice, reflect.ValueOf(item))
						}
					})
			}

			if slice.IsNil() {
				slice = reflect.MakeSlice(fieldType, 0, 0)
			}
			field.Set(slice)

			//--------------------------------------------------------------
			// 1b) Destination is map[string]*Struct
			//--------------------------------------------------------------
		} else if reflection.IsMapOfStringToPointerToStruct(fieldType) {
			switch self.Mode {

			case ReadFieldModeList:
				ctx := self.Context.FieldChild(self.Key, childData)
				ctx.ReadListItems(self.Reader, func(item any) {
					if item != nil { // <-- guard
						ctx.setMapItem(field, item)
					}
				})

			case ReadFieldModeSequencedList, ReadFieldModeUniqueSequencedList:
				ctx := self.Context.FieldChild(self.Key, childData)
				ctx.ReadSequencedListItems(self.Reader, func(item any) {
					if item != nil { // <-- guard
						ctx.setMapItem(field, item)
					}
				})

			case ReadFieldModeItem:
				ctx := self.Context.FieldChild(self.Key, childData)
				if itm := self.Reader(ctx); itm != nil {
					ctx.setMapItem(field, itm)
				}

			default:
				ctx := self.Context.FieldChild(self.Key, childData)
				ctx.ReadMapItems(self.Reader, func(item any) {
					if item != nil { // <-- guard
						ctx.setMapItem(field, item)
					}
				})
			}

			//--------------------------------------------------------------
			// 1c) Destination is *Struct
			//--------------------------------------------------------------
		} else {
			if item := self.Reader(self.Context.FieldChild(self.Key, childData)); item != nil {
				field.Set(reflect.ValueOf(item))
			}
		}

		//----------------------------------------------------------------------
		// 2) Paths that rely on the built-in scalar readers
		//----------------------------------------------------------------------
	} else {
		fieldEntityPtr := field.Interface()

		switch {
		case reflection.IsPointerToString(fieldEntityPtr):
			if item := self.Context.FieldChild(self.Key, childData).ReadString(); item != nil {
				field.Set(reflect.ValueOf(item))
			}
		case reflection.IsPointerToInt64(fieldEntityPtr):
			if item := self.Context.FieldChild(self.Key, childData).ReadInteger(); item != nil {
				field.Set(reflect.ValueOf(item))
			}
		case reflection.IsPointerToFloat64(fieldEntityPtr):
			if item := self.Context.FieldChild(self.Key, childData).ReadFloat(); item != nil {
				field.Set(reflect.ValueOf(item))
			}
		case reflection.IsPointerToBool(fieldEntityPtr):
			if item := self.Context.FieldChild(self.Key, childData).ReadBoolean(); item != nil {
				field.Set(reflect.ValueOf(item))
			}
		case reflection.IsPointerToSliceOfString(fieldEntityPtr):
			if item := self.Context.FieldChild(self.Key, childData).ReadStringList(); item != nil {
				field.Set(reflect.ValueOf(item))
			}
		case reflection.IsPointerToMapOfStringToString(fieldEntityPtr):
			if item := self.Context.FieldChild(self.Key, childData).ReadStringStringMap(); item != nil {
				field.Set(reflect.ValueOf(item))
			}
		default:
			panic(fmt.Sprintf("\"read\" tag's field type \"%T\" is not supported in struct: %T.%s",
				fieldEntityPtr, self.Entity.Interface(), self.FieldName))
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

func (self *Context) ReadStringListMinLength(length int) *[]string {
	strings := self.ReadStringList()
	if (strings != nil) && (len(*strings) < length) {
		self.ReportValueWrongLength("list", length)
		return nil
	}
	return strings
}

func (self *Context) ReadStringMap() *map[string]any {
	if self.ValidateType(ard.TypeMap) {
		map_ := make(map[string]any)
		for key, data := range self.Data.(ard.Map) {
			var ok bool
			var key_ string
			if key_, ok = key.(string); ok {
			} else if self.HasQuirk(QuirkDataTypesStringPermissive) {
				key_ = ard.ValueToString(data)
			} else {
				self.MapChild(key, data).ReportValueAspectWrongType("map key", key, ard.TypeString)
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
				self.MapChild(key, data).ReportValueAspectWrongType("map key", key, ard.TypeString)
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
		case uint64:
			value = int64(d)
		case uint32:
			value = int64(d)
		case uint16:
			value = int64(d)
		case uint8:
			value = int64(d)
		case uint:
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
				self.ReportKeynameMalformedSequencedList()
				return false
			}
			item := data.(ard.Map)
			if len(item) != 1 {
				self.ReportKeynameMalformedSequencedList()
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

func (self *Context) getReadTagKey(entity reflect.Value, fieldName string) (string, bool) {
	var tag string
	var ok bool

	if self.ReadTagOverrides != nil {
		tag, ok = self.ReadTagOverrides[fieldName]
	}

	if !ok {
		if structField, ok_ := entity.Type().FieldByName(fieldName); ok_ {
			tag, ok = structField.Tag.Lookup("read")
		}
	}

	if ok {
		t := strings.Split(tag, ",")
		return t[0], true
	} else {
		return "", false
	}
}
