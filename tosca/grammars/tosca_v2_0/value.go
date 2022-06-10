package tosca_v2_0

import (
	"fmt"

	"github.com/tliron/kutil/ard"
	"github.com/tliron/kutil/util"
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
	"github.com/tliron/yamlkeys"
)

//
// Value
//
// (attribute, property, and parameter assignments)
//
// [TOSCA-v2.0] @ ?
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.6.11, 3.6.13
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.10, 3.6.12
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.9, 3.5.11
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.5.9, 3.5.11
//

type Value struct {
	*Entity `name:"value"`
	Name    string

	ConstraintClauses ConstraintClauses
	Description       *string

	DataType    *DataType                `traverse:"ignore" json:"-" yaml:"-"`
	Information *normal.ValueInformation `traverse:"ignore" json:"-" yaml:"-"`
	Converter   *tosca.FunctionCall      `traverse:"ignore" json:"-" yaml:"-"`
}

func NewValue(context *tosca.Context) *Value {
	return &Value{
		Entity:      NewEntity(context),
		Name:        context.Name,
		Information: normal.NewValueInformation(),
	}
}

// tosca.Reader signature
func ReadValue(context *tosca.Context) tosca.EntityPtr {
	ToFunctionCall(context)
	return NewValue(context)
}

// tosca.Reader signature
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.12.2.2
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.11.2.2
func ReadAttributeValue(context *tosca.Context) tosca.EntityPtr {
	self := NewValue(context)

	// Unpack long notation
	if context.Is(ard.TypeMap) {
		map_ := context.Data.(ard.Map)
		if len(map_) == 2 {
			if description, ok := map_["description"]; ok {
				if value, ok := map_["value"]; ok {
					self.Description = context.FieldChild("description", description).ReadString()
					context.Data = value
				}
			}
		}
	}

	ToFunctionCall(context)

	return self
}

// tosca.Mappable interface
func (self *Value) GetKey() string {
	return self.Name
}

// fmt.Stringer interface
func (self *Value) String() string {
	return yamlkeys.KeyString(self.Context.Data)
}

func (self *Value) RenderDataType(dataTypeName string) {
	if e, ok := self.Context.Namespace.Lookup(dataTypeName); ok {
		if dataType, ok := e.(*DataType); ok {
			self.Render(dataType, nil, false, false)
		} else {
			self.Context.ReportUnknownDataType(dataTypeName)
		}
	} else {
		self.Context.ReportUnknownDataType(dataTypeName)
	}
}

