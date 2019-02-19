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

func (self *Resource) Normalize(s *normal.ServiceTemplate) *normal.NodeTemplate {
	log.Infof("{normalize} resource: %s", self.Name)

	n := s.NewNodeTemplate(self.Name)

	if self.Type != nil {
		n.Types = normal.Types{*self.Type: normal.NewType(*self.Type)}
	}

	return n
}
