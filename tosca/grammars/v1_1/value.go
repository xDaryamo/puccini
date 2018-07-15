package v1_1

import (
	"github.com/tliron/puccini/ard"
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
)

//
// Value
//

type Value struct {
	*Entity `name:"value"`
	Name    string

	Data              interface{}
	ConstraintClauses ConstraintClauses
	Description       *string

	rendered bool
}

func NewValue(context *tosca.Context) *Value {
	return &Value{
		Entity:            NewEntity(context),
		Name:              context.Name,
		Data:              context.Data,
		ConstraintClauses: make(ConstraintClauses, 0),
	}
}

// tosca.Reader signature
func ReadValue(context *tosca.Context) interface{} {
	if function, ok := GetFunction(context); ok {
		context.Data = function
	}
	return NewValue(context)
}

// tosca.Mappable interface
func (self *Value) GetKey() string {
	return self.Name
}

func (self *Value) RenderDataType(dataTypeName string) {
	e, ok := self.Context.Namespace.Lookup(dataTypeName)
	if !ok {
		self.Context.ReportUnknownDataType(dataTypeName)
		return
	}

	dataType, ok := e.(*DataType)
	if !ok {
		self.Context.ReportUnknownDataType(dataTypeName)
		return
	}

	self.RenderAttribute(dataType, nil, false)
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

	if (definition != nil) && (definition.Description != nil) {
		self.Description = definition.Description
	} else {
		self.Description = dataType.Description
	}

	if allowNil && (self.Data == nil) {
		return
	}

	if _, ok := self.Data.(*tosca.Function); ok {
		return
	}

	dataType.Complete(self.Data, self.Context)
	dataType.ConstraintClauses.Render(&self.ConstraintClauses, dataType)

	// Internal types
	if typeName, ok := dataType.GetInternalTypeName(); ok {
		if validator, ok := tosca.PrimitiveTypeValidators[typeName]; ok {
			if self.Data == nil {
				// Nil data only happens when an attribute is added despite not having a
				// "default" value; we will give it a valid zero value instead
				self.Data = tosca.PrimitiveTypeZeroes[typeName]
			}

			// Primitive types
			if validator(self.Data) {
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
						slice := self.Data.(ard.List)
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
						map_ := self.Data.(ard.Map)
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
			read, ok := Readers[typeName]
			if !ok {
				// Avoid reporting more than once
				if !dataType.typeProblemReported {
					dataType.Context.ReportUnsupportedType()
					dataType.typeProblemReported = true
				}
				return
			}

			self.Data = read(self.Context)
		}

		return
	}

	// Complex data types

	if !self.Context.ValidateType("map") {
		return
	}

	map_ := self.Data.(ard.Map)

	// All properties must be defined in type
	for key := range map_ {
		_, ok := dataType.PropertyDefinitions[key]
		if !ok {
			self.Context.MapChild(key, nil).ReportUndefined("property")
			delete(map_, key)
		}
	}

	// Render properties
	for key, definition := range dataType.PropertyDefinitions {
		data, ok := map_[key]
		if !ok {
			// PropertyDefinition.Required defaults to true
			required := (definition.Required == nil) || *definition.Required
			if required {
				self.Context.MapChild(key, data).ReportPropertyRequired("property")
			}
		} else {
			var value *Value
			if value, ok = data.(*Value); !ok {
				// Convert to value
				value = NewValue(self.Context.MapChild(key, data))
				map_[key] = value
			}
			if definition.DataType != nil {
				value.RenderProperty(definition.DataType, definition)
			}
		}
	}
}

func (self *Value) Normalize() normal.Constrainable {
	var constrainable normal.Constrainable

	if list, ok := self.Data.(ard.List); ok {
		l := normal.NewConstrainableList(len(list))
		for index, value := range list {
			if v, ok := value.(*Value); ok {
				l.List[index] = v.Normalize()
			} else {
				l.List[index] = normal.NewValue(value)
			}
		}
		constrainable = l
	} else if map_, ok := self.Data.(ard.Map); ok {
		m := normal.NewConstrainableMap()
		for key, value := range map_ {
			if v, ok := value.(*Value); ok {
				m.Map[key] = v.Normalize()
			} else {
				m.Map[key] = normal.NewValue(value)
			}
		}
		constrainable = m
	} else if function, ok := self.Data.(*tosca.Function); ok {
		NormalizeFunctionArguments(function, self.Context)
		constrainable = normal.NewFunction(function)
	} else {
		constrainable = normal.NewValue(self.Data)
	}

	self.ConstraintClauses.Normalize(self.Context, constrainable)

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
		value, ok := self[key]
		if !ok {
			// PropertyDefinition.Required defaults to true
			required := (definition.Required == nil) || *definition.Required
			self.RenderMissingValue(definition.AttributeDefinition, kind, required, context)
			// (If the above assigns the "default" value -- it has already been coerced elsewhere)
		} else if definition.DataType != nil {
			value.RenderProperty(definition.DataType, definition)
		}
	}

	for key, value := range self {
		_, ok := definitions[key]
		if !ok {
			value.Context.ReportUndefined(kind)
			delete(self, key)
		}
	}
}

func (self Values) RenderAttributes(definitions AttributeDefinitions, context *tosca.Context) {
	for key, definition := range definitions {
		_, ok := self[key]
		if !ok {
			self.RenderMissingValue(definition, "attribute", false, context)
		}
	}

	for key, value := range self {
		definition, ok := definitions[key]
		if !ok {
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
