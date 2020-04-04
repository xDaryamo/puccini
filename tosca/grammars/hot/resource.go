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
func ReadResource(context *tosca.Context) tosca.EntityPtr {
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

var capabilityTypeName = "Resource"
var capabilityTypes = normal.NewTypes(capabilityTypeName)
var relationshipTypes = normal.NewTypes("DependsOn")

func (self *Resource) Normalize(normalServiceTemplate *normal.ServiceTemplate) *normal.NodeTemplate {
	log.Infof("{normalize} resource: %s", self.Name)

	normalNodeTemplate := normalServiceTemplate.NewNodeTemplate(self.Name)

	if self.ToscaType != nil {
		normalNodeTemplate.Types = normal.NewTypes(*self.ToscaType)
	}

	self.Properties.Normalize(normalNodeTemplate.Properties)

	normalNodeTemplate.NewCapability("resource").Types = capabilityTypes

	return normalNodeTemplate
}

func (self *Resource) NormalizeDependencies(normalServiceTemplate *normal.ServiceTemplate) {
	log.Infof("{normalize} resource dependencies: %s", self.Name)

	normalNodeTemplate := normalServiceTemplate.NodeTemplates[self.Name]
	path := self.Context.FieldChild("depends_on", nil).Path.String()

	for _, resource := range self.DependsOnResources {
		normalRequirement := normalNodeTemplate.NewRequirement("dependency", path)
		normalRequirement.NodeTemplate = normalServiceTemplate.NodeTemplates[resource.Name]
		normalRequirement.CapabilityTypeName = &capabilityTypeName

		normalRelationship := normalRequirement.NewRelationship()
		normalRelationship.Types = relationshipTypes
	}
}

//
// Resources
//

type Resources []*Resource

func (self Resources) Normalize(normalServiceTemplate *normal.ServiceTemplate) {
	for _, resource := range self {
		normalServiceTemplate.NodeTemplates[resource.Name] = resource.Normalize(normalServiceTemplate)
	}

	// Dependencies must be normalized after resources
	// (because they may reference other resources)
	for _, resource := range self {
		resource.NormalizeDependencies(normalServiceTemplate)
	}
}
