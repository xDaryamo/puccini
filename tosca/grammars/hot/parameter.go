package hot

import (
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
)

//
// Parameter
//
// [https://docs.openstack.org/heat/wallaby/template_guide/hot_spec.html#parameters-section]
//

type Parameter struct {
	*Entity `name:"parameter"`
	Name    string `namespace:""`

	Type        *string     `read:"type" mandatory:""`
	Label       *string     `read:"label"`
	Description *string     `read:"description"`
	Default     *Value      `read:"default,Value"`
	Hidden      *bool       `read:"hidden"`
	Constraints Constraints `read:"constraints,[]Constraint"`
	Immutable   *bool       `read:"immutable"`
	Tags        *[]string   `read:"tags"`

	Value *Value
}

func NewParameter(context *tosca.Context) *Parameter {
	return &Parameter{
		Entity: NewEntity(context),
		Name:   context.Name,
	}
}

// tosca.Reader signature
func ReadParameter(context *tosca.Context) tosca.EntityPtr {
	self := NewParameter(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))

	if self.Type != nil {
		type_ := *self.Type
		if IsParameterTypeValid(type_) {
			if self.Default != nil {
				self.Default.CoerceParameterType(type_)
				self.Default.ValidateParameterType(type_)
				self.Default.Constraints = self.Constraints
			}
		} else {
			context.FieldChild("type", type_).ReportKeynameUnsupportedValue()
		}
	}

	return self
}

// tosca.Mappable interface
func (self *Parameter) GetKey() string {
	return self.Name
}

// tosca.Renderable interface
func (self *Parameter) Render() {
	self.renderOnce.Do(self.render)
}

func (self *Parameter) render() {
	logRender.Debugf("parameter: %s", self.Name)

	if self.Value != nil {
		if self.Type != nil {
			type_ := *self.Type
			if IsParameterTypeValid(type_) {
				self.Value.CoerceParameterType(type_)
				self.Value.ValidateParameterType(type_)
			}
			self.Value.Constraints = self.Constraints
		}
	} else if self.Default == nil {
		self.Context.ReportValueRequired("parameter")
	}
}

func (self *Parameter) Normalize(context *tosca.Context) normal.Value {
	value := self.Value
	if value == nil {
		if self.Default != nil {
			value = self.Default
		} else {
			// Parameters should always appear, even if they have no default value
			// (Note that in HOT they are *always* required, so it would be abnormal)
			value = NewValue(context.MapChild(self.Name, nil))
		}
	}
	// TODO: normalize Hidden, Mutable, etc.
	return value.Normalize()
}

//
// Parameters
//

type Parameters map[string]*Parameter

func (self Parameters) Normalize(c normal.Values, context *tosca.Context) {
	for key, parameter := range self {
		c[key] = parameter.Normalize(context)
	}
}
