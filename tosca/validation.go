package tosca

import (
	"reflect"

	"github.com/tliron/puccini/ard"
	"github.com/tliron/puccini/tosca/reflection"
)

func (self *Context) ValidateUnsupportedFields(keys []string) {
	if !self.Is("!!map") {
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

func (self *Context) ValidateType(requiredTypeNames ...string) bool {
	is := self.Is(requiredTypeNames...)
	if !is {
		self.ReportValueWrongType(requiredTypeNames...)
	}
	return is
}

// From "require" tags
func ValidateRequiredFields(entityPtr interface{}) bool {
	context := GetContext(entityPtr)
	entity := reflect.ValueOf(entityPtr).Elem()
	for fieldName, tag := range reflection.GetFieldTagsForValue(entity, "require") {
		field := entity.FieldByName(fieldName)
		if reflection.IsNil(field) {
			context.FieldChild(tag, nil).ReportFieldMissing()
		}
	}
	return true
}
