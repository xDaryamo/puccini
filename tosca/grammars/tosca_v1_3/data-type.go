package tosca_v1_3

import (
	"github.com/tliron/puccini/ard"
	"github.com/tliron/puccini/tosca"
)

//
// DataType
//
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.7.6
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.6.6
//

type DataType struct {
	*Type `name:"data type"`

	PropertyDefinitions PropertyDefinitions `read:"properties,PropertyDefinition" inherit:"properties,Parent"`
	ConstraintClauses   ConstraintClauses   `read:"constraints,[]ConstraintClause" inherit:"constraints,Parent"`

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

// tosca.Inherits interface
func (self *DataType) Inherit() {
	log.Infof("{inherit} data type: %s", self.Name)

	if _, ok := self.GetInternalTypeName(); ok && (len(self.PropertyDefinitions) > 0) {
		// Doesn't make sense to be an internal type (non-complex) and also have properties (complex)
		self.Context.ReportPrimitiveType()
		self.PropertyDefinitions = make(PropertyDefinitions)
		return
	}

	if self.Parent == nil {
		return
	}

	self.PropertyDefinitions.Inherit(self.Parent.PropertyDefinitions)
}

func (self *DataType) GetInternalTypeName() (string, bool) {
	typeName, ok := self.GetMetadataValue("puccini-tosca.type")
	if !ok && (self.Parent != nil) {
		// The internal type metadata is inherited
		return self.Parent.GetInternalTypeName()
	}
	return typeName, ok
}

// Note that this may change the data (if it's a map), but that should be fine, because we intend to
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

		var d interface{}
		if d, ok = map_[key]; ok {
			childContext.Data = d
		} else if definition.Default != nil {
			// Assign default value
			d = definition.Default.Context.Data
			childContext.Data = d
			map_[key] = d
		}

		if ToFunction(childContext) {
			map_[key] = childContext.Data
		} else {
			definition.DataType.Complete(childContext)
		}
	}
}

//
// DataTypes
//

type DataTypes []*DataType
