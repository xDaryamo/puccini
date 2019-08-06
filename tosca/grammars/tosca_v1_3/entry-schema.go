package tosca_v1_3

import (
	"github.com/tliron/puccini/tosca"
)

//
// EntrySchema
//

type EntrySchema struct {
	*Entity `name:"entry schema"`

	DataTypeName      *string           `read:"type" require:"type"`
	Description       *string           `read:"description" inherit:"description,DataType"`
	ConstraintClauses ConstraintClauses `read:"constraints,[]ConstraintClause"`

	DataType *DataType `lookup:"type,DataTypeName" json:"-" yaml:"-"`
}

func NewEntrySchema(context *tosca.Context) *EntrySchema {
	return &EntrySchema{Entity: NewEntity(context)}
}

// tosca.Reader signature
func ReadEntrySchema(context *tosca.Context) interface{} {
	self := NewEntrySchema(context)

	if context.Is("map") {
		// Long notation
		context.ValidateUnsupportedFields(context.ReadFields(self))
	} else if context.ValidateType("map", "string") {
		// Short notation
		self.DataTypeName = context.FieldChild("type", context.Data).ReadString()
	}

	return self
}
