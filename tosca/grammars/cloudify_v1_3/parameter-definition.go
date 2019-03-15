package cloudify_v1_3

import (
	"github.com/tliron/puccini/tosca"
)

//
// ParameterDefinition
//

type ParameterDefinition struct {
	*Entity `name:"parameter definition"`

	Description  *string `read:"description" inherit:"description,DataType"`
	DataTypeName *string `read:"type"`
	Default      *Value  `read:"default,Value"`

	DataType *DataType `lookup:"type,DataTypeName" json:"-" yaml:"-"`
}

func NewParameterDefinition(context *tosca.Context) *ParameterDefinition {
	return &ParameterDefinition{
		Entity: NewEntity(context),
	}
}

// tosca.Reader signature
func ReadParameterDefinition(context *tosca.Context) interface{} {
	self := NewParameterDefinition(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))
	return self
}

// tosca.Renderable interface
// func (self *ParameterDefinition) Render() {
// 	log.Info("{render} parameter definition")

// 	if self.Value != nil {
// 		if self.Type != nil {
// 			type_ := *self.Type
// 			if IsParameterTypeValid(type_) {
// 				self.Value.CoerceParameterType(type_)
// 				self.Value.ValidateParameterType(type_)
// 			}
// 			self.Value.Constraints = self.Constraints
// 		}
// 	} else if self.Default == nil {
// 		self.Context.ReportPropertyRequired("parameter")
// 	}
// }

//
// ParameterDefinitions
//

type ParameterDefinitions map[string]*ParameterDefinition
