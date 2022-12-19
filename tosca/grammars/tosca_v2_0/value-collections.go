package tosca_v2_0

import (
	"github.com/tliron/kutil/ard"
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
)

//
// ValueList
//

type ValueList struct {
	EntryConstraints ConstraintClauses
	Slice            []any
}

func NewValueList(length int, entryConstraints ConstraintClauses) *ValueList {
	return &ValueList{
		EntryConstraints: entryConstraints,
		Slice:            make([]any, length),
	}
}

func (self *ValueList) Set(index int, value any) {
	self.Slice[index] = value
}

func (self *ValueList) Normalize(context *tosca.Context) *normal.List {
	normalList := normal.NewList(len(self.Slice))

	for index, value := range self.Slice {
		if value_, ok := value.(*Value); ok {
			normalList.Set(index, value_.normalize(true))
		} else {
			normalList.Set(index, normal.NewPrimitive(value))
		}
	}

	return normalList
}

//
// ValueMap
//

type ValueMap struct {
	KeyConstraints   ConstraintClauses
	ValueConstraints ConstraintClauses
	Map              ard.Map
}

func NewValueMap(keyConstraints ConstraintClauses, valueConstraints ConstraintClauses) *ValueMap {
	return &ValueMap{
		KeyConstraints:   keyConstraints,
		ValueConstraints: valueConstraints,
		Map:              make(ard.Map),
	}
}

func (self *ValueMap) Put(key any, value any) {
	self.Map[key] = value
}

func (self *ValueMap) Normalize(context *tosca.Context) *normal.Map {
	normalMap := normal.NewMap()

	for key, value := range self.Map {
		if key_, ok := key.(*Value); ok {
			key = key_.normalize(true)
		}
		if value_, ok := value.(*Value); ok {
			normalMap.Put(key, value_.normalize(true))
		} else {
			normalMap.Put(key, normal.NewPrimitive(value))
		}
	}

	return normalMap
}

// Utils

func GetListSchemas(dataType *DataType, dataDefinition DataDefinition) *Schema {
	var entrySchema *Schema
	if dataDefinition != nil {
		entrySchema = dataDefinition.GetEntrySchema()
	}
	if entrySchema == nil {
		entrySchema = dataType.EntrySchema
	}
	if (entrySchema != nil) && (entrySchema.DataType != nil) {
		return entrySchema
	} else {
		return nil
	}
}

func GetMapSchemas(dataType *DataType, dataDefinition DataDefinition) (*Schema, *Schema) {
	var keySchema *Schema
	var valueSchema *Schema
	if dataDefinition != nil {
		keySchema = dataDefinition.GetKeySchema()
		valueSchema = dataDefinition.GetEntrySchema()
	}
	if keySchema == nil {
		keySchema = dataType.KeySchema
	}
	if valueSchema == nil {
		valueSchema = dataType.EntrySchema
	}

	if (keySchema != nil) && (keySchema.DataType != nil) && (valueSchema != nil) && (valueSchema.DataType != nil) {
		return keySchema, valueSchema
	} else {
		return nil, nil
	}
}
