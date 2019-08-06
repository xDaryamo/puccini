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
		Entity:           NewEntity(context),
		Properties:       make(Values),
		SourceInterfaces: make(InterfaceAssignments),
		TargetInterfaces: make(InterfaceAssignments),
	}
}

// tosca.Reader signature
func ReadRelationshipAssignment(context *tosca.Context) interface{} {
	self := NewRelationshipAssignment(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))
	return self
}

// tosca.Renderable interface
func (self *RelationshipAssignment) Render() {
	log.Info("{render} relationship")

	if self.RelationshipType != nil {
		self.Properties.RenderProperties(self.RelationshipType.PropertyDefinitions, "property", self.Context.FieldChild("properties", nil))
		self.SourceInterfaces.Render(self.RelationshipType.SourceInterfaceDefinitions, self.Context.FieldChild("source_interfaces", nil))
		self.TargetInterfaces.Render(self.RelationshipType.TargetInterfaceDefinitions, self.Context.FieldChild("target_interfaces", nil))
	}

	// TODO: this should apply only to derivatives of cloudify.relationships.connected_to and cloudify.relationships.depends_on
	for key, value := range self.Properties {
		switch key {
		case "connection_type":
			if connectionType := value.Context.ReadString(); connectionType != nil {
				switch *connectionType {
				case "all_to_all", "all_to_one":
				default:
					value.Context.ReportFieldUnsupportedValue()
				}
			}
		default:
			value.Context.ReportFieldUnsupported()
		}
	}

	self.Properties.SetIfNil(self.Context.FieldChild("properties", nil), "connection_type", "all_to_all")
}

func (self *RelationshipAssignment) Normalize(nodeTemplate *NodeTemplate, s *normal.ServiceTemplate, n *normal.NodeTemplate) *normal.Requirement {
	r := n.NewRequirement("relationship", self.Context.Path.String())

	if self.TargetNodeTemplate != nil {
		r.NodeTemplate = s.NodeTemplates[self.TargetNodeTemplate.Name]
	}
	r.CapabilityTypeName = &capabilityTypeName

	rr := r.NewRelationship()

	if types, ok := normal.GetTypes(self.Context.Hierarchy, self.RelationshipType); ok {
		rr.Types = types
	}

	self.Properties.Normalize(rr.Properties, "")
	self.SourceInterfaces.NormalizeForRelationshipSource(self, rr)
	self.TargetInterfaces.NormalizeForRelationshipTarget(self, rr)

	return r
}

//
// RelationshipAssignments
//

type RelationshipAssignments []*RelationshipAssignment
