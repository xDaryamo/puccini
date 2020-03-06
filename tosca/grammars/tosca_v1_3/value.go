package tosca_v1_3

import (
	"fmt"
	"strconv"
	"time"

	"github.com/tliron/puccini/ard"
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
	"github.com/tliron/yamlkeys"
)

//
// Value
//
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
	Type              string

	rendered bool
}

func NewValue(context *tosca.Context) *Value {
	return &Value{
		Entity: NewEntity(context),
		Name:   context.Name,
	}
}

// tosca.Reader signature
func ReadValue(context *tosca.Context) interface{} {
	ToFunctionCall(context)
	return NewValue(context)
}

// tosca.Reader signature
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.12.2.2
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.11.2.2
func ReadAttributeValue(context *tosca.Context) interface{} {
	self := NewValue(context)

	// Unpack long notation
	if context.Is("!!map") {
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

func (self *Value) RenderDataType(dataTypeName string) {
	if e, ok := self.Context.Namespace.Lookup(dataTypeName); ok {
		if dataType, ok := e.(*DataType); ok {
			self.RenderAttribute(dataType, nil, false, false)
		} else {
			self.Context.ReportUnknownDataType(dataTypeName)
		}
	} else {
		self.Context.ReportUnknownDataType(dataTypeName)
	}
}

func (self *Value) RenderAttribute(dataType *DataType, definition *AttributeDefinition, bare bool, allowNil bool) {
	if self.rendered {
		// Avoid rendering more than once (can happen if we were copied from definition "default")
		return
	}
	self.rendered = true

	if !bare {
		if self.Description == nil {
			if (definition != nil) && (definition.Description != nil) {
				self.Description = definition.Description
			} else {
				self.Description = dataType.Description
			}
		}
	}

	self.Type = dataType.Name

	dataType.Complete(self.Context)
	if !bare {
		dataType.ConstraintClauses.RenderAndAppend(&self.ConstraintClauses, dataType)
	}

	if _, ok := self.Context.Data.(*tosca.FunctionCall); ok {
		return
	}

	if allowNil && (self.Context.Data == nil) {
		return
	}

	// Internal types
	if internalTypeName, ok := dataType.GetInternalTypeName(); ok {
		if typeValidator, ok := ard.TypeValidators[internalTypeName]; ok {
			if self.Context.Data == nil {
				// Nil data only happens when an attribute is added despite not having a
				// "default" value; we will give it a valid zero value instead
				self.Context.Data = ard.TypeZeroes[internalTypeName]
			}

			if (internalTypeName == "!!str") && self.Context.HasQuirk(tosca.QuirkDataTypesStringPermissive) {
				switch data := self.Context.Data.(type) {
				case bool:
					self.Context.Data = strconv.FormatBool(data)

				case int:
					self.Context.Data = strconv.FormatInt(int64(data), 10)

				case int64:
					self.Context.Data = strconv.FormatInt(data, 10)

				case float64: // YAML parser returns float64
					self.Context.Data = strconv.FormatFloat(data, 'g', -1, 64)

				case time.Time:
					self.Context.Data = data.String()
				}
			}

			// Primitive types
			if typeValidator(self.Context.Data) {
				// Render list and map elements according to entry schema
				// (The entry schema may also have additional constraints)
				switch internalTypeName {
				case "!!seq", "!!map":
					if (definition == nil) || (definition.EntrySchema == nil) || (definition.EntrySchema.DataType == nil) {
						// This problem is reported in AttributeDefinition.Render
						return
					}

					entryDataType := definition.EntrySchema.DataType
					entryConstraints := definition.EntrySchema.RenderConstraints()

					if internalTypeName == "!!seq" {
						slice := self.Context.Data.(ard.List)

						valueList := NewValueList(definition, len(slice), self.Description, entryConstraints)

						for index, data := range slice {
							value := ReadAndRenderBareAttribute(self.Context.ListChild(index, data), entryDataType)
							valueList.Set(index, value)
						}

						self.Context.Data = valueList
					} else { // "!!map"
						if (definition == nil) || (definition.KeySchema == nil) {
							// TODO: check if this is reported elsewhere
							return
						}

						keyDataType := definition.KeySchema.DataType
						keyConstraints := definition.KeySchema.RenderConstraints()

						valueMap := NewValueMap(definition, self.Description, keyConstraints, entryConstraints)

						for key, data := range self.Context.Data.(ard.Map) {
							// Complex keys are stringified for the purpose of the contexts

							// Validate key schema
							keyContext := self.Context.MapChild(key, yamlkeys.KeyData(key))
							key = ReadAndRenderBareAttribute(keyContext, keyDataType)

							context := self.Context.MapChild(key, data)
							value := ReadAndRenderBareAttribute(context, entryDataType)
							value.ConstraintClauses = ConstraintClauses{}
							valueMap.Put(key, value)
						}

						self.Context.Data = valueMap
					}
				}
			} else {
				self.Context.ReportValueWrongType(tosca.GetCanonicalName(dataType))
			}
		} else {
			// Special types
			if read, ok := self.Context.Grammar.Readers[internalTypeName]; ok {
				self.Context.Data = read(self.Context)
			} else {
				// Avoid reporting more than once
				if !dataType.typeProblemReported {
					dataType.Context.ReportUnsupportedType()
					dataType.typeProblemReported = true
				}
			}
		}

		if comparer, ok := dataType.GetMetadataValue("puccini.comparer"); ok {
			if hasComparer, ok := self.Context.Data.(HasComparer); ok {
				hasComparer.SetComparer(comparer)
			} else {
				panic(fmt.Sprintf("type has \"puccini.comparer\" metadata but does not support HasComparer interface: %T", self.Context.Data))
			}
		}
	} else if self.Context.ValidateType("!!map") {
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
			} else {
				// PropertyDefinition.Required defaults to true
				required := (definition.Required == nil) || *definition.Required
				if required {
					self.Context.MapChild(key, data).ReportPropertyRequired("property")
				}
			}
		}
	}
}

