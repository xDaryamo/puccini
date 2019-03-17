package cloudify_v1_3

import (
	"github.com/tliron/puccini/tosca"
)

//
// DataType
//
// [https://docs.cloudify.co/4.5.5/developer/blueprints/spec-data-types/]
//

type DataType struct {
	*Type `name:"data type"`

	Description         *string             `read:"description" inherit:"description,Parent"`
	PropertyDefinitions PropertyDefinitions `read:"properties,PropertyDefinition" inherit:"properties,Parent"`

	Parent *DataType `lookup:"derived_from,ParentName" json:"-" yaml:"-"`

	typeProblemReported bool
}

func NewDataType(context *tosca.Context) *DataType {
	return &DataType{
		Type:                NewType(context),
		PropertyDefinitions: make(PropertyDefinitions),
	}
}

// tosca.Reader signature
func ReadDataType(context *tosca.Context) interface{} {
	self := NewDataType(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))
	return self
}

// tosca.Hierarchical interface
func (self *DataType) GetParent() interface{} {
	return self.Parent
}

func (self *DataType) GetInternalTypeName() (string, bool) {
	switch self.Name {
	case "string", "integer", "float", "boolean":
		return self.Name, true
	}
	return "", false
}
