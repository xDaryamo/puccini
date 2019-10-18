package compiler

import (
	"github.com/tliron/puccini/ard"
	"github.com/tliron/puccini/clout"
	"github.com/tliron/puccini/common"
	"github.com/tliron/puccini/tosca/normal"
)

func Compile(s *normal.ServiceTemplate) (*clout.Clout, error) {
	clout_ := clout.NewClout()

	metadata := make(ard.StringMap)
	for name, scriptlet := range s.ScriptletNamespace {
		scriptlet, err := scriptlet.Read()
		if err != nil {
			return nil, err
		}
		if err = ard.StringMapPutNested(metadata, name, scriptlet); err != nil {
			return nil, err
		}
	}
	clout_.Metadata["puccini-js"] = metadata

	history := ard.List{ard.StringMap{
		"timestamp":   common.Timestamp(),
		"description": "compile",
	}}

	metadata = make(ard.StringMap)
	metadata["version"] = VERSION
	metadata["history"] = history
	clout_.Metadata["puccini-tosca"] = metadata

	tosca := make(ard.StringMap)
	tosca["description"] = s.Description
	tosca["metadata"] = s.Metadata
	tosca["inputs"] = s.Inputs
	tosca["outputs"] = s.Outputs
	clout_.Properties["tosca"] = tosca

	nodeTemplates := make(map[string]*clout.Vertex)

	// Node templates
	for _, nodeTemplate := range s.NodeTemplates {
		v := clout_.NewVertex(clout.NewKey())

		nodeTemplates[nodeTemplate.Name] = v

		SetMetadata(v, "nodeTemplate")
		v.Properties["name"] = nodeTemplate.Name
		v.Properties["description"] = nodeTemplate.Description
		v.Properties["types"] = nodeTemplate.Types
		v.Properties["directives"] = nodeTemplate.Directives
		v.Properties["properties"] = nodeTemplate.Properties
		v.Properties["attributes"] = nodeTemplate.Attributes
		v.Properties["requirements"] = nodeTemplate.Requirements
		v.Properties["capabilities"] = nodeTemplate.Capabilities
		v.Properties["interfaces"] = nodeTemplate.Interfaces
		v.Properties["artifacts"] = nodeTemplate.Artifacts
	}

	groups := make(map[string]*clout.Vertex)

	// Groups
	for _, group := range s.Groups {
		v := clout_.NewVertex(clout.NewKey())

		groups[group.Name] = v

		SetMetadata(v, "group")
		v.Properties["name"] = group.Name
		v.Properties["description"] = group.Description
		v.Properties["types"] = group.Types
		v.Properties["properties"] = group.Properties
		v.Properties["interfaces"] = group.Interfaces

		for _, nodeTemplate := range group.Members {
			nv := nodeTemplates[nodeTemplate.Name]
			e := v.NewEdgeTo(nv)

			SetMetadata(e, "member")
		}
	}

	workflows := make(map[string]*clout.Vertex)

	// Workflows
	for _, workflow := range s.Workflows {
		v := clout_.NewVertex(clout.NewKey())

		workflows[workflow.Name] = v

		SetMetadata(v, "workflow")
		v.Properties["name"] = workflow.Name
		v.Properties["description"] = workflow.Description
	}

	// Workflow steps
	for name, workflow := range s.Workflows {
		v := workflows[name]

		steps := make(map[string]*clout.Vertex)

		for _, step := range workflow.Steps {
			sv := clout_.NewVertex(clout.NewKey())

			steps[step.Name] = sv

			SetMetadata(sv, "workflowStep")
			sv.Properties["name"] = step.Name

			e := v.NewEdgeTo(sv)
			SetMetadata(e, "workflowStep")

			if step.TargetNodeTemplate != nil {
				nv := nodeTemplates[step.TargetNodeTemplate.Name]
				e = sv.NewEdgeTo(nv)
				SetMetadata(e, "nodeTemplateTarget")
			} else if step.TargetGroup != nil {
				gv := groups[step.TargetGroup.Name]
				e = sv.NewEdgeTo(gv)
				SetMetadata(e, "groupTarget")
			} else {
				// This would happen only if there was a parsing error
				continue
			}

			// Workflow activities
			for sequence, activity := range step.Activities {
				av := clout_.NewVertex(clout.NewKey())

				e = sv.NewEdgeTo(av)
				SetMetadata(e, "workflowActivity")
				e.Properties["sequence"] = sequence

				SetMetadata(av, "workflowActivity")
				if activity.DelegateWorkflow != nil {
					wv := workflows[activity.DelegateWorkflow.Name]
					e = av.NewEdgeTo(wv)
					SetMetadata(e, "delegateWorflow")
				} else if activity.InlineWorkflow != nil {
					wv := workflows[activity.InlineWorkflow.Name]
					e = av.NewEdgeTo(wv)
					SetMetadata(e, "inlineWorflow")
				} else if activity.SetNodeState != "" {
					av.Properties["setNodeState"] = activity.SetNodeState
				} else if activity.CallOperation != nil {
					m := make(ard.StringMap)
					m["interface"] = activity.CallOperation.Interface.Name
					m["operation"] = activity.CallOperation.Name
					av.Properties["callOperation"] = m
				}
			}
		}

		for _, step := range workflow.Steps {
			sv := steps[step.Name]

			for _, next := range step.OnSuccessSteps {
				nsv := steps[next.Name]
				e := sv.NewEdgeTo(nsv)
				SetMetadata(e, "onSuccess")
			}

			for _, next := range step.OnFailureSteps {
				nsv := steps[next.Name]
				e := sv.NewEdgeTo(nsv)
				SetMetadata(e, "onFailure")
			}
		}
	}

	// Policies
	for _, policy := range s.Policies {
		v := clout_.NewVertex(clout.NewKey())

		SetMetadata(v, "policy")
		v.Properties["name"] = policy.Name
		v.Properties["description"] = policy.Description
		v.Properties["types"] = policy.Types
		v.Properties["properties"] = policy.Properties

		for _, nodeTemplate := range policy.NodeTemplateTargets {
			nv := nodeTemplates[nodeTemplate.Name]
			e := v.NewEdgeTo(nv)

			SetMetadata(e, "nodeTemplateTarget")
		}

		for _, group := range policy.GroupTargets {
			gv := groups[group.Name]
			e := v.NewEdgeTo(gv)

			SetMetadata(e, "groupTarget")
		}

		for _, trigger := range policy.Triggers {
			if trigger.Operation != nil {
				to := clout_.NewVertex(clout.NewKey())

				SetMetadata(to, "operation")
				to.Properties["description"] = trigger.Operation.Description
				to.Properties["implementation"] = trigger.Operation.Implementation
				to.Properties["dependencies"] = trigger.Operation.Dependencies
				to.Properties["inputs"] = trigger.Operation.Inputs

				e := v.NewEdgeTo(to)
				SetMetadata(e, "policyTriggerOperation")
			} else if trigger.Workflow != nil {
				wv := workflows[trigger.Workflow.Name]

				e := v.NewEdgeTo(wv)
				SetMetadata(e, "policyTriggerWorkflow")
			}
		}
	}

	// Substitution
	if s.Substitution != nil {
		v := clout_.NewVertex(clout.NewKey())

		SetMetadata(v, "substitution")
		v.Properties["type"] = s.Substitution.Type
		v.Properties["typeMetadata"] = s.Substitution.TypeMetadata

		for nodeTemplate, capability := range s.Substitution.CapabilityMappings {
			vv := nodeTemplates[nodeTemplate.Name]
			e := v.NewEdgeTo(vv)

			SetMetadata(e, "capabilityMapping")
			e.Properties["capability"] = capability.Name
		}

		for nodeTemplate, requirement := range s.Substitution.RequirementMappings {
			vv := nodeTemplates[nodeTemplate.Name]
			e := v.NewEdgeTo(vv)

			SetMetadata(e, "requirementMapping")
			e.Properties["requirement"] = requirement
		}

		for nodeTemplate, property := range s.Substitution.PropertyMappings {
			vv := nodeTemplates[nodeTemplate.Name]
			e := v.NewEdgeTo(vv)

			SetMetadata(e, "propertyMapping")
			e.Properties["property"] = property
		}

		for nodeTemplate, interface_ := range s.Substitution.InterfaceMappings {
			vv := nodeTemplates[nodeTemplate.Name]
			e := v.NewEdgeTo(vv)

			SetMetadata(e, "interfaceMapping")
			e.Properties["interface"] = interface_
		}
	}

	// Normalize
	var err error
	clout_, err = clout_.Normalize()
	if err != nil {
		return clout_, err
	}

	// TODO: call JavaScript plugins

	return clout_, nil
}

func SetMetadata(entity clout.Entity, kind string) {
	metadata := make(ard.StringMap)
	metadata["version"] = VERSION
	metadata["kind"] = kind
	entity.GetMetadata()["puccini-tosca"] = metadata
}
