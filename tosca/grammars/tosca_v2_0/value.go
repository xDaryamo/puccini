package tosca_v2_0

import (
	"fmt"

	"github.com/tliron/go-ard"
	"github.com/tliron/puccini/normal"
	"github.com/tliron/puccini/tosca/parsing"
	"github.com/tliron/yamlkeys"
)

//
// Value
//
// (attribute, property, and parameter assignments)
//
// [TOSCA-v2.0] @ 9.11
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.6.11, 3.6.13
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.10, 3.6.12
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.9, 3.5.11
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.5.9, 3.5.11
//

type Value struct {
	*Entity `name:"value"`
	Name    string

	ValidationClause *ValidationClause
	Description      *string // not used since TOSCA 2.0

	DataType *DataType         `traverse:"ignore" json:"-" yaml:"-"`
	Meta     *normal.ValueMeta `traverse:"ignore" json:"-" yaml:"-"`
}

func NewValue(context *parsing.Context) *Value {
	return &Value{
		Entity: NewEntity(context),
		Name:   context.Name,
		Meta:   normal.NewValueMeta(),
	}
}

// ([parsing.Reader] signature)
func ReadValue(context *parsing.Context) parsing.EntityPtr {
	ParseFunctionCall(context)
	return NewValue(context)
}

func ReadAndRenderBare(context *parsing.Context, dataType *DataType, definition DataDefinition) *Value {
	self := ReadValue(context).(*Value)
	self.render(dataType, definition, true, false)
	return self
}

// ([parsing.Mappable] interface)
func (self *Value) GetKey() string {
	return self.Name
}

// ([fmt.Stringer] interface)
func (self *Value) String() string {
	return yamlkeys.KeyString(self.Context.Data)
}

func (self *Value) RenderDataType(dataTypeName string) {
	if entityPtr, ok := self.Context.Namespace.Lookup(dataTypeName); ok {
		if dataType, ok := entityPtr.(*DataType); ok {
			self.Render(dataType, nil, false, false)
		} else {
			self.Context.ReportUnknownDataType(dataTypeName)
		}
	} else {
		self.Context.ReportUnknownDataType(dataTypeName)
	}
}

func (self *Value) Render(dataType *DataType, dataDefinition DataDefinition, bare bool, allowNil bool) {
	// Avoid rendering more than once
	self.renderOnce.Do(func() {
		self.render(dataType, dataDefinition, bare, allowNil)
	})
}

