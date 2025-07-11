package tosca_v1_3

import (
	"github.com/tliron/go-ard"
	"github.com/tliron/puccini/tosca/grammars/tosca_v2_0"
	"github.com/tliron/puccini/tosca/parsing"
)

//
// InterfaceMapping
//
// [TOSCA-Simple-Profile-YAML-v1.3] @ 3.8.12
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.8.11
//

type InterfaceMapping struct {
	*tosca_v2_0.Entity `name:"interface mapping"`
	Name               string

	// TOSCA 1.3 specific fields
	NodeTemplateName *string
	InterfaceName    *string

	// For compatibility, we'll have empty OperationMappings
	OperationMappings map[string]string `json:"operationMappings" yaml:"operationMappings"`

	// Flag to distinguish TOSCA 1.3 format from TOSCA 2.0
	IsTosca13Format bool `json:"-" yaml:"-"`
}

func NewInterfaceMapping(context *parsing.Context) *InterfaceMapping {
	return &InterfaceMapping{
		Entity:            tosca_v2_0.NewEntity(context),
		Name:              context.Name,
		OperationMappings: make(map[string]string),
	}
}

// ([parsing.Reader] signature)
func ReadInterfaceMapping(context *parsing.Context) parsing.EntityPtr {
	self := NewInterfaceMapping(context)

	// TOSCA 1.3 format: [node_template, interface_name]
	if context.ValidateType(ard.TypeList) {
		strings := context.ReadStringListFixed(2)
		if strings != nil {
			nodeTemplateName := (*strings)[0]
			interfaceName := (*strings)[1]

			self.NodeTemplateName = &nodeTemplateName
			self.InterfaceName = &interfaceName
			self.IsTosca13Format = true
		}
	}

	return self
}

// ([parsing.Mappable] interface)
func (self *InterfaceMapping) GetKey() string {
	return self.Name
}

// ([parsing.Renderable] interface)
func (self *InterfaceMapping) Render() {
	// TOSCA 1.3 interface mappings don't need workflow validation
	// They just map interfaces to node templates
}

// Convert to tosca_v2_0.InterfaceMapping for compatibility
func (self *InterfaceMapping) ToV2InterfaceMapping() *tosca_v2_0.InterfaceMapping {
	v2Mapping := tosca_v2_0.NewInterfaceMapping(self.Context)
	v2Mapping.Name = self.Name
	// Don't copy OperationMappings to avoid validation issues
	return v2Mapping
}
