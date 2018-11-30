package v1_1

import (
	"github.com/tliron/puccini/tosca"
)

//
// PropertyFilter
//
// [TOSCA 1.1 3.5.3]
//

type PropertyFilter struct {
	*Entity `name:"property filter"`
	Name    string

	ConstraintClauses ConstraintClauses
}

func NewPropertyFilter(context *tosca.Context) *PropertyFilter {
	return &PropertyFilter{
		Entity:            NewEntity(context),
		Name:              context.Name,
		ConstraintClauses: make(ConstraintClauses, 0),
	}
}

// tosca.Reader signature
func ReadPropertyFilter(context *tosca.Context) interface{} {
	self := NewPropertyFilter(context)
	context.ReadListItems(ReadConstraintClause, func(item interface{}) {
		self.ConstraintClauses = append(self.ConstraintClauses, item.(*ConstraintClause))
	})
	return self
}

// tosca.Mappable interface
func (self *PropertyFilter) GetKey() string {
	return self.Name
}

// TODO: move to Clout
func (self *PropertyFilter) Apply(nodeTemplate *NodeTemplate) bool {
	if _, ok := nodeTemplate.Properties[self.Name]; ok {
		//var constrainable = value.Normalize()
		//self.ConstraintClauses.Normalize(value.Context, constrainable)
	}

	return true
}

//
// PropertyFilters
//

type PropertyFilters map[string]*PropertyFilter

// TODO: move to Clout
func (self PropertyFilters) Apply(nodeTemplate *NodeTemplate) bool {
	if len(self) == 0 {
		return true
	}

	for _, propertyFilter := range self {
		if !propertyFilter.Apply(nodeTemplate) {
			return false
		}
	}

	return true
}
