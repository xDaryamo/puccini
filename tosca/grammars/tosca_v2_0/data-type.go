package tosca_v2_0

import (
	"github.com/tliron/go-ard"
	"github.com/tliron/puccini/normal"
	"github.com/tliron/puccini/tosca/parsing"
)

type HasComparer interface {
	SetComparer(comparer string)
}

//
// DataType
//
// [TOSCA-v2.0] @ 9.11
//

type DataType struct {
	*Type `name:"data type"`

	PropertyDefinitions PropertyDefinitions `read:"properties,PropertyDefinition" inherit:"properties,Parent"`
	ValidationClause    *ValidationClause   `read:"validation,ValidationClause" traverse:"ignore"`
	KeySchema           *Schema             `read:"key_schema,Schema"`   // introduced in TOSCA 1.3
	EntrySchema         *Schema             `read:"entry_schema,Schema"` // introduced in TOSCA 1.3

	// TOSCA 2.0 scalar type fields - usa reader custom
	DataTypeName  *string `read:"data_type"`               // introduced in TOSCA 2.0
	Units         ard.Map `read:"units,UnitsReader"`       // introduced in TOSCA 2.0
	CanonicalUnit *string `read:"canonical_unit"`          // introduced in TOSCA 2.0
	Prefixes      ard.Map `read:"prefixes,PrefixesReader"` // introduced in TOSCA 2.0

	Parent *DataType `lookup:"derived_from,ParentName" traverse:"ignore" json:"-" yaml:"-"`
}

// Custom reader for units field
func ReadUnitsField(context *parsing.Context) parsing.EntityPtr {
	if context.ValidateType(ard.TypeMap) {
		return context.Data.(ard.Map)
	}
	return nil
}

// Custom reader for prefixes field
func ReadPrefixesField(context *parsing.Context) parsing.EntityPtr {
	if context.ValidateType(ard.TypeMap) {
		return context.Data.(ard.Map)
	}
	return nil
}

func NewDataType(context *parsing.Context) *DataType {
	return &DataType{
		Type:                NewType(context),
		PropertyDefinitions: make(PropertyDefinitions),
	}
}

// ([parsing.Reader] signature)
func ReadDataType(context *parsing.Context) parsing.EntityPtr {
	self := NewDataType(context)

	// Read fields normally
	context.ValidateUnsupportedFields(context.ReadFields(self))

	// Manually handle Units and Prefixes fields
	if context.ValidateType(ard.TypeMap) {
		if data, ok := context.Data.(ard.Map)["units"]; ok {
			if unitsMap, ok := data.(ard.Map); ok {
				self.Units = unitsMap
			}
		}

		if data, ok := context.Data.(ard.Map)["prefixes"]; ok {
			if prefixesMap, ok := data.(ard.Map); ok {
				self.Prefixes = prefixesMap
			}
		}
	}

	return self
}

// ([parsing.Hierarchical] interface)
func (self *DataType) GetParent() parsing.EntityPtr {
	return self.Parent
}

// ([parsing.Inherits] interface)
func (self *DataType) Inherit() {
	logInherit.Debugf("data type: %s", self.Name)

	if _, ok := self.GetInternalTypeName(); ok && (len(self.PropertyDefinitions) > 0) {
		// Doesn't make sense to be an internal type (non-complex) and also have properties (complex)
		self.Context.ReportPrimitiveType()
		self.PropertyDefinitions = make(PropertyDefinitions)
		return
	}

	if self.Parent == nil {
		return
	}

	if (self.KeySchema == nil) && (self.Parent.KeySchema != nil) {
		self.KeySchema = self.Parent.KeySchema
	}
	if (self.EntrySchema == nil) && (self.Parent.EntrySchema != nil) {
		self.EntrySchema = self.Parent.EntrySchema
	}
	if (self.ValidationClause == nil) && (self.Parent.ValidationClause != nil) {
		self.ValidationClause = self.Parent.ValidationClause
	}

	// TOSCA 2.0 scalar inheritance
	if (self.DataTypeName == nil) && (self.Parent.DataTypeName != nil) {
		self.DataTypeName = self.Parent.DataTypeName
	}
	if (self.Units == nil) && (self.Parent.Units != nil) {
		self.Units = make(ard.Map)
		for k, v := range self.Parent.Units {
			self.Units[k] = v
		}
	}
	if (self.CanonicalUnit == nil) && (self.Parent.CanonicalUnit != nil) {
		self.CanonicalUnit = self.Parent.CanonicalUnit
	}
	if (self.Prefixes == nil) && (self.Parent.Prefixes != nil) {
		self.Prefixes = make(ard.Map)
		for k, v := range self.Parent.Prefixes {
			self.Prefixes[k] = v
		}
	}

	self.PropertyDefinitions.Inherit(self.Parent.PropertyDefinitions)
}

