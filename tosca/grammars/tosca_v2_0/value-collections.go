package tosca_v2_0

import (
	"github.com/tliron/go-ard"
	"github.com/tliron/puccini/normal"
	"github.com/tliron/puccini/tosca/parsing"
)

//
// ValueList
//

type ValueList struct {
	ValidationClause *ValidationClause
	Slice            []any
}

func NewValueList(length int, validationClause *ValidationClause) *ValueList {
	return &ValueList{
		ValidationClause: validationClause,
		Slice:            make([]any, length),
	}
}

func (self *ValueList) Set(index int, value any) {
	self.Slice[index] = value
}

func (self *ValueList) Normalize(context *parsing.Context) *normal.List {
	normalList := normal.NewList(len(self.Slice))

	for index, value := range self.Slice {
		if value_, ok := value.(*Value); ok {
			bare := value_.Meta == nil || value_.Meta.Empty()
			normalList.Set(index, value_.normalize(bare))
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
	KeyValidation   *ValidationClause
	ValueValidation *ValidationClause
	Map             ard.Map
}

func NewValueMap(keyValidation *ValidationClause, valueValidation *ValidationClause) *ValueMap {
	return &ValueMap{
		KeyValidation:   keyValidation,
		ValueValidation: valueValidation,
		Map:             make(ard.Map),
	}
}

func (self *ValueMap) Put(key any, value any) {
	self.Map[key] = value
}

func (self *ValueMap) Normalize(context *parsing.Context) *normal.Map {
	normalMap := normal.NewMap()

	for key, value := range self.Map {
		if key_, ok := key.(*Value); ok {
			keyBare := key_.Meta == nil || key_.Meta.Empty()
			key = key_.normalize(keyBare)
		}
		if value_, ok := value.(*Value); ok {
			valueBare := value_.Meta == nil || value_.Meta.Empty()
			normalMap.Put(key, value_.normalize(valueBare))
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