// "bare" means without meta
func (self *Value) render(dataType *DataType, dataDefinition DataDefinition, bare bool, allowNil bool) {
	self.DataType = dataType

	dataType.CompleteData(self.Context)

	if !bare {
		self.Meta = NewValueMeta(self.Context, dataType, dataDefinition, self.ValidationClause)

		// Not used since TOSCA 2.0
		if self.Description != nil {
			self.Meta.LocalDescription = *self.Description
		}
	}

	if _, ok := self.Context.Data.(*parsing.FunctionCall); ok {
		return
	}

	if allowNil && (self.Context.Data == nil) {
		return
	}

	// Check if this is a scalar type (TOSCA 2.0)
	if dataType.IsScalarType() {
		self.Context.Data = ReadScalar(self.Context, dataType)
		return
	}

	// Internal types
	if internal := dataType.GetInternal(); internal != nil {
		if internal.Validator != nil {
			if self.Context.Data == nil {
				// Nil data only happens when an attribute is added despite not having a
				// "default" value; we will give it a valid zero value instead
				var ok bool
				if self.Context.Data, ok = ScalarUnitTypeZeroes[internal.Name]; !ok {
					if self.Context.Data, ok = ard.TypeZeroes[internal.Name]; !ok {
						panic(fmt.Sprintf("unsupported internal type name: %s", internal.Name))
					}
				}
			}

			if (internal.Name == ard.TypeString) && self.Context.HasQuirk(parsing.QuirkDataTypesStringPermissive) {
				self.Context.Data = ard.ValueToString(self.Context.Data)
			}

			// Primitive types
			if internal.Validator(self.Context.Data) {
				// Render list and map elements according to entry schema
				// (The entry schema may also have additional constraints)
				switch internal.Name {
				case ard.TypeList:
					if entrySchema := GetListSchemas(dataType, dataDefinition); entrySchema != nil {
						slice := self.Context.Data.(ard.List)
						valueList := NewValueList(len(slice), entrySchema.GetValidation())

						for index, data := range slice {
							value := ReadAndRenderBare(self.Context.ListChild(index, data), entrySchema.DataType, entrySchema)
							valueList.Set(index, value)
						}

						self.Context.Data = valueList
					}

				case ard.TypeMap:
					if keySchema, valueSchema := GetMapSchemas(dataType, dataDefinition); (keySchema != nil) && (valueSchema != nil) {
						valueMap := NewValueMap(keySchema.GetValidation(), valueSchema.GetValidation())

						for key, data := range self.Context.Data.(ard.Map) {
							// Complex keys are stringified for the purpose of the contexts
							key = ReadAndRenderBare(self.Context.MapChild(key, yamlkeys.KeyData(key)), keySchema.DataType, keySchema)
							value := ReadAndRenderBare(self.Context.MapChild(key, data), valueSchema.DataType, valueSchema)
							valueMap.Put(key, value)
						}

						self.Context.Data = valueMap
					}
				}
			} else {
				self.Context.ReportValueWrongType(internal.Name)
			}
		} else {
			// Special types
			self.Context.Data = internal.Reader(self.Context)
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
					value.RenderProperty(definition.DataType, definition)
				}

				// Move meta
				if (value.Meta != nil) && !value.Meta.Empty() {
					self.Meta.Fields[key] = value.Meta
					value.Meta = nil
				}
			} else if definition.IsRequired() {
				self.Context.MapChild(key, data).ReportValueRequired("property")
			}
		}
	}

	if self.Meta != nil {
		var converter *parsing.FunctionCall
		if dataDefinition != nil {
			if metadata := dataDefinition.GetTypeMetadata(); metadata != nil {
				if converter_, ok := metadata[parsing.MetadataConverter]; ok {
					converter = self.Context.NewFunctionCall(converter_, nil)
				}
			}
		}
		if converter == nil {
			if converter_, ok := dataType.GetMetadataValue(parsing.MetadataConverter); ok {
				converter = self.Context.NewFunctionCall(converter_, nil)
			}
		}
		if converter != nil {
			self.Meta.SetConverter(converter)
		}
	}

	if comparer, ok := dataType.GetMetadataValue(parsing.MetadataComparer); ok {
		if hasComparer, ok := self.Context.Data.(HasComparer); ok {
			hasComparer.SetComparer(comparer)
		} else {
			panic(fmt.Sprintf("type has %q metadata but does not support HasComparer interface: %T", parsing.MetadataComparer, self.Context.Data))
		}
	}
}

func (self *Value) RenderProperty(dataType *DataType, dataDefinition *PropertyDefinition) {
	self.Render(dataType, dataDefinition, false, false)
}

func (self *Value) Normalize() normal.Value {
	return self.normalize(false)
}