// Avoid rendering more than once (can happen if we were copied from PropertyDefinition.Default)
func (self *Value) Render(dataType *DataType, definition *AttributeDefinition, bare bool, allowNil bool) {
	self.DataType = dataType

	if !bare {
		if self.Description != nil {
			self.Information.Description = *self.Description
		}

		if definition != nil {
			self.Information.Definition = definition.GetTypeInformation()
		}

		if dataType != nil {
			self.Information.Type = dataType.GetTypeInformation()
		}
	}

	dataType.Complete(self.Context)
	if !bare {
		self.ConstraintClauses.Render(dataType, definition)
		dataType.ConstraintClauses.Render(dataType, definition)
		self.ConstraintClauses = dataType.ConstraintClauses.Append(self.ConstraintClauses)
	}

	if _, ok := self.Context.Data.(*tosca.FunctionCall); ok {
		return
	}

	if allowNil && (self.Context.Data == nil) {
		return
	}

	// Internal types
	if internalTypeName, typeValidator, reader, ok := dataType.GetInternal(); ok {
		if typeValidator != nil {
			if self.Context.Data == nil {
				// Nil data only happens when an attribute is added despite not having a
				// "default" value; we will give it a valid zero value instead
				if self.Context.Data, ok = ScalarUnitTypeZeroes[internalTypeName]; !ok {
					if self.Context.Data, ok = ard.TypeZeroes[internalTypeName]; !ok {
						panic(fmt.Sprintf("unsupported internal type name: %s", internalTypeName))
					}
				}
			}

			if (internalTypeName == ard.TypeString) && self.Context.HasQuirk(tosca.QuirkDataTypesStringPermissive) {
				self.Context.Data = ard.ValueToString(self.Context.Data)
			}

			// Primitive types
			if typeValidator(self.Context.Data) {
				// Render list and map elements according to entry schema
				// (The entry schema may also have additional constraints)
				switch internalTypeName {
				case ard.TypeList, ard.TypeMap:
					var entrySchema *Schema
					if definition != nil {
						entrySchema = definition.EntrySchema
					} else {
						entrySchema = dataType.EntrySchema
					}

					if entrySchema == nil {
						return
					} else if entrySchema.DataType == nil {
						return
					}

					if internalTypeName == ard.TypeList {
						// Information
						entryDataType := entrySchema.DataType
						entryConstraints := entrySchema.GetConstraints()
						self.Information.Entry = entryDataType.GetTypeInformation()
						if entrySchema.Description != nil {
							self.Information.Entry.SchemaDescription = *entrySchema.Description
						}

						slice := self.Context.Data.(ard.List)
						valueList := NewValueList(definition, len(slice), entryConstraints)

						for index, data := range slice {
							value := ReadAndRenderBare(self.Context.ListChild(index, data), entryDataType)
							valueList.Set(index, value)
						}

						self.Context.Data = valueList
					} else { // ard.TypeMap
						var keySchema *Schema
						if definition != nil {
							keySchema = definition.KeySchema
						} else {
							keySchema = dataType.KeySchema
						}

						if keySchema == nil {
							return
						} else if keySchema.DataType == nil {
							return
						}

						// Information

						keyDataType := keySchema.DataType
						keyConstraints := keySchema.GetConstraints()
						self.Information.Key = keyDataType.GetTypeInformation()
						if keySchema.Description != nil {
							self.Information.Key.SchemaDescription = *keySchema.Description
						}

						valueDataType := entrySchema.DataType
						valueConstraints := entrySchema.GetConstraints()
						self.Information.Value = valueDataType.GetTypeInformation()
						if entrySchema.Description != nil {
							self.Information.Value.SchemaDescription = *entrySchema.Description
						}

						valueMap := NewValueMap(definition, keyConstraints, valueConstraints)

						for key, data := range self.Context.Data.(ard.Map) {
							// Complex keys are stringified for the purpose of the contexts

							// Validate key schema
							keyContext := self.Context.MapChild(key, yamlkeys.KeyData(key))
							key = ReadAndRenderBare(keyContext, keyDataType)

							context := self.Context.MapChild(key, data)
							value := ReadAndRenderBare(context, valueDataType)
							value.ConstraintClauses = ConstraintClauses{}
							valueMap.Put(key, value)
						}

						self.Context.Data = valueMap
					}
				}
			} else {
				self.Context.ReportValueWrongType(internalTypeName)
			}
		} else {
			// Special types
			self.Context.Data = reader(self.Context)
		}
	} else if self.Context.ValidateType(ard.TypeMap) {
		// Complex data types

		map_ := self.Context.Data.(ard.Map)

		// All properties must be defined in type
		for key := range map_ {
			name := yamlkeys.KeyString(key)
			if _, ok := dataType.PropertyDefinitions[name]; !ok {
				self.Context.MapChild(name, nil).ReportUndeclared("property")
				delete(map_, key)
			}
		}

		// Render properties
		for key, definition := range dataType.PropertyDefinitions {
			definition.Render()
			if data, ok := map_[key]; ok {
				var value *Value
				if value, ok = data.(*Value); !ok {
					// Convert to value
					value = ReadValue(self.Context.MapChild(key, data)).(*Value)
					map_[key] = value
				}
				if definition.DataType != nil {
					var lock util.RWLocker
					value.RenderProperty(definition.DataType, definition)
					if lock != nil {
						lock.RUnlock()
					}
				}

				// Grab information
				if (value.Information != nil) && !value.Information.Empty() {
					self.Information.Fields[key] = value.Information
				}
			} else if definition.IsRequired() {
				self.Context.MapChild(key, data).ReportValueRequired("property")
			}
		}
	}

	if (definition != nil) && (definition.Metadata != nil) {
		if converter, ok := definition.Metadata["puccini.converter"]; ok {
			self.Converter = self.Context.NewFunctionCall(converter, nil)
		}
	}
	if self.Converter == nil {
		if converter, ok := dataType.GetMetadataValue("puccini.converter"); ok {
			self.Converter = self.Context.NewFunctionCall(converter, nil)
		}
	}

	if comparer, ok := dataType.GetMetadataValue("puccini.comparer"); ok {
		if hasComparer, ok := self.Context.Data.(HasComparer); ok {
			hasComparer.SetComparer(comparer)
		} else {
			panic(fmt.Sprintf("type has \"puccini.comparer\" metadata but does not support HasComparer interface: %T", self.Context.Data))
		}
	}
}

func ReadAndRenderBare(context *tosca.Context, dataType *DataType) *Value {
	self := ReadValue(context).(*Value)
	self.Render(dataType, nil, true, false)
	return self
}

func (self *Value) RenderProperty(dataType *DataType, definition *PropertyDefinition) {
	if definition == nil {
		self.Render(dataType, nil, false, false)
	} else {
		self.ConstraintClauses.Render(dataType, definition.AttributeDefinition)
		definition.ConstraintClauses.Render(definition.DataType, definition.AttributeDefinition)
		self.ConstraintClauses = definition.ConstraintClauses.Append(self.ConstraintClauses)
		self.Render(dataType, definition.AttributeDefinition, false, false)
	}
}

func (self *Value) Normalize() normal.Constrainable {
	return self.normalize(true)
}

