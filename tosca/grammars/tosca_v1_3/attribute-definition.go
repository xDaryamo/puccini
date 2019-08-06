package tosca_v1_3

import (
	"github.com/tliron/puccini/tosca"
)

//
// AttributeDefinition
//
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.11
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.10
//

type AttributeDefinition struct {
	*Entity `name:"attribute definition"`
	Name    string

	Description  *string      `read:"description" inherit:"description,DataType"`
	DataTypeName *string      `read:"type"` // required only if cannot be inherited or discovered
	EntrySchema  *EntrySchema `read:"entry_schema,EntrySchema"`
	Default      *Value       `read:"default,Value"`
	Status       *string      `read:"status"`

	DataType *DataType `lookup:"type,DataTypeName" json:"-" yaml:"-"`

	typeMissingProblemReported bool
}

func NewAttributeDefinition(context *tosca.Context) *AttributeDefinition {
	return &AttributeDefinition{
		Entity: NewEntity(context),
		Name:   context.Name,
	}
}

// tosca.Reader signature
func ReadAttributeDefinition(context *tosca.Context) interface{} {
	self := NewAttributeDefinition(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))
	return self
}

// tosca.Mappable interface
func (self *AttributeDefinition) GetKey() string {
	return self.Name
}

func (self *AttributeDefinition) Inherit(parentDefinition *AttributeDefinition) {
	if parentDefinition != nil {
		if ((self.Description == nil) || ((self.DataType != nil) && (self.Description == self.DataType.Description))) && (parentDefinition.Description != nil) {
			self.Description = parentDefinition.Description
		}
		if (self.DataTypeName == nil) && (parentDefinition.DataTypeName != nil) {
			self.DataTypeName = parentDefinition.DataTypeName
		}
		if (self.EntrySchema == nil) && (parentDefinition.EntrySchema != nil) {
			self.EntrySchema = parentDefinition.EntrySchema
		}
		if (self.Default == nil) && (parentDefinition.Default != nil) {
			self.Default = parentDefinition.Default
		}
		if (self.Status == nil) && (parentDefinition.Status != nil) {
			self.Status = parentDefinition.Status
		}
		if (self.DataType == nil) && (parentDefinition.DataType != nil) {
			self.DataType = parentDefinition.DataType
		}

		// Validate type compatibility
		if (self.DataType != nil) && (parentDefinition.DataType != nil) && !self.Context.Hierarchy.IsCompatible(parentDefinition.DataType, self.DataType) {
			self.Context.ReportIncompatibleType(self.DataType.Name, parentDefinition.DataType.Name)
			return
		}
	}

	if self.DataTypeName == nil {
		// Avoid reporting more than once
		if !self.typeMissingProblemReported {
			self.Context.FieldChild("type", nil).ReportFieldMissing()
			self.typeMissingProblemReported = true
		}
		return
	}

	if self.DataType == nil {
		return
	}

	switch self.DataType.Name {
	case "list", "map":
		// Make sure we have an entry schema
		if (self.EntrySchema == nil) || (self.EntrySchema.DataType == nil) {
			self.Context.ReportMissingEntrySchema(self.DataType.Name)
		}
	}

	if self.Default != nil {
		// The "default" value must be a valid value of the type
		self.Default.RenderAttribute(self.DataType, self, false)
	}
}

//
// AttributeDefinitions
//

type AttributeDefinitions map[string]*AttributeDefinition

func (self AttributeDefinitions) Inherit(parentDefinitions AttributeDefinitions) {
	for name, definition := range self {
		if parentDefinitions != nil {
			if parentDefinition, ok := parentDefinitions[name]; ok {
				if definition != parentDefinition {
					definition.Inherit(parentDefinition)
				}
				continue
			}
		}

		definition.Inherit(nil)
	}
}
