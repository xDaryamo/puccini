package hot

import (
	"github.com/tliron/puccini/ard"
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
)

var Types = []string{
	"string",
	"number",
	"json",
	"comma_delimited_list",
	"boolean",
}

func IsTypeValid(type_ string) bool {
	for _, t := range Types {
		if t == type_ {
			return true
		}
	}
	return false
}

func ValidateValue(context *tosca.Context, type_ string) bool {
	switch type_ {
	case "string":
		return context.ValidateType("string")
	case "number":
		return context.ValidateType("integer", "float")
	case "json":
		return context.ValidateType("map", "list")
	case "comma_delimited_list":
		if context.ValidateType("list") {
			for index, e := range context.Data.(ard.List) {
				if !context.ListChild(index, e).ValidateType("string") {
					return false
				}
			}
			return true
		} else {
			return false
		}
	case "boolean":
		return context.ValidateType("boolean")
	default:
		panic("unsupported parameter type")
	}
}

//
// Parameter
//
// [https://docs.openstack.org/heat/rocky/template_guide/hot_spec.html#parameters-section]
//

type Parameter struct {
	*Entity `name:"parameter"`
	Name    string `namespace:""`

	Type        *string     `read:"type" require:"type"`
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
func ReadParameter(context *tosca.Context) interface{} {
	self := NewParameter(context)
	context.ValidateUnsupportedFields(append(context.ReadFields(self, Readers)))

	if self.Type != nil {
		type_ := *self.Type
		if IsTypeValid(type_) {
			if self.Default != nil {
				self.Default.Fix(type_)
				ValidateValue(context.FieldChild("default", self.Default.Data), type_)
			}
		} else {
			context.FieldChild("type", type_).ReportFieldUnsupportedValue()
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
	log.Info("{render} parameter")

	if self.Type == nil {
		return
	}

	type_ := *self.Type

	if self.Value != nil {
		ValidateValue(self.Context.WithData(self.Value.Data), type_)
	}
}

func (self *Parameter) Normalize(context *tosca.Context) normal.Constrainable {
	value := self.Value
	if value == nil {
		if self.Default != nil {
			value = self.Default
		} else {
			// Parameters should always appear, even if they have no default value
			value = NewValue(context.MapChild(self.Name, nil))
		}
	}
	return value.Normalize()
}

//
// Parameters
//

type Parameters map[string]*Parameter

func (self Parameters) Normalize(c normal.Constrainables, context *tosca.Context) {
	for key, parameter := range self {
		c[key] = parameter.Normalize(context)
	}
}
