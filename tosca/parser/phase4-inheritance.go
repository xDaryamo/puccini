package parser

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/tliron/kutil/reflection"
	"github.com/tliron/puccini/tosca/parsing"
)

func (self *Context) Inherit(debug func(Tasks)) {
	self.Parser.lock.Lock()
	defer self.Parser.lock.Unlock()

	tasks := self.getInheritTasks()
	if debug != nil {
		debug(tasks)
	}
	tasks.Drain()
}

func (self *Context) getInheritTasks() Tasks {
	inheritContext := NewInheritContext()
	self.TraverseEntities(logInheritance, self.Parser.getInheritTaskWork, func(entityPtr parsing.EntityPtr) bool {
		inheritContext.GetInheritTask(entityPtr)
		return true
	})
	return inheritContext.Tasks
}

//
// InheritContext
//

type InheritContext struct {
	Tasks            Tasks
	TasksForEntities TasksForEntities
	InheritFields    InheritFields
}

func NewInheritContext() *InheritContext {
	return &InheritContext{make(Tasks), make(TasksForEntities), make(InheritFields)}
}

func (self *InheritContext) GetInheritTask(entityPtr parsing.EntityPtr) *Task {
	task, ok := self.TasksForEntities[entityPtr]
	if !ok {
		path := parsing.GetContext(entityPtr).Path.String()
		if path == "" {
			path = "<root>"
		}

		task = NewTask(path)
		self.Tasks.Add(task)
		self.TasksForEntities[entityPtr] = task

		for dependencyEntityPtr := range self.GetDependencies(entityPtr) {
			task.AddDependency(self.GetInheritTask(dependencyEntityPtr))
		}

		task.Executor = self.NewExecutor(entityPtr)

		logInheritance.Debugf("new task: %s (%d)", task.Name, len(task.Dependencies))
	} else {
		logInheritance.Debugf("task cache hit: %s (%d)", task.Name, len(task.Dependencies))
	}
	return task
}

func (self *InheritContext) NewExecutor(entityPtr parsing.EntityPtr) Executor {
	return func(task *Task) {
		defer task.Done()

		logInheritance.Debugf("task: %s", task.Name)

		for _, inheritField := range self.InheritFields.Get(entityPtr) {
			inheritField.Inherit()
		}

		// Custom inheritance after all fields have been inherited
		parsing.Inherit(entityPtr)
	}
}

func (self *InheritContext) GetDependencies(entityPtr parsing.EntityPtr) parsing.EntityPtrSet {
	dependencies := make(parsing.EntityPtrSet)

	// From "inherit" tags
	for _, inheritField := range self.InheritFields.Get(entityPtr) {
		dependencies.Add(inheritField.ParentEntityPtr)
	}

	// From field values
	entity := reflect.ValueOf(entityPtr).Elem()
	for _, structField := range reflection.GetStructFields(entity.Type()) {
		// Does this case ever happen?
		// Would conflict with anonymous pointer fields (Go "inheritance")
		//		if reflection.IsPointerToStruct(structField.Type) {
		//			// Compatible with *any
		//			field := entity.FieldByName(structField.Name)
		//			if !field.IsNil() {
		//				e := field.Interface()
		//				// We sometimes have pointers to non-entities, so make sure
		//				if _, ok := e.(tosca.Contextual); ok {
		//					dependencies[e] = true
		//				}
		//			}
		//		}

		if reflection.IsMapOfStringToPointerToStruct(structField.Type) {
			// Compatible with map[string]*any
			field := entity.FieldByName(structField.Name)
			for _, mapKey := range field.MapKeys() {
				element := field.MapIndex(mapKey)
				dependencies.Add(element.Interface())
			}
		} else if reflection.IsSliceOfPointerToStruct(structField.Type) {
			// Compatible with []*any
			field := entity.FieldByName(structField.Name)
			length := field.Len()
			for i := 0; i < length; i++ {
				element := field.Index(i)
				dependencies.Add(element.Interface())
			}
		}
	}

	return dependencies
}

//
// InheritField
//

type InheritField struct {
	Entity          reflect.Value
	Key             string
	Field           reflect.Value
	ParentEntityPtr parsing.EntityPtr
	ParentField     reflect.Value
}

func (self *InheritField) Inherit() {
	// TODO: do we really need all of these? some of them aren't used in TOSCA
	fieldEntityPtr := self.Field.Interface()
	if reflection.IsPointerToString(fieldEntityPtr) {
		self.InheritEntity()
	} else if reflection.IsPointerToInt64(fieldEntityPtr) {
		self.InheritEntity()
	} else if reflection.IsPointerToBool(fieldEntityPtr) {
		self.InheritEntity()
	} else if reflection.IsPointerToSliceOfString(fieldEntityPtr) {
		self.InheritStringsFromSlice()
	} else if reflection.IsPointerToMapOfStringToString(fieldEntityPtr) {
		self.InheritStringsFromMap()
	} else {
		fieldType := self.Field.Type()
		if reflection.IsPointerToStruct(fieldType) {
			self.InheritEntity()
		} else if reflection.IsSliceOfPointerToStruct(fieldType) {
			self.InheritStructsFromSlice()
		} else if reflection.IsMapOfStringToPointerToStruct(fieldType) {
			self.InheritStructsFromMap()
		} else {
			panic(fmt.Sprintf("\"inherit\" tag's field type %q is not supported in struct: %T", fieldType, self.Entity.Interface()))
		}
	}
}

