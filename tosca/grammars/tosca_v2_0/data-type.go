package tosca_v2_0

import (
	"reflect"

	"github.com/tliron/kutil/ard"
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
)

type HasComparer interface {
	SetComparer(comparer string)
}

//
// DataType
//
// [TOSCA-v2.0] @ ?
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.7.6
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.7.6
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.6.6
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.6.5
//

type DataType struct {
	*Type `name:"data type"`

	PropertyDefinitions PropertyDefinitions `read:"properties,PropertyDefinition" inherit:"properties,Parent"`
	ConstraintClauses   ConstraintClauses   `read:"constraints,[]ConstraintClause"`
	KeySchema           *Schema             `read:"key_schema,Schema"`   // introduced in TOSCA 1.3
	EntrySchema         *Schema             `read:"entry_schema,Schema"` // introduced in TOSCA 1.3

	Parent *DataType `lookup:"derived_from,ParentName" json:"-" yaml:"-"`
}

func NewDataType(context *tosca.Context) *DataType {
	return &DataType{
		Type:                NewType(context),
		PropertyDefinitions: make(PropertyDefinitions),
	}
}

// tosca.Reader signature
func ReadDataType(context *tosca.Context) tosca.EntityPtr {
	self := NewDataType(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))
	return self
}

// tosca.Hierarchical interface
func (self *DataType) GetParent() tosca.EntityPtr {
	return self.Parent
}

// tosca.Inherits interface
func (self *DataType) Inherit() {
	logInherit.Debugf("data type: %s", self.Name)

	if _, ok := self.GetInternalTypeName(); ok && (len(self.PropertyDefinitions) > 0) {
		// Doesn't make sense to be an internal type (non-complex) and also have properties (complex)
		self.Context.ReportPrimitiveType()
		self.PropertyDefinitions = make(PropertyDefinitions)
		return
	}

	if self.Parent == nil {
		return
	}

	if (self.KeySchema == nil) && (self.Parent.KeySchema != nil) {
		self.KeySchema = self.Parent.KeySchema
	}
	if (self.EntrySchema == nil) && (self.Parent.EntrySchema != nil) {
		self.EntrySchema = self.Parent.EntrySchema
	}
	if self.Parent.ConstraintClauses != nil {
		self.ConstraintClauses = self.Parent.ConstraintClauses.Append(self.ConstraintClauses)
	}

	self.PropertyDefinitions.Inherit(self.Parent.PropertyDefinitions)
}

// parser.Renderable interface
func (self *DataType) Render() {
	logRender.Debugf("data type: %s", self.Name)

	self.ConstraintClauses.Render(self)

	if internalTypeName, ok := self.GetInternalTypeName(); ok {
		if _, ok := ard.TypeValidators[internalTypeName]; !ok {
			if _, ok := self.Context.Grammar.Readers[string(internalTypeName)]; !ok {
				self.Context.ReportUnsupportedType()
			}
		}
	}
}

func (self *DataType) GetInternalTypeName() (ard.TypeName, bool) {
	if typeName, ok := self.GetMetadataValue("puccini.type"); ok {
		return ard.TypeName(typeName), ok
	} else if self.Parent != nil {
		// The internal type metadata is inherited
		return self.Parent.GetInternalTypeName()
	} else {
		return ard.NoType, false
	}
}

func (self *DataType) GetInternal() (ard.TypeName, ard.TypeValidator, tosca.Reader, bool) {
	if internalTypeName, ok := self.GetInternalTypeName(); ok {
		if typeValidator, ok := ard.TypeValidators[internalTypeName]; ok {
			return internalTypeName, typeValidator, nil, true
		} else if reader, ok := self.Context.Grammar.Readers[string(internalTypeName)]; ok {
			return internalTypeName, nil, reader, true
		}
	}
	return ard.NoType, nil, nil, false
}

// Note that this may change the data (if it's a map), but that should be fine, because we intend
// for the data to be complete. For the same reason, this action is idempotent (subsequent calls to
// the same data will not have an effect).
func (self *DataType) Complete(context *tosca.Context) {
	map_, ok := context.Data.(ard.Map)
	if !ok {
		// Only for complex data types
		return
	}

	for key, definition := range self.PropertyDefinitions {
		childContext := context.MapChild(key, nil)

		var data ard.Value
		if data, ok = map_[key]; ok {
			childContext.Data = data
		} else if definition.Default != nil {
			// Assign default value
			data = definition.Default.Context.Data
			childContext.Data = data
			map_[key] = data
		}

		if ToFunctionCall(childContext) {
			map_[key] = childContext.Data
		} else if definition.DataType != nil {
			definition.DataType.Complete(childContext)
		}
	}
}

func (self *DataType) GetTypeInformation() *normal.TypeInformation {
	information := normal.NewTypeInformation()
	information.Name = tosca.GetCanonicalName(self)
	information.Metadata = tosca.GetInformationMetadata(self.Metadata)
	if self.Description != nil {
		information.Description = *self.Description
	}
	return information
}

func GetDataType(context *tosca.Context, name string) (*DataType, bool) {
	var dataType *DataType
	if entityPtr, ok := context.Namespace.LookupForType(name, reflect.TypeOf(dataType)); ok {
		return entityPtr.(*DataType), true
	} else {
		return nil, false
	}
}

//
// DataTypes
//

type DataTypes []*DataType
