package tosca_v2_0

import (
	"github.com/tliron/go-ard"
	"github.com/tliron/puccini/normal"
	"github.com/tliron/puccini/tosca/parsing"
)

//
// PropertyFilter
//
// [TOSCA-v2.0] @ ?
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.6.4
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.4
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.3
// [TOSCA-Simple-Profile-YAML-v1.0] @ 3.5.3
//

type PropertyFilter struct {
	*Entity `name:"property filter"`
	Name    string

	ValidationClauses ValidationClauses
}

func NewPropertyFilter(context *parsing.Context) *PropertyFilter {
	return &PropertyFilter{
		Entity: NewEntity(context),
		Name:   context.Name,
	}
}

// ([parsing.Reader] signature)
func ReadPropertyFilter(context *parsing.Context) parsing.EntityPtr {
	self := NewPropertyFilter(context)

	if context.Is(ard.TypeList) {
		context.ReadListItems(ReadValidationClause, func(item ard.Value) {
			self.ValidationClauses = append(self.ValidationClauses, item.(*ValidationClause))
		})
	} else {
		self.ValidationClauses = ValidationClauses{ReadValidationClause(context).(*ValidationClause)}
	}

	return self
}

func (self *PropertyFilter) Normalize(normalFunctionCallMap normal.FunctionCallMap) normal.FunctionCalls {
	if len(self.ValidationClauses) == 0 {
		return nil
	}

	normalFunctionCalls := self.ValidationClauses.Normalize(self.Context)
	if existing, ok := normalFunctionCallMap[self.Name]; ok {
		normalFunctionCallMap[self.Name] = append(existing, normalFunctionCalls...)
	} else {
		normalFunctionCallMap[self.Name] = normalFunctionCalls
	}
	return normalFunctionCalls
}

//
// PropertyFilters
//

type PropertyFilters []*PropertyFilter

func (self PropertyFilters) Normalize(normalFunctionCallMap normal.FunctionCallMap) {
	for _, propertyFilter := range self {
		propertyFilter.Normalize(normalFunctionCallMap)
	}
}