func (self *Value) normalize(withInformation bool) normal.Constrainable {
	var normalConstrainable normal.Constrainable

	switch data := self.Context.Data.(type) {
	case ard.Map:
		// This is for complex types (the "map" type is a ValueMap, below)
		normalMap := normal.NewMap()
		for key, value := range data {
			if v, ok := value.(*Value); ok {
				normalMap.Put(key, v.normalize(true))
			} else {
				normalMap.Put(key, normal.NewValue(value))
			}
		}
		normalConstrainable = normalMap

	case *ValueList:
		normalConstrainable = data.Normalize(self.Context)

	case *ValueMap:
		normalConstrainable = data.Normalize(self.Context)

	case *tosca.FunctionCall:
		NormalizeFunctionCallArguments(data, self.Context)
		normalConstrainable = normal.NewFunctionCall(data)

	default:
		value := normal.NewValue(data)
		normalConstrainable = value
	}

	if withInformation {
		normalConstrainable.SetInformation(self.Information)
	}

	self.ConstraintClauses.NormalizeConstrainable(self.Context, normalConstrainable)

	if self.Converter != nil {
		normalConstrainable.SetConverter(self.Converter)
	}

	return normalConstrainable
}

//
// Values
//

type Values map[string]*Value

func (self Values) CopyUnassigned(values Values) {
	for key, value := range values {
		if _, ok := self[key]; !ok {
			self[key] = value
		}
	}
}

func (self Values) RenderAttributes(definitions AttributeDefinitions, context *tosca.Context) {
	for key, definition := range definitions {
		definition.Render()
		if _, ok := self[key]; !ok {
			if definition.Default != nil {
				self[definition.Name] = definition.Default
			} else {
				// Attributes should always appear, at least as nil, even if they have no default value
				self[definition.Name] = NewValue(context.MapChild(definition.Name, nil))
			}
		}
	}

	for key, value := range self {
		if definition, ok := definitions[key]; ok {
			// Avoid re-rendering default
			if value != definition.Default {
				if definition.DataType != nil {
					value.Render(definition.DataType, definition, false, true)
				}
			}
		} else {
			value.Context.ReportUndeclared("attribute")
			delete(self, key)
		}
	}
}

func (self Values) RenderProperties(definitions PropertyDefinitions, context *tosca.Context) {
	for key, definition := range definitions {
		definition.Render()
		if _, ok := self[key]; !ok {
			if definition.Default != nil {
				self[definition.Name] = definition.Default
			} else if definition.IsRequired() {
				context.MapChild(definition.Name, nil).ReportValueRequired("property")
			}
		}
	}

	for key, value := range self {
		if definition, ok := definitions[key]; ok {
			// Avoid re-rendering default
			if value != definition.Default {
				if definition.DataType != nil {
					value.RenderProperty(definition.DataType, definition)
				}
			}
		} else {
			value.Context.ReportUndeclared("property")
			delete(self, key)
		}
	}
}

func (self Values) RenderInputs(definitions ParameterDefinitions, context *tosca.Context) {
	for key, definition := range definitions {
		definition.Render()
		if _, ok := self[key]; !ok {
			if definition.Value != nil {
				self[definition.Name] = definition.Value
			} else if definition.IsRequired() {
				context.MapChild(definition.Name, nil).ReportValueRequired("input")
			}
		}
	}

	for key, value := range self {
		if definition, ok := definitions[key]; ok {
			// Avoid re-rendering
			if value != definition.Value {
				if definition.DataType != nil {
					value.RenderProperty(definition.DataType, definition.PropertyDefinition)
				}
			}
		}
		// No "else": we allow ad-hoc input assignments
	}
}

func (self Values) Normalize(normalConstrainables normal.Constrainables) {
	for key, value := range self {
		normalConstrainables[key] = value.Normalize()
	}
}

//
// ValueList
//

type ValueList struct {
	EntryConstraints ConstraintClauses
	Slice            []any
}

func NewValueList(definition *AttributeDefinition, length int, entryConstraints ConstraintClauses) *ValueList {
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

	self.EntryConstraints.NormalizeListEntries(context, normalList)

	for index, value := range self.Slice {
		if v, ok := value.(*Value); ok {
			normalList.Set(index, v.normalize(false))
		} else {
			normalList.Set(index, normal.NewValue(value))
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

func NewValueMap(definition *AttributeDefinition, keyConstraints ConstraintClauses, valueConstraints ConstraintClauses) *ValueMap {
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

	self.KeyConstraints.NormalizeMapKeys(context, normalMap)
	self.ValueConstraints.NormalizeMapValues(context, normalMap)

	for key, value := range self.Map {
		if k, ok := key.(*Value); ok {
			key = k.normalize(false)
		}
		if v, ok := value.(*Value); ok {
			normalMap.Put(key, v.normalize(false))
		} else {
			normalMap.Put(key, normal.NewValue(value))
		}
	}

	return normalMap
}