func (self *Value) RenderProperty(dataType *DataType, definition *PropertyDefinition) {
	if definition == nil {
		self.RenderAttribute(dataType, nil, false, false)
	} else {
		self.RenderAttribute(dataType, definition.AttributeDefinition, false, false)
		definition.ConstraintClauses.RenderAndAppend(&self.ConstraintClauses, dataType)
	}
}

func ReadAndRenderBareAttribute(context *tosca.Context, dataType *DataType) *Value {
	self := ReadValue(context).(*Value)
	self.RenderAttribute(dataType, nil, true, false)
	return self
}

func (self *Value) Normalize() normal.Constrainable {
	var normalConstrainable normal.Constrainable

	switch data := self.Context.Data.(type) {
	case ard.Map:
		// This is for complex types (the "map" type is a ValueMap, below)
		normalMap := normal.NewMap()
		for key, value := range data {
			if v, ok := value.(*Value); ok {
				normalMap.Put(key, v.Normalize())
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
		value.Type = self.Type
		normalConstrainable = value
	}

	self.ConstraintClauses.NormalizeConstrainable(self.Context, normalConstrainable)

	if self.Description != nil {
		normalConstrainable.SetDescription(*self.Description)
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

func (self Values) RenderMissingValue(definition *AttributeDefinition, kind string, required bool, context *tosca.Context) {
	if definition.Default != nil {
		self[definition.Name] = definition.Default
	} else if required {
		context.MapChild(definition.Name, nil).ReportPropertyRequired(kind)
	} else if kind == "attribute" {
		// Attributes should always appear, even if they have no default value
		self[definition.Name] = NewValue(context.MapChild(definition.Name, nil))
	}
}

func (self Values) RenderProperties(definitions PropertyDefinitions, kind string, context *tosca.Context) {
	for key, definition := range definitions {
		if value, ok := self[key]; !ok {
			// PropertyDefinition.Required defaults to true
			required := (definition.Required == nil) || *definition.Required
			self.RenderMissingValue(definition.AttributeDefinition, kind, required, context)
			// (If the above assigns the "default" value -- it has already been rendered elsewhere)
		} else if definition.DataType != nil {
			value.RenderProperty(definition.DataType, definition)
		}
	}

	for key, value := range self {
		if _, ok := definitions[key]; !ok {
			value.Context.ReportUndeclared(kind)
			delete(self, key)
		}
	}
}

func (self Values) RenderAttributes(definitions AttributeDefinitions, context *tosca.Context) {
	for key, definition := range definitions {
		if _, ok := self[key]; !ok {
			self.RenderMissingValue(definition, "attribute", false, context)
		}
	}

	for key, value := range self {
		if definition, ok := definitions[key]; !ok {
			value.Context.ReportUndeclared("attribute")
			delete(self, key)
		} else if definition.DataType != nil {
			value.RenderAttribute(definition.DataType, definition, false, true)
		}
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
	Description      *string
	EntryDescription *string
	EntryConstraints ConstraintClauses
	Slice            []interface{}
}

func NewValueList(definition *AttributeDefinition, length int, description *string, entryConstraints ConstraintClauses) *ValueList {
	return &ValueList{
		Description:      description,
		EntryDescription: definition.EntrySchema.Description,
		EntryConstraints: entryConstraints,
		Slice:            make([]interface{}, length),
	}
}

func (self *ValueList) Set(index int, value interface{}) {
	self.Slice[index] = value
}

func (self *ValueList) Normalize(context *tosca.Context) *normal.List {
	normalList := normal.NewList(len(self.Slice))

	if self.Description != nil {
		normalList.Description = *self.Description
	}
	if self.EntryDescription != nil {
		normalList.EntryDescription = *self.EntryDescription
	}

	self.EntryConstraints.NormalizeListEntries(context, normalList)

	for index, value := range self.Slice {
		if v, ok := value.(*Value); ok {
			normalList.Set(index, v.Normalize())
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
	Description      *string
	KeyDescription   *string
	EntryDescription *string
	KeyConstraints   ConstraintClauses
	EntryConstraints ConstraintClauses
	Map              ard.Map
}

func NewValueMap(definition *AttributeDefinition, description *string, keyConstraints ConstraintClauses, entryConstraints ConstraintClauses) *ValueMap {
	var keyDescription *string
	if definition.KeySchema != nil {
		keyDescription = definition.KeySchema.Description
	}
	return &ValueMap{
		Description:      description,
		KeyDescription:   keyDescription,
		EntryDescription: definition.EntrySchema.Description,
		KeyConstraints:   keyConstraints,
		EntryConstraints: entryConstraints,
		Map:              make(ard.Map),
	}
}

func (self *ValueMap) Put(key interface{}, value interface{}) {
	self.Map[key] = value
}

func (self *ValueMap) Normalize(context *tosca.Context) *normal.Map {
	normalMap := normal.NewMap()

	if self.Description != nil {
		normalMap.Description = *self.Description
	}
	if self.KeyDescription != nil {
		normalMap.KeyDescription = *self.KeyDescription
	}
	if self.EntryDescription != nil {
		normalMap.ValueDescription = *self.EntryDescription
	}

	self.KeyConstraints.NormalizeMapKeys(context, normalMap)
	self.EntryConstraints.NormalizeMapValues(context, normalMap)

	for key, value := range self.Map {
		if k, ok := key.(*Value); ok {
			key = k.Normalize()
		}
		if v, ok := value.(*Value); ok {
			normalMap.Put(key, v.Normalize())
		} else {
			normalMap.Put(key, normal.NewValue(value))
		}
	}

	return normalMap
}
