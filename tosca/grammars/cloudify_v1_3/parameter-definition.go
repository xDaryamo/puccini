package cloudify_v1_3

import (
	"github.com/tliron/puccini/tosca"
)

//
// ParameterDefinition
//

type ParameterDefinition struct {
	*Entity `name:"parameter definition"`
	Name    string

	Description  *string `read:"description" inherit:"description,DataType"`
	DataTypeName *string `read:"type"`
	Default      *Value  `read:"default,Value"`

	DataType *DataType `lookup:"type,DataTypeName" json:"-" yaml:"-"`
}

func NewParameterDefinition(context *tosca.Context) *ParameterDefinition {
	return &ParameterDefinition{
		Entity: NewEntity(context),
		Name:   context.Name,
	}
}

// tosca.Reader signature
func ReadParameterDefinition(context *tosca.Context) interface{} {
	self := NewParameterDefinition(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))
	return self
}

// tosca.Mappable interface
func (self *ParameterDefinition) GetKey() string {
	return self.Name
}

//
// ParameterDefinitions
//

type ParameterDefinitions map[string]*ParameterDefinition
