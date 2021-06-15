package tosca_v2_0

import (
	"reflect"

	"github.com/tliron/kutil/ard"
	"github.com/tliron/puccini/tosca"
)

//
// PropertyMapping
//
// Attaches to SubstitutionMappings
//
// [TOSCA-v2.0] @ ?
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.8.8
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.8.8
//

type PropertyMapping struct {
	*Entity `name:"property mapping"`
	Name    string

	InputName        *string
	NodeTemplateName *string // deprecated in TOSCA 1.3
	PropertyName     *string // deprecated in TOSCA 1.3

	InputDefinition *ParameterDefinition `traverse:"ignore" json:"-" yaml:"-"`
	NodeTemplate    *NodeTemplate        `traverse:"ignore" json:"-" yaml:"-"` // deprecated in TOSCA 1.3
	Property        *Value               `traverse:"ignore" json:"-" yaml:"-"` // deprecated in TOSCA 1.3
}

func NewPropertyMapping(context *tosca.Context) *PropertyMapping {
	return &PropertyMapping{
		Entity: NewEntity(context),
		Name:   context.Name,
	}
}

// tosca.Reader signature
func ReadPropertyMapping(context *tosca.Context) tosca.EntityPtr {
	self := NewPropertyMapping(context)

	var read bool
	if context.Is(ard.TypeList) {
		if strings := context.ReadStringList(); strings != nil {
			switch len(*strings) {
			case 1:
				self.InputName = &(*strings)[0]
				read = true

			case 2:
				// Deprecated in TOSCA 1.3
				self.NodeTemplateName = &(*strings)[0]
				self.PropertyName = &(*strings)[1]
				read = true
			}
		}
	}

	if !read {
		// Fallback to constant value (deprecated in TOSCA 1.3)
		self.Property = ReadValue(context).(*Value)

		// self.Context.ReportValueMalformed("property mapping", "must be list of 1 or 2 strings")
	}

	return self
}

// tosca.Mappable interface
func (self *PropertyMapping) GetKey() string {
	return self.Name
}

func (self *PropertyMapping) Render(inputDefinitions ParameterDefinitions) {
	logRender.Debug("property mapping")

	if self.InputName != nil {
		inputName := *self.InputName
		var ok bool
		if self.InputDefinition, ok = inputDefinitions[inputName]; !ok {
			self.Context.ListChild(0, inputName).ReportUnknown("input")
		}
	} else if (self.NodeTemplateName != nil) && (self.PropertyName != nil) {
		// Deprecated in TOSCA 1.3
		nodeTemplateName := *self.NodeTemplateName
		var nodeTemplateType *NodeTemplate
		if nodeTemplate, ok := self.Context.Namespace.LookupForType(nodeTemplateName, reflect.TypeOf(nodeTemplateType)); ok {
			self.NodeTemplate = nodeTemplate.(*NodeTemplate)
			self.NodeTemplate.Render()
		} else {
			self.Context.ListChild(0, nodeTemplateName).ReportUnknown("node template")
			return
		}

		name := *self.PropertyName
		var ok bool
		if self.Property, ok = self.NodeTemplate.Properties[name]; !ok {
			self.Context.ListChild(1, name).ReportReferenceNotFound("property", self.NodeTemplate)
		}
	}
}

//
// PropertyMappings
//

type PropertyMappings map[string]*PropertyMapping

func (self PropertyMappings) Render(inputDefinitions ParameterDefinitions) {
	for _, mapping := range self {
		mapping.Render(inputDefinitions)
	}
}
