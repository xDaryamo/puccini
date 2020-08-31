package tosca

import (
	"reflect"

	"github.com/tliron/kutil/ard"
	"github.com/tliron/kutil/reflection"
)

func (self *Context) ValidateUnsupportedFields(keys []string) {
	if !self.Is(ard.TypeMap) {
		return
	}
	for key := range self.Data.(ard.Map) {
		found := false
		for _, key_ := range keys {
			if key == key_ {
				found = true
				break
			}
		}
		if !found {
			self.FieldChild(key, nil).ReportFieldUnsupported()
		}
	}
}

func (self *Context) ValidateType(requiredTypeNames ...ard.TypeName) bool {
	is := self.Is(requiredTypeNames...)
	if !is {
		self.ReportValueWrongType(requiredTypeNames...)
	}
	return is
}

// From "require" tags
func ValidateRequiredFields(entityPtr EntityPtr) bool {
	context := GetContext(entityPtr)
	entity := reflect.ValueOf(entityPtr).Elem()
	for fieldName, tag := range reflection.GetFieldTagsForValue(entity, "require") {
		field := entity.FieldByName(fieldName)
		if reflection.IsNil(field) {
			/// Try to use the "read" tag for the problem report
			if readTag, ok := getReadTagKey(entity, fieldName); ok {
				tag = readTag
			}

			context.FieldChild(tag, nil).ReportFieldMissing()
		}
	}
	return true
}
