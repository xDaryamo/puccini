package cloudify_v1_3

import (
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
)

//
// RelationshipAssignment
//

type RelationshipAssignment struct {
	*Entity `name:"relationship assignment"`

	RelationshipTypeName   *string              `read:"type" require:"type"`
	TargetNodeTemplateName *string              `read:"target" require:"target"`
	Properties             Values               `read:"properties,Value"`
	SourceInterfaces       InterfaceAssignments `read:"source_interfaces,InterfaceAssignment"`
	TargetInterfaces       InterfaceAssignments `read:"target_interfaces,InterfaceAssignment"`

	RelationshipType   *RelationshipType `lookup:"type,RelationshipTypeName" json:"-" yaml:"-"`
	TargetNodeTemplate *NodeTemplate     `lookup:"target,TargetNodeTemplateName" json:"-" yaml:"-"`
}

func NewRelationshipAssignment(context *tosca.Context) *RelationshipAssignment {
	return &RelationshipAssignment{
		Entity:     NewEntity(context),
		Properties: make(Values),
	}
}

// tosca.Reader signature
func ReadRelationshipAssignment(context *tosca.Context) interface{} {
	self := NewRelationshipAssignment(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))
	ValidateRelationshipProperties(context, self.Properties)
	return self
}

func ValidateRelationshipProperties(context *tosca.Context, properties Values) {
	propertiesContext := context.FieldChild("properties", nil)
	for key, value := range properties {
		childContext := propertiesContext.MapChild(key, value.Data)
		switch key {
		case "connection_type":
			if connectionType := childContext.ReadString(); connectionType != nil {
				switch *connectionType {
				case "all_to_all", "all_to_one":
				default:
					childContext.ReportFieldUnsupportedValue()
				}
			}
		default:
			childContext.ReportFieldUnsupported()
		}
	}

	properties.SetIfNil(propertiesContext, "connection_type", "all_to_all")
}

func (self *RelationshipAssignment) Normalize(nodeTemplate *NodeTemplate, s *normal.ServiceTemplate, n *normal.NodeTemplate) *normal.Requirement {
	r := n.NewRequirement("relationship", self.Context.Path)

	if self.TargetNodeTemplate != nil {
		r.NodeTemplate = s.NodeTemplates[self.TargetNodeTemplate.Name]
	}
	r.CapabilityTypeName = &capabilityTypeName

	rr := r.NewRelationship()

	if types, ok := normal.GetTypes(self.Context.Hierarchy, self.RelationshipType); ok {
		rr.Types = types
	}

	self.Properties.Normalize(rr.Properties)

	// for key, intr := range self.SourceInterfaces {
	// 	if definition, ok := intr.GetDefinitionForRelationship(self); ok {
	// 		i := rr.NewInterface(key)
	// 		intr.Normalize(i, definition)
	// 	}
	// }

	return r
}

//
// RelationshipAssignments
//

type RelationshipAssignments []*RelationshipAssignment