func (self *Value) normalize(bare bool) normal.Value {
	var normalValue normal.Value

	switch data := self.Context.Data.(type) {
	case ard.Map:
		// This is for complex types (the "map" type is a ValueMap, below)
		normalMap := normal.NewMap()
		for key, value := range data {
			if value_, ok := value.(*Value); ok {
				normalMap.Put(key, value_.normalize(false))
			} else {
				normalMap.Put(key, normal.NewPrimitive(value))
			}
		}
		normalValue = normalMap

	case *ValueList:
		normalValue = data.Normalize(self.Context)

	case *ValueMap:
		normalValue = data.Normalize(self.Context)

	case *parsing.FunctionCall:
		NormalizeFunctionCallArguments(data, self.Context)
		normalValue = normal.NewFunctionCall(data)

	default:
		normalValue = normal.NewPrimitive(data)
	}

	if !bare {
		normalValue.SetMeta(self.Meta)
	}

	return normalValue
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

func (self Values) RenderAttributes(definitions AttributeDefinitions, context *parsing.Context) {
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

func (self Values) RenderProperties(definitions PropertyDefinitions, context *parsing.Context) {
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

func (self Values) RenderInputs(definitions ParameterDefinitions, context *parsing.Context) {
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

func (self Values) Normalize(normalConstrainables normal.Values) {
	for key, value := range self {
		normalConstrainables[key] = value.Normalize()
	}
}

// Utils

func NewValueMeta(context *parsing.Context, dataType *DataType, dataDefinition DataDefinition, validationClause *ValidationClause) *normal.ValueMeta {
	if dataType == nil {
		return nil
	}

	meta := dataType.NewValueMeta()

	if dataDefinition != nil {
		meta.Description = dataDefinition.GetDescription()
		meta.Metadata = dataDefinition.GetTypeMetadata()

		// Get ValidationClause from definition
		definitionValidation := dataDefinition.GetValidationClause()
		if definitionValidation != nil {
			if validationClause == nil {
				validationClause = definitionValidation
			} else {
				// Create a composite validation using $and
				compositeValidation := NewValidationClause(context)
				compositeValidation.Operator = "and"
				compositeValidation.Arguments = ard.List{validationClause, definitionValidation}
				validationClause = compositeValidation
			}
		}
	}

	// Get ValidationClause from dataType
	var dataTypeValidationClause *ValidationClause
	if dataType.ValidationClause != nil {
		dataTypeValidationClause = dataType.ValidationClause
		if validationClause == nil {
			validationClause = dataType.ValidationClause
		}
		// Note: If both property and datatype have validation clauses,
		// we'll add them separately to the meta instead of combining them
	}

	// Helper function to add a validation clause to meta
	addValidationToMeta := func(clause *ValidationClause) {
		if clause == nil {
			return
		}
		clause.DataType = dataType
		clause.Definition = dataDefinition

		// Add validation function to meta using the same logic as ValidationClauses.AddToMeta
		// Check if this is a map and should apply validation to values instead
		if meta.Type == "map" && meta.Value != nil {
			// Set the DataType to the element type for proper $value handling
			originalDataType := clause.DataType
			originalDefinition := clause.Definition
			if _, valueSchema := GetMapSchemas(dataType, dataDefinition); valueSchema != nil {
				clause.DataType = valueSchema.DataType
				clause.Definition = valueSchema
			}
			functionCall := clause.ToFunctionCall(context, true)
			// Restore original values
			clause.DataType = originalDataType
			clause.Definition = originalDefinition
			NormalizeFunctionCallArguments(functionCall, context)
			meta.Value.AddValidator(functionCall)
		} else if meta.Type == "list" && meta.Element != nil {
			// Set the DataType to the element type for proper $value handling
			originalDataType := clause.DataType
			originalDefinition := clause.Definition
			if entrySchema := GetListSchemas(dataType, dataDefinition); entrySchema != nil {
				clause.DataType = entrySchema.DataType
				clause.Definition = entrySchema
			}
			functionCall := clause.ToFunctionCall(context, true)
			// Restore original values
			clause.DataType = originalDataType
			clause.Definition = originalDefinition
			NormalizeFunctionCallArguments(functionCall, context)
			meta.Element.AddValidator(functionCall)
		} else {
			// Add validation function to meta
			functionCall := clause.ToFunctionCall(context, true)
			NormalizeFunctionCallArguments(functionCall, context)
			meta.AddValidator(functionCall)
		}
	}

	// Process data type validation clause first (if different from property validation clause)
	if dataTypeValidationClause != nil && dataTypeValidationClause != validationClause {
		addValidationToMeta(dataTypeValidationClause)
	}

	// Process property validation clause
	addValidationToMeta(validationClause)

	if internalTypeName, ok := dataType.GetInternalTypeName(); ok {
		switch internalTypeName {
		case ard.TypeList:
			if entrySchema := GetListSchemas(dataType, dataDefinition); entrySchema != nil {
				meta.Element = NewValueMeta(context, entrySchema.DataType, entrySchema, nil)
			}

		case ard.TypeMap:
			if keySchema, valueSchema := GetMapSchemas(dataType, dataDefinition); (keySchema != nil) && (valueSchema != nil) {
				meta.Key = NewValueMeta(context, keySchema.DataType, keySchema, nil)
				meta.Value = NewValueMeta(context, valueSchema.DataType, valueSchema, nil)
			}
		}
	}

	return meta
}