// Check if this is a scalar type
func (self *DataType) IsScalarType() bool {
	// Check if derived from scalar
	current := self
	for current != nil {
		if current.Name == "scalar" {
			return true
		}
		current = current.Parent
	}
	return false
}

// ([parsing.Renderable] interface)
func (self *DataType) Render() {
	// Avoid rendering more than once
	self.renderOnce.Do(self.render)
}

func (self *DataType) render() {
	logRender.Debugf("data type: %s", self.Name)

	// Verify that if it's internal type that it is supported
	if internalTypeName, ok := self.GetInternalTypeName(); ok {
		if _, ok := ard.TypeValidators[internalTypeName]; !ok {
			if _, ok := self.Context.Grammar.Readers[string(internalTypeName)]; !ok {
				self.Context.ReportUnsupportedType()
			}
		}
	}
}

func (self *DataType) GetInternalTypeName() (ard.TypeName, bool) {
	if typeName, ok := self.GetMetadataValue(parsing.MetadataType); ok {
		return ard.TypeName(typeName), ok
	} else if self.Parent != nil {
		// The internal type metadata is inherited
		return self.Parent.GetInternalTypeName()
	} else {
		return ard.NoType, false
	}
}

type DataTypeInternal struct {
	Name      ard.TypeName
	Validator ard.TypeValidator
	Reader    parsing.Reader
}

func (self *DataType) GetInternal() *DataTypeInternal {
	if internalTypeName, ok := self.GetInternalTypeName(); ok {
		if typeValidator, ok := ard.TypeValidators[internalTypeName]; ok {
			return &DataTypeInternal{internalTypeName, typeValidator, nil}
		} else if reader, ok := self.Context.Grammar.Readers[string(internalTypeName)]; ok {
			return &DataTypeInternal{internalTypeName, nil, reader}
		}
	}
	return nil
}

// Note that this may change the data (if it's a map), but that should be fine, because we intend
// for the data to be complete. For the same reason, this action is idempotent (subsequent calls with
// the same data will not have an effect).
func (self *DataType) CompleteData(context *parsing.Context) {
	map_, ok := context.Data.(ard.Map)
	if !ok {
		// Only for complex data types
		return
	}

	for key, definition := range self.PropertyDefinitions {
		childContext := context.MapChild(key, nil)

		if data, ok := map_[key]; ok {
			childContext.Data = data
		} else if definition.Default != nil {
			// Assign default value
			childContext.Data = definition.Default.Context.Data
			ParseFunctionCall(childContext)
			map_[key] = childContext.Data
		}

		if definition.DataType != nil {
			definition.DataType.CompleteData(childContext)
		}
	}
}

func (self *DataType) NewValueMeta() *normal.ValueMeta {
	normalValueMeta := normal.NewValueMeta()
	normalValueMeta.Type = parsing.GetCanonicalName(self)
	normalValueMeta.TypeMetadata = parsing.GetDataTypeMetadata(self.Metadata)
	if self.Description != nil {
		normalValueMeta.TypeDescription = *self.Description
	}
	return normalValueMeta
}

func LookupDataType(context *parsing.Context, name string) (*DataType, bool) {
	if entityPtr, ok := context.Namespace.LookupForType(name, dataTypePtrType); ok {
		return entityPtr.(*DataType), true
	} else {
		return nil, false
	}
}

//
// DataTypes
//

type DataTypes []*DataType
