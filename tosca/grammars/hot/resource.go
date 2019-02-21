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
// [https://docs.openstack.org/heat/rocky/template_guide/hot_spec.html#resources-section]
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

	DependsOnResources []*Resource `lookup:"depends_on,DependsOn" json:"-" yaml:"-"`
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
	context.ValidateUnsupportedFields(append(context.ReadFields(self, Readers)))

	if childContext, ok := context.GetFieldChild("depends_on"); ok {
		self.DependsOn = childContext.ReadStringOrStringList()
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

	if self.Type != nil {
		n.Types = normal.NewTypes(*self.Type)
	}

	self.Properties.Normalize(n.Properties)

	n.NewCapability("resource").Types = capabilityTypes

	return n
}

func (self *Resource) NormalizeDependencies(s *normal.ServiceTemplate) {
	n := s.NodeTemplates[self.Name]
	path := self.Context.FieldChild("depends_on", nil).Path

	for _, resource := range self.DependsOnResources {
		r := n.NewRequirement("dependency", path)
		r.NodeTemplate = s.NodeTemplates[resource.Name]
		r.CapabilityTypeName = &capabilityTypeName

		rr := r.NewRelationship()
		rr.Types = relationshipTypes
	}
}
