package parser

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/tliron/kutil/reflection"
	"github.com/tliron/puccini/tosca"
)

func (self *Context) LookupNames() {
	self.Traverse(logLookup, self.LookupFields)
}

// From "lookup" tags
func (self *Context) LookupFields(entityPtr tosca.EntityPtr) bool {
	lookupProblems := make(LookupProblems)

	context := tosca.GetContext(entityPtr)
	entity := reflect.ValueOf(entityPtr).Elem()

	for fieldName, tag := range reflection.GetFieldTagsForValue(entity, "lookup") {
		lookupFieldKey, lookupFieldName, report := parseLookupTag(tag)

		// Field to fill in
		targetField := entity.FieldByName(fieldName)

		// Field with lookup name
		lookupField := entity.FieldByName(lookupFieldName)
		if !lookupField.IsValid() {
			panicLookupTag(entityPtr)
		}
		lookupFieldType := lookupField.Type()
		lookupFieldPtrType := lookupFieldType.Elem()
		lookupFieldPtrTypeKind := lookupFieldPtrType.Kind()
		if lookupFieldType.Kind() != reflect.Ptr {
			panicLookupTag(entityPtr)
		}

		if lookupFieldPtrTypeKind == reflect.String {
			// Lookup field is *string

			// TODO: panic if targetField is not of the right type

			lookupType := targetField.Type()
			lookupProblems.AddType(lookupFieldKey, lookupType)

			// Name to lookup
			lookupName := lookupField.Interface().(*string)
			if lookupName == nil {
				continue
			}

			targetPtr, ok := context.Namespace.LookupForType(*lookupName, lookupType)
			if ok {
				targetField.Set(reflect.ValueOf(targetPtr))
			}
			if report {
				lookupProblems.SetFound(lookupFieldKey, -1, *lookupName, ok)
			}
		} else if (lookupFieldPtrTypeKind == reflect.Slice) && (lookupFieldPtrType.Elem().Kind() == reflect.String) {
			// Lookup field is *[]string

			// TODO: panic if targetField is not of the right type

			lookupType := targetField.Type().Elem()
			lookupProblems.AddType(lookupFieldKey, lookupType)

			// Names to lookup
			lookupNames := lookupField.Interface().(*[]string)
			if (lookupNames == nil) || (len(*lookupNames) == 0) {
				continue
			}

			targetPtrs := targetField
			for index, lookupName := range *lookupNames {
				targetPtr, ok := context.Namespace.LookupForType(lookupName, lookupType)
				if ok {
					targetPtrs = reflect.Append(targetPtrs, reflect.ValueOf(targetPtr))
				}
				if report {
					lookupProblems.SetFound(lookupFieldKey, index, lookupName, ok)
				}
			}
			targetField.Set(targetPtrs)
		} else {
			panicLookupTag(entityPtr)
		}
	}

	lookupProblems.Report(context)

	return true
}

func parseLookupTag(tag string) (string, string, bool) {
	t := strings.Split(tag, ",")

	if len(t) != 2 {
		panic("must be 2")
	}

	lookupFieldKey := t[0]
	lookupFieldName := t[1]

	report := true
	if strings.HasPrefix(lookupFieldName, "?") {
		lookupFieldName = lookupFieldName[1:]
		report = false
	}

	return lookupFieldKey, lookupFieldName, report
}

func panicLookupTag(entityPtr tosca.EntityPtr) {
	panic(fmt.Sprintf("\"lookup\" tag refers to a field that is not of type \"*string\" or \"*[]string\" in struct %T", entityPtr))
}

//
// LookupField
//

type LookupField struct {
	Types []reflect.Type
	Names []LookupName
}

type LookupName struct {
	Index int
	Name  string
	Found bool
}

func (self *LookupField) addType(type_ reflect.Type) {
	type_ = type_.Elem()
	for _, t := range self.Types {
		if t == type_ {
			return
		}
	}
	self.Types = append(self.Types, type_)
}

func (self *LookupField) setFound(index int, name string, found bool) {
	for i, n := range self.Names {
		if (n.Index == index) && (n.Name == name) {
			// Don't change an existing true to false
			if found {
				self.Names[i].Found = true
			}
			return
		}
	}
	self.Names = append(self.Names, LookupName{index, name, found})
}

//
// LookupProblems
//

type LookupProblems map[string]*LookupField

func (self LookupProblems) Field(key string) *LookupField {
	field, ok := self[key]
	if !ok {
		field = new(LookupField)
		self[key] = field
	}
	return field
}

func (self LookupProblems) AddType(key string, type_ reflect.Type) {
	self.Field(key).addType(type_)
}

func (self LookupProblems) SetFound(key string, index int, name string, found bool) {
	self.Field(key).setFound(index, name, found)
}

func (self LookupProblems) Report(context *tosca.Context) {
	for key, field := range self {
		for _, name := range field.Names {
			if !name.Found {
				if name.Index == -1 {
					context.FieldChild(key, name.Name).ReportFieldReferenceNotFound(field.Types...)
				} else {
					context.FieldChild(key, nil).ListChild(name.Index, name.Name).ReportFieldReferenceNotFound(field.Types...)
				}
			}
		}
	}
}