// Field is compatible with *any
func (self *InheritField) InheritEntity() {
	if self.Field.IsNil() && !self.ParentField.IsNil() {
		self.Field.Set(self.ParentField)
	}
}

// Field is *[]string
func (self *InheritField) InheritStringsFromSlice() {
	slicePtr := self.Field.Interface().(*[]string)
	parentSlicePtr := self.ParentField.Interface().(*[]string)

	if parentSlicePtr == nil {
		return
	}

	parentSlice := *parentSlicePtr
	length := len(parentSlice)
	if length == 0 {
		return
	}

	var slice []string
	if slicePtr != nil {
		slice = *slicePtr
	} else {
		slice = make([]string, 0, length)
	}

	for _, element := range parentSlice {
		slice = append(slice, element)
	}

	self.Field.Set(reflect.ValueOf(&slice))
}

// Field is compatible with []*any
func (self *InheritField) InheritStructsFromSlice() {
	slice := self.Field

	length := self.ParentField.Len()
	for i := 0; i < length; i++ {
		element := self.ParentField.Index(i)

		if _, ok := element.Interface().(parsing.Mappable); ok {
			// For mappable elements only, *don't* inherit the same key
			// (We'll merge everything else)
			key := parsing.GetKey(element.Interface())
			if ii, ok := getSliceElementIndexForKey(self.Field, key); ok {
				element_ := self.Field.Index(ii)
				logInheritance.Debugf("override: %s", parsing.GetContext(element_.Interface()).Path)
				continue // skip this key
			}
		}

		slice = reflect.Append(slice, element)
	}

	self.Field.Set(slice)
}

// Field is *map[string]string
func (self *InheritField) InheritStringsFromMap() {
	mapPtr := self.Field.Interface().(*map[string]string)
	parentMapPtr := self.ParentField.Interface().(*map[string]string)

	if parentMapPtr == nil {
		return
	}

	parentMap := *parentMapPtr
	if length := len(parentMap); length == 0 {
		return
	}

	var map_ map[string]string
	if mapPtr != nil {
		map_ = *mapPtr
	} else {
		map_ = make(map[string]string)
	}

	for key, value := range parentMap {
		if _, ok := map_[key]; !ok {
			map_[key] = value
		}
	}

	self.Field.Set(reflect.ValueOf(&map_))
}

// Field is compatible with map[string]*any
func (self *InheritField) InheritStructsFromMap() {
	for _, mapKey := range self.ParentField.MapKeys() {
		element := self.ParentField.MapIndex(mapKey)
		element_ := self.Field.MapIndex(mapKey)
		if element_.IsValid() {
			// We are overriding this element, so don't inherit it
			logInheritance.Debugf("override: %s", parsing.GetContext(element_.Interface()).Path)
		} else {
			self.Field.SetMapIndex(mapKey, element)
		}
	}
}

func getSliceElementIndexForKey(slice reflect.Value, key string) (int, bool) {
	length := slice.Len()
	for i := 0; i < length; i++ {
		element := slice.Index(i)
		if parsing.GetKey(element.Elem()) == key {
			return i, true
		}
	}
	return -1, false
}

//
// InheritFields
//

type InheritFields map[parsing.EntityPtr][]*InheritField

// From "inherit" tags
func NewInheritFields(entityPtr parsing.EntityPtr) []*InheritField {
	var inheritFields []*InheritField

	entity := reflect.ValueOf(entityPtr).Elem()
	for fieldName, tag := range reflection.GetFieldTagsForValue(entity, "inherit") {
		key, referenceFieldName := parseInheritTag(tag)

		if referenceField, referredField, ok := reflection.GetReferredField(entity, referenceFieldName, fieldName); ok {
			field := entity.FieldByName(fieldName)
			inheritFields = append(inheritFields, &InheritField{entity, key, field, referenceField.Interface(), referredField})
		}
	}

	return inheritFields
}

func (self InheritFields) Get(entityPtr parsing.EntityPtr) []*InheritField {
	// TODO: cache these, because we call twice for each entity
	inheritFields, ok := self[entityPtr]
	if !ok {
		inheritFields = NewInheritFields(entityPtr)
		self[entityPtr] = inheritFields
	}
	return inheritFields
}

func parseInheritTag(tag string) (string, string) {
	t := strings.Split(tag, ",")
	if len(t) != 2 {
		panic("must be 2")
	}

	key := t[0]
	referenceFieldName := t[1]

	return key, referenceFieldName
}
