package cloudify_v1_3

import (
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
)

//
// Workflow
//
// [https://docs.cloudify.co/4.5.5/developer/blueprints/spec-workflows/]
//

type Workflow struct {
	*Entity `name:"workflow"`
	Name    string `namespace:""`

	Mapping              *string              `read:"mapping" require:"mapping"`
	ParameterDefinitions ParameterDefinitions `read:"parameters,ParameterDefinition"`
}

func NewWorkflow(context *tosca.Context) *Workflow {
	return &Workflow{
		Entity:               NewEntity(context),
		Name:                 context.Name,
		ParameterDefinitions: make(ParameterDefinitions),
	}
}

// tosca.Reader signature
func ReadWorkflow(context *tosca.Context) interface{} {
	self := NewWorkflow(context)

	if context.Is("map") {
		// Long notation
		context.ValidateUnsupportedFields(context.ReadFields(self))
	} else if context.ValidateType("map", "string") {
		// Short notation
		self.Mapping = context.FieldChild("mapping", context.Data).ReadString()
	}

	return self
}

func (self *Workflow) Normalize(s *normal.ServiceTemplate) *normal.Workflow {
	log.Infof("{normalize} workflow: %s", self.Name)

	w := s.NewWorkflow(self.Name)

	// TODO: mapping

	// TODO: support property definitions
	//self.ParameterDefinitions.Normalize(w.Inputs)

	return w
}

//
// Workflows
//

type Workflows []*Workflow

func (self Workflows) Normalize(s *normal.ServiceTemplate) {
	for _, workflow := range self {
		s.Workflows[workflow.Name] = workflow.Normalize(s)
	}
}
