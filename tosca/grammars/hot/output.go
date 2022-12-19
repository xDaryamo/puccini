package hot

import (
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
)

//
// Output
//
// [https://docs.openstack.org/heat/wallaby/template_guide/hot_spec.html#outputs-section]
//

type Output struct {
	*Entity `name:"output"`
	Name    string `namespace:""`

	Description *string    `read:"description"`
	Value       *Value     `read:"value,Value" mandatory:""`
	Condition   *Condition `read:"condition,Condition"`
}

func NewOutput(context *tosca.Context) *Output {
	return &Output{
		Entity: NewEntity(context),
		Name:   context.Name,
	}
}

// tosca.Reader signature
func ReadOutput(context *tosca.Context) tosca.EntityPtr {
	self := NewOutput(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))
	return self
}

// tosca.Mappable interface
func (self *Output) GetKey() string {
	return self.Name
}

func (self *Output) Normalize(context *tosca.Context) normal.Value {
	var value *Value
	if self.Value != nil {
		value = self.Value
	} else {
		// Parameters should always appear, even if they have no default value
		value = NewValue(context.MapChild(self.Name, nil))
	}
	return value.Normalize()
}

//
// Outputs
//

type Outputs map[string]*Output

func (self Outputs) Normalize(normalConstrainables normal.Values, context *tosca.Context) {
	for key, output := range self {
		normalConstrainables[key] = output.Value.Normalize()
	}
}
