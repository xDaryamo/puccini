package normal

import (
	"github.com/tliron/puccini/tosca"
)

//
// ServiceTemplate
//

type ServiceTemplate struct {
	Description        string                   `json:"description" yaml:"description"`
	NodeTemplates      NodeTemplates            `json:"nodeTemplates" yaml:"nodeTemplates"`
	Groups             Groups                   `json:"groups" yaml:"groups"`
	Policies           Policies                 `json:"policies" yaml:"policies"`
	Inputs             Constrainables           `json:"inputs" yaml:"inputs"`
	Outputs            Constrainables           `json:"outputs" yaml:"outputs"`
	Workflows          Workflows                `json:"workflows" yaml:"workflows"`
	Substitution       *Substitution            `json:"substitution" yaml:"substitution"`
	Metadata           map[string]string        `json:"metadata" yaml:"metadata"`
	ScriptletNamespace tosca.ScriptletNamespace `json:"scriptletNamespace" yaml:"scriptletNamespace"`
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
		ScriptletNamespace: make(tosca.ScriptletNamespace),
	}
}
