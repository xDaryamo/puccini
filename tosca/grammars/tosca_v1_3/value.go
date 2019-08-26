package tosca_v1_3

import (
	"strconv"
	"time"

	"github.com/tliron/puccini/ard"
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
)

//
// Value
//
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.10
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.12
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.9
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.11
//

type Value struct {
	*Entity `name:"value"`
	Name    string

	ConstraintClauses ConstraintClauses
	Description       *string

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
	if context.Is("map") {
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
			self.RenderAttribute(dataType, nil, false)
		} else {
			self.Context.ReportUnknownDataType(dataTypeName)
		}
	} else {
		self.Context.ReportUnknownDataType(dataTypeName)
	}
}

func (self *Value) RenderProperty(dataType *DataType, definition *PropertyDefinition) {
	if definition == nil {
		self.RenderAttribute(dataType, nil, false)
	} else {
		self.RenderAttribute(dataType, definition.AttributeDefinition, false)
		definition.ConstraintClauses.Render(&self.ConstraintClauses, dataType)
	}
}

func (self *Value) RenderAttribute(dataType *DataType, definition *AttributeDefinition, allowNil bool) {
	if self.rendered {
		// Avoid rendering more than once (can happen if we use the "default" value)
		return
	}
	self.rendered = true

	if self.Description == nil {
		if (definition != nil) && (definition.Description != nil) {
			self.Description = definition.Description
		} else {
			self.Description = dataType.Description
		}
	}

	if allowNil && (self.Context.Data == nil) {
		return
	}

	if _, ok := self.Context.Data.(*tosca.FunctionCall); ok {
		return
	}

	dataType.Complete(self.Context)
	dataType.ConstraintClauses.Render(&self.ConstraintClauses, dataType)

	// Internal types
	if typeName, ok := dataType.GetInternalTypeName(); ok {
		if validator, ok := tosca.PrimitiveTypeValidators[typeName]; ok {
			if self.Context.Data == nil {
				// Nil data only happens when an attribute is added despite not having a
				// "default" value; we will give it a valid zero value instead
				self.Context.Data = tosca.PrimitiveTypeZeroes[typeName]
			}

			if (typeName == "string") && self.Context.HasQuirk("data_types.string.permissive") {
				switch self.Context.Data.(type) {
				case bool:
					self.Context.Data = strconv.FormatBool(self.Context.Data.(bool))
				case int:
					self.Context.Data = strconv.FormatInt(int64(self.Context.Data.(int)), 10)
				case float64: // YAML parser returns float64
					self.Context.Data = strconv.FormatFloat(self.Context.Data.(float64), 'f', -1, 64)
				case time.Time:
					self.Context.Data = self.Context.Data.(time.Time).String()
				}
			}

			// Primitive types
			if validator(self.Context.Data) {
				// Render list and map elements according to entry schema
				// (The entry schema may also have additional constraints)
				switch typeName {
				case "list", "map":
					if (definition == nil) || (definition.EntrySchema == nil) || (definition.EntrySchema.DataType == nil) {
						// This problem is reported in AttributeDefinition.Render
						return
					}

					entryDataType := definition.EntrySchema.DataType
					constraints := definition.EntrySchema.ConstraintClauses
					description := definition.EntrySchema.Description
					switch typeName {
					case "list":
						slice := self.Context.Data.(ard.List)
						for index, data := range slice {
							value := ReadValue(self.Context.ListChild(index, data)).(*Value)
							value.RenderAttribute(entryDataType, nil, false)
							constraints.Render(&value.ConstraintClauses, entryDataType)
							if description != nil {
								value.Description = description
							}
							slice[index] = value
						}
					case "map":
						map_ := self.Context.Data.(ard.Map)
						for key, data := range map_ {
							value := ReadValue(self.Context.MapChild(key, data)).(*Value)
							value.RenderAttribute(entryDataType, nil, false)
							constraints.Render(&value.ConstraintClauses, entryDataType)
							if description != nil {
								value.Description = description
							}
							map_[key] = value
						}
					}
				}
			} else {
				self.Context.ReportValueWrongType(dataType.Name)
			}
		} else {
			// Special types
			if read, ok := self.Context.Grammar[typeName]; ok {
				self.Context.Data = read(self.Context)
			} else {
				// Avoid reporting more than once
				if !dataType.typeProblemReported {
					dataType.Context.ReportUnsupportedType()
					dataType.typeProblemReported = true
				}
			}
		}

		return
	}

	// Complex data types

	if !self.Context.ValidateType("map") {
		return
	}

	map_ := self.Context.Data.(ard.Map)

	// All properties must be defined in type
	for key := range map_ {
		if _, ok := dataType.PropertyDefinitions[key]; !ok {
			self.Context.MapChild(key, nil).ReportUndefined("property")
			delete(map_, key)
		}
	}

	// Render properties
	for key, definition := range dataType.PropertyDefinitions {
		if data, ok := map_[key]; ok {
			var value *Value
			if value, ok = data.(*Value); !ok {
				// Convert to value
				value = NewValue(self.Context.MapChild(key, data))
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

func (self *Value) Normalize() normal.Constrainable {
	var constrainable normal.Constrainable

	if list, ok := self.Context.Data.(ard.List); ok {
		l := normal.NewConstrainableList(len(list))
		for index, value := range list {
			if v, ok := value.(*Value); ok {
				l.List[index] = v.Normalize()
			} else {
				l.List[index] = normal.NewValue(value)
			}
		}
		constrainable = l
	} else if map_, ok := self.Context.Data.(ard.Map); ok {
		m := normal.NewConstrainableMap()
		for key, value := range map_ {
			if v, ok := value.(*Value); ok {
				m.Map[key] = v.Normalize()
			} else {
				m.Map[key] = normal.NewValue(value)
			}
		}
		constrainable = m
	} else if functionCall, ok := self.Context.Data.(*tosca.FunctionCall); ok {
		NormalizeFunctionCallArguments(functionCall, self.Context)
		constrainable = normal.NewFunctionCall(functionCall)
	} else {
		constrainable = normal.NewValue(self.Context.Data)
	}

	self.ConstraintClauses.NormalizeConstrainable(self.Context, constrainable)

	if self.Description != nil {
		constrainable.SetDescription(*self.Description)
	}

	return constrainable
}

//
// Values
//

type Values map[string]*Value

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
			value.Context.ReportUndefined(kind)
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
			value.Context.ReportUndefined("attribute")
			delete(self, key)
		} else if definition.DataType != nil {
			value.RenderAttribute(definition.DataType, definition, true)
		}
	}
}

func (self Values) Normalize(c normal.Constrainables) {
	for key, value := range self {
		c[key] = value.Normalize()
	}
}
