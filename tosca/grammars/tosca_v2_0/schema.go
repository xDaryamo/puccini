package tosca_v2_0

import (
	"sync"

	"github.com/tliron/kutil/ard"
	"github.com/tliron/puccini/tosca"
)

//
// Schema
//

type Schema struct {
	*Entity `name:"schema"`

	DataTypeName      *string           `read:"type" mandatory:""`
	Description       *string           `read:"description"`
	ConstraintClauses ConstraintClauses `read:"constraints,[]ConstraintClause" traverse:"ignore"`

	DataType *DataType `lookup:"type,DataTypeName" json:"-" yaml:"-"`

	renderOnce sync.Once
}

func NewSchema(context *tosca.Context) *Schema {
	return &Schema{Entity: NewEntity(context)}
}

// tosca.Reader signature
func ReadSchema(context *tosca.Context) tosca.EntityPtr {
	self := NewSchema(context)

	if context.Is(ard.TypeMap) {
		// Long notation
		context.ValidateUnsupportedFields(context.ReadFields(self))
	} else if context.ValidateType(ard.TypeMap, ard.TypeString) {
		// Short notation
		self.DataTypeName = context.FieldChild("type", context.Data).ReadString()
	}

	return self
}

// tosca.Renderable interface
// Avoid rendering more than once (can happen if we were called from Schema.GetConstraints)
func (self *Schema) Render() {
	self.renderOnce.Do(self.render)
}

func (self *Schema) render() {
	logRender.Debug("schema")

	if self.DataType == nil {
		return
	}
}

func (self *Schema) LookupDataType() bool {
	if self.DataTypeName != nil {
		dataTypeName := *self.DataTypeName
		var ok bool
		if self.DataType, ok = GetDataType(self.Context, dataTypeName); ok {
			return true
		} else {
			self.Context.ReportMissingEntrySchema(dataTypeName)
		}
	}

	return false
}

func (self *Schema) GetConstraints() ConstraintClauses {
	if self.DataType != nil {
		self.ConstraintClauses.Render(self.DataType, nil)
		self.DataType.ConstraintClauses.Render(self.DataType, nil)
		return self.DataType.ConstraintClauses.Append(self.ConstraintClauses)
	} else {
		return self.ConstraintClauses
	}
}
