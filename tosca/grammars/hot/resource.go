package hot

import (
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
)

var DeletionPolicies = map[string]string{
	"Delete":   "Delete",
	"Retain":   "Retain",
	"Snapshot": "Snapshot",
	"delete":   "Delete",
	"retain":   "Retain",
	"snapshot": "Snapshot",
}

func GetDeletionPolicy(policy string) (string, bool) {
	if p, ok := DeletionPolicies[policy]; ok {
		return p, true
	}
	return "", false
}

//
// Resource
//
// [https://docs.openstack.org/heat/stein/template_guide/hot_spec.html#resources-section]
//

type Resource struct {
	*Entity `name:"resource"`
	Name    string `namespace:""`

	Type           *string    `read:"type" require:"type"`
	Properties     Values     `read:"properties,Value"`
	Metadata       *Data      `read:"metadata,Data"`
	DependsOn      *[]string  `read:"depends_on"`
	UpdatePolicy   *Data      `read:"update_policy,Data"`
	DeletionPolicy *string    `read:"deletion_policy"`
	ExternalID     *string    `read:"external_id"`
	Condition      *Condition `read:"condition,Condition"`

	ToscaType          *string   `json:"-" yaml:"-"`
	DependsOnResources Resources `lookup:"depends_on,DependsOn" json:"-" yaml:"-"`
}

func NewResource(context *tosca.Context) *Resource {
	return &Resource{
		Entity:     NewEntity(context),
		Name:       context.Name,
		Properties: make(Values),
	}
}

// tosca.Reader signature
func ReadResource(context *tosca.Context) interface{} {
	self := NewResource(context)
	context.ValidateUnsupportedFields(append(context.ReadFields(self)))

	if childContext, ok := context.GetFieldChild("depends_on"); ok {
		self.DependsOn = childContext.ReadStringOrStringList()
	}

	if self.Type != nil {
		if type_, ok := ResourceTypes[*self.Type]; ok {
			self.ToscaType = &type_
		} else {
			context.FieldChild("type", *self.Type).ReportFieldUnsupportedValue()
		}
	}

	if self.DeletionPolicy != nil {
		if policy, ok := GetDeletionPolicy(*self.DeletionPolicy); ok {
			self.DeletionPolicy = &policy
		} else {
			context.FieldChild("deletion_policy", *self.DeletionPolicy).ReportFieldUnsupportedValue()
		}
	}

	return self
}

var capabilityTypeName = "openstack.Resource"
var capabilityTypes = normal.NewTypes(capabilityTypeName)
var relationshipTypes = normal.NewTypes("openstack.DependsOn")

func (self *Resource) Normalize(s *normal.ServiceTemplate) *normal.NodeTemplate {
	log.Infof("{normalize} resource: %s", self.Name)

	n := s.NewNodeTemplate(self.Name)

	if self.ToscaType != nil {
		n.Types = normal.NewTypes(*self.ToscaType)
	}

	self.Properties.Normalize(n.Properties)

	n.NewCapability("resource").Types = capabilityTypes

	return n
}

func (self *Resource) NormalizeDependencies(s *normal.ServiceTemplate) {
	log.Infof("{normalize} resource dependencies: %s", self.Name)

	n := s.NodeTemplates[self.Name]
	path := self.Context.FieldChild("depends_on", nil).Path.String()

	for _, resource := range self.DependsOnResources {
		r := n.NewRequirement("dependency", path)
		r.NodeTemplate = s.NodeTemplates[resource.Name]
		r.CapabilityTypeName = &capabilityTypeName

		rr := r.NewRelationship()
		rr.Types = relationshipTypes
	}
}

//
// Resources
//

type Resources []*Resource

func (self Resources) Normalize(s *normal.ServiceTemplate) {
	for _, resource := range self {
		s.NodeTemplates[resource.Name] = resource.Normalize(s)
	}

	// Dependencies must be normalized after resources
	// (because they may reference other resources)
	for _, resource := range self {
		resource.NormalizeDependencies(s)
	}
}
