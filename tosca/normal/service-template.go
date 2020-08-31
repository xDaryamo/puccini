package normal

import (
	"github.com/tliron/kutil/reflection"
	"github.com/tliron/puccini/tosca"
)

//
// ServiceTemplate
//

type ServiceTemplate struct {
	Description        string                    `json:"description" yaml:"description"`
	NodeTemplates      NodeTemplates             `json:"nodeTemplates" yaml:"nodeTemplates"`
	Groups             Groups                    `json:"groups" yaml:"groups"`
	Policies           Policies                  `json:"policies" yaml:"policies"`
	Inputs             Constrainables            `json:"inputs" yaml:"inputs"`
	Outputs            Constrainables            `json:"outputs" yaml:"outputs"`
	Workflows          Workflows                 `json:"workflows" yaml:"workflows"`
	Substitution       *Substitution             `json:"substitution" yaml:"substitution"`
	Metadata           map[string]string         `json:"metadata" yaml:"metadata"`
	ScriptletNamespace *tosca.ScriptletNamespace `json:"scriptletNamespace" yaml:"scriptletNamespace"`
}

func NewServiceTemplate() *ServiceTemplate {
	return &ServiceTemplate{
		NodeTemplates:      make(NodeTemplates),
		Groups:             make(Groups),
		Policies:           make(Policies),
		Inputs:             make(Constrainables),
		Outputs:            make(Constrainables),
		Workflows:          make(Workflows),
		Metadata:           make(map[string]string),
		ScriptletNamespace: tosca.NewScriptletNamespace(),
	}
}

// From Normalizable interface
func NormalizeServiceTemplate(entityPtr tosca.EntityPtr) (*ServiceTemplate, bool) {
	var s *ServiceTemplate

	reflection.Traverse(entityPtr, func(entityPtr tosca.EntityPtr) bool {
		if normalizable, ok := entityPtr.(Normalizable); ok {
			s = normalizable.NormalizeServiceTemplate()

			// Only one entity should implement the interface
			return false
		}
		return true
	})

	if s == nil {
		return nil, false
	}

	return s, true
}

//
// Normalizable
//

type Normalizable interface {
	NormalizeServiceTemplate() *ServiceTemplate
}
