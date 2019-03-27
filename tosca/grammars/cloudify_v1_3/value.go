package cloudify_v1_3

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

	Data        interface{}
	Description *string

	rendered bool
}

func NewValue(context *tosca.Context) *Value {
	return &Value{
		Entity: NewEntity(context),
		Name:   context.Name,
		Data:   context.Data,
	}
}

// tosca.Reader signature
func ReadValue(context *tosca.Context) interface{} {
	self := NewValue(context)
	self.Data = ToFunctions(context)
	return self
}

// tosca.Mappable interface
func (self *Value) GetKey() string {
	return self.Name
}

func (self *Value) RenderProperty(dataType *DataType, definition *PropertyDefinition) {
	if definition == nil {
		self.RenderParameter(dataType, nil, false)
	} else {
		self.RenderParameter(dataType, definition.ParameterDefinition, false)
	}
}

func (self *Value) RenderParameter(dataType *DataType, definition *ParameterDefinition, allowNil bool) {
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

	//dataType.Complete(self.Data, self.Context)

	// Internal types
	if typeName, ok := dataType.GetInternalTypeName(); ok {
		if validator, ok := tosca.PrimitiveTypeValidators[typeName]; ok {
			if self.Data == nil {
				// Nil data only happens when an parameter is added despite not having a
				// "default" value; we will give it a valid zero value instead
				self.Data = tosca.PrimitiveTypeZeroes[typeName]
			}

			// Primitive types
			if !validator(self.Data) {
				self.Context.ReportValueWrongType(dataType.Name)
			}
		} else {
			// Special types
			read, ok := Grammar[typeName]
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
			l.List[index] = NewValue(self.Context.ListChild(index, value)).Normalize()
		}
		constrainable = l
	} else if map_, ok := self.Data.(ard.Map); ok {
		m := normal.NewConstrainableMap()
		for key, value := range map_ {
			m.Map[key] = NewValue(self.Context.MapChild(key, value)).Normalize()
		}
		constrainable = m
	} else if function, ok := self.Data.(*tosca.Function); ok {
		NormalizeFunctionArguments(function, self.Context)
		constrainable = normal.NewFunction(function)
	} else {
		constrainable = normal.NewValue(self.Data)
	}

	if self.Description != nil {
		constrainable.SetDescription(*self.Description)
	}

	return constrainable
}

//
// Values
//

type Values map[string]*Value

func (self Values) SetIfNil(context *tosca.Context, key string, data interface{}) {
	if _, ok := self[key]; !ok {
		self[key] = NewValue(context.MapChild(key, data))
	}
}

func (self Values) RenderMissingValue(definition *ParameterDefinition, kind string, required bool, context *tosca.Context) {
	if definition.Default != nil {
		self[definition.Name] = definition.Default
	} else if required {
		context.MapChild(definition.Name, nil).ReportPropertyRequired(kind)
	} else if kind == "attribute" {
		// Parameters should always appear, even if they have no default value
		self[definition.Name] = NewValue(context.MapChild(definition.Name, nil))
	}
}

func (self Values) RenderProperties(definitions PropertyDefinitions, kind string, context *tosca.Context) {
	for key, definition := range definitions {
		if value, ok := self[key]; !ok {
			// PropertyDefinition.Required defaults to true
			required := (definition.Required == nil) || *definition.Required
			self.RenderMissingValue(definition.ParameterDefinition, kind, required, context)
			// (If the above assigns the "default" value -- it has already been rendered elsewhere)
		} else if definition.DataType != nil {
			value.RenderProperty(definition.DataType, definition)
		}
	}
}

func (self Values) RenderParameters(definitions ParameterDefinitions, kind string, context *tosca.Context) {
	for key, definition := range definitions {
		_, ok := self[key]
		if !ok {
			self.RenderMissingValue(definition, kind, false, context)
		}
	}

	for key, value := range self {
		definition, ok := definitions[key]
		if !ok {
			value.Context.ReportUndefined(kind)
			delete(self, key)
		} else if definition.DataType != nil {
			value.RenderParameter(definition.DataType, definition, true)
		}
	}
}

func (self Values) Normalize(c normal.Constrainables, prefix string) {
	for key, value := range self {
		c[prefix+key] = value.Normalize()
	}
}
