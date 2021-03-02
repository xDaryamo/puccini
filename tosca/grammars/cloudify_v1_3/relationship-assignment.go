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

	RelationshipTypeName   *string              `read:"type" require:""`
	TargetNodeTemplateName *string              `read:"target" require:""`
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
func ReadRelationshipAssignment(context *tosca.Context) tosca.EntityPtr {
	self := NewRelationshipAssignment(context)
	context.ValidateUnsupportedFields(context.ReadFields(self))
	return self
}

// parser.Renderable interface
func (self *RelationshipAssignment) Render() {
	logRender.Debug("relationship")

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

func (self *RelationshipAssignment) Normalize(nodeTemplate *NodeTemplate, normalNodeTemplate *normal.NodeTemplate) *normal.Requirement {
	normalRequirement := normalNodeTemplate.NewRequirement("relationship", normal.NewLocationForContext(self.Context))

	if self.TargetNodeTemplate != nil {
		normalRequirement.NodeTemplate = normalNodeTemplate.ServiceTemplate.NodeTemplates[self.TargetNodeTemplate.Name]
	}
	normalRequirement.CapabilityTypeName = &capabilityTypeName

	normalRelationship := normalRequirement.NewRelationship()

	if types, ok := normal.GetTypes(self.Context.Hierarchy, self.RelationshipType); ok {
		normalRelationship.Types = types
	}

	self.Properties.Normalize(normalRelationship.Properties, "")
	self.SourceInterfaces.NormalizeForRelationshipSource(self, normalRelationship)
	self.TargetInterfaces.NormalizeForRelationshipTarget(self, normalRelationship)

	return normalRequirement
}

//
// RelationshipAssignments
//

type RelationshipAssignments []*RelationshipAssignment
