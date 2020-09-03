package compiler

import (
	"github.com/tliron/kutil/ard"
	"github.com/tliron/kutil/util"
	cloutpkg "github.com/tliron/puccini/clout"
	"github.com/tliron/puccini/clout/js"
	"github.com/tliron/puccini/tosca"
	"github.com/tliron/puccini/tosca/normal"
)

func Compile(serviceTemplate *normal.ServiceTemplate, allowTimestamps bool) (*cloutpkg.Clout, error) {
	clout := cloutpkg.NewClout()

	puccini := make(ard.StringMap)
	puccini["version"] = VERSION

	scriptlets := make(ard.StringMap)
	var err error = nil
	serviceTemplate.ScriptletNamespace.Range(func(name string, scriptlet *tosca.Scriptlet) bool {
		if scriptlet_, err_ := scriptlet.Read(); err == nil {
			if err_ = ard.StringMapPutNested(scriptlets, name, js.CleanupScriptlet(scriptlet_)); err_ != nil {
				err = err_
				return false
			}
		} else {
			err = err_
			return false
		}
		return true
	})
	if err != nil {
		return nil, err
	}
	puccini["scriptlets"] = scriptlets

	clout.Metadata["puccini"] = puccini

	history := ard.List{ard.StringMap{
		"timestamp":   util.Timestamp(!allowTimestamps),
		"description": "compile",
	}}
	clout.Metadata["history"] = history

	tosca := make(ard.StringMap)
	tosca["description"] = serviceTemplate.Description
	tosca["metadata"] = serviceTemplate.Metadata
	tosca["inputs"] = serviceTemplate.Inputs
	tosca["outputs"] = serviceTemplate.Outputs
	clout.Properties["tosca"] = tosca

	nodeTemplates := make(map[string]*cloutpkg.Vertex)

	// Node templates
	for _, nodeTemplate := range serviceTemplate.NodeTemplates {
		vertex := clout.NewVertex(cloutpkg.NewKey())

		nodeTemplates[nodeTemplate.Name] = vertex

		SetMetadata(vertex, "NodeTemplate")
		vertex.Properties["name"] = nodeTemplate.Name
		vertex.Properties["metadata"] = nodeTemplate.Metadata
		vertex.Properties["description"] = nodeTemplate.Description
		vertex.Properties["types"] = nodeTemplate.Types
		vertex.Properties["directives"] = nodeTemplate.Directives
		vertex.Properties["properties"] = nodeTemplate.Properties
		vertex.Properties["attributes"] = nodeTemplate.Attributes
		vertex.Properties["requirements"] = nodeTemplate.Requirements
		vertex.Properties["capabilities"] = nodeTemplate.Capabilities
		vertex.Properties["interfaces"] = nodeTemplate.Interfaces
		vertex.Properties["artifacts"] = nodeTemplate.Artifacts
	}

	groups := make(map[string]*cloutpkg.Vertex)

	// Groups
	for _, group := range serviceTemplate.Groups {
		vertex := clout.NewVertex(cloutpkg.NewKey())

		groups[group.Name] = vertex

		SetMetadata(vertex, "Group")
		vertex.Properties["name"] = group.Name
		vertex.Properties["metadata"] = group.Metadata
		vertex.Properties["description"] = group.Description
		vertex.Properties["types"] = group.Types
		vertex.Properties["properties"] = group.Properties
		vertex.Properties["interfaces"] = group.Interfaces

		for _, nodeTemplate := range group.Members {
			nodeTemplateVertex := nodeTemplates[nodeTemplate.Name]
			edge := vertex.NewEdgeTo(nodeTemplateVertex)

			SetMetadata(edge, "Member")
		}
	}

	workflows := make(map[string]*cloutpkg.Vertex)

	// Workflows
	for _, workflow := range serviceTemplate.Workflows {
		vertex := clout.NewVertex(cloutpkg.NewKey())

		workflows[workflow.Name] = vertex

		SetMetadata(vertex, "Workflow")
		vertex.Properties["name"] = workflow.Name
		vertex.Properties["description"] = workflow.Description
	}

	// Workflow steps
	for name, workflow := range serviceTemplate.Workflows {
		vertex := workflows[name]

		steps := make(map[string]*cloutpkg.Vertex)

		for _, step := range workflow.Steps {
			stepVertex := clout.NewVertex(cloutpkg.NewKey())

			steps[step.Name] = stepVertex

			SetMetadata(stepVertex, "WorkflowStep")
			stepVertex.Properties["name"] = step.Name

			edge := vertex.NewEdgeTo(stepVertex)
			SetMetadata(edge, "WorkflowStep")

			if step.TargetNodeTemplate != nil {
				nodeTemplateVertex := nodeTemplates[step.TargetNodeTemplate.Name]
				edge = stepVertex.NewEdgeTo(nodeTemplateVertex)
				SetMetadata(edge, "NodeTemplateTarget")
			} else if step.TargetGroup != nil {
				groupVertex := groups[step.TargetGroup.Name]
				edge = stepVertex.NewEdgeTo(groupVertex)
				SetMetadata(edge, "GroupTarget")
			} else {
				// This would happen only if there was a parsing error
				continue
			}

			// Workflow activities
			for sequence, activity := range step.Activities {
				activityVertex := clout.NewVertex(cloutpkg.NewKey())

				edge = stepVertex.NewEdgeTo(activityVertex)
				SetMetadata(edge, "WorkflowActivity")
				edge.Properties["sequence"] = sequence

				SetMetadata(activityVertex, "WorkflowActivity")
				if activity.DelegateWorkflow != nil {
					workflowVertex := workflows[activity.DelegateWorkflow.Name]
					edge = activityVertex.NewEdgeTo(workflowVertex)
					SetMetadata(edge, "DelegateWorkflow")
				} else if activity.InlineWorkflow != nil {
					workflowVertex := workflows[activity.InlineWorkflow.Name]
					edge = activityVertex.NewEdgeTo(workflowVertex)
					SetMetadata(edge, "InlineWorkflow")
				} else if activity.SetNodeState != "" {
					activityVertex.Properties["setNodeState"] = activity.SetNodeState
				} else if activity.CallOperation != nil {
					map_ := make(ard.StringMap)
					map_["interface"] = activity.CallOperation.Interface.Name
					map_["operation"] = activity.CallOperation.Name
					activityVertex.Properties["callOperation"] = map_
				}
			}
		}

		for _, step := range workflow.Steps {
			stepVertex := steps[step.Name]

			for _, next := range step.OnSuccessSteps {
				nextStepVertex := steps[next.Name]
				edge := stepVertex.NewEdgeTo(nextStepVertex)
				SetMetadata(edge, "OnSuccess")
			}

			for _, next := range step.OnFailureSteps {
				nextStepVertex := steps[next.Name]
				edge := stepVertex.NewEdgeTo(nextStepVertex)
				SetMetadata(edge, "OnFailure")
			}
		}
	}

	// Policies
	for _, policy := range serviceTemplate.Policies {
		vertex := clout.NewVertex(cloutpkg.NewKey())

		SetMetadata(vertex, "Policy")
		vertex.Properties["name"] = policy.Name
		vertex.Properties["metadata"] = policy.Metadata
		vertex.Properties["description"] = policy.Description
		vertex.Properties["types"] = policy.Types
		vertex.Properties["properties"] = policy.Properties

		for _, nodeTemplate := range policy.NodeTemplateTargets {
			nodeTemplateVertex := nodeTemplates[nodeTemplate.Name]
			edge := vertex.NewEdgeTo(nodeTemplateVertex)

			SetMetadata(edge, "NodeTemplateTarget")
		}

		for _, group := range policy.GroupTargets {
			groupVertex := groups[group.Name]
			edge := vertex.NewEdgeTo(groupVertex)

			SetMetadata(edge, "GroupTarget")
		}

		for _, trigger := range policy.Triggers {
			if trigger.Operation != nil {
				toVertex := clout.NewVertex(cloutpkg.NewKey())

				SetMetadata(toVertex, "Operation")
				toVertex.Properties["description"] = trigger.Operation.Description
				toVertex.Properties["implementation"] = trigger.Operation.Implementation
				toVertex.Properties["dependencies"] = trigger.Operation.Dependencies
				toVertex.Properties["inputs"] = trigger.Operation.Inputs

				edge := vertex.NewEdgeTo(toVertex)
				SetMetadata(edge, "PolicyTriggerOperation")
			} else if trigger.Workflow != nil {
				workflowVertex := workflows[trigger.Workflow.Name]

				edge := vertex.NewEdgeTo(workflowVertex)
				SetMetadata(edge, "PolicyTriggerWorkflow")
			}
		}
	}

	// Substitution
	if serviceTemplate.Substitution != nil {
		vertex := clout.NewVertex(cloutpkg.NewKey())

		SetMetadata(vertex, "Substitution")
		vertex.Properties["type"] = serviceTemplate.Substitution.Type
		vertex.Properties["typeMetadata"] = serviceTemplate.Substitution.TypeMetadata

		for nodeTemplate, capability := range serviceTemplate.Substitution.CapabilityMappings {
			nodeTemplateVertex := nodeTemplates[nodeTemplate.Name]
			edge := vertex.NewEdgeTo(nodeTemplateVertex)

			SetMetadata(edge, "CapabilityMapping")
			edge.Properties["capability"] = capability.Name
		}

		for nodeTemplate, requirement := range serviceTemplate.Substitution.RequirementMappings {
			nodeTemplateVertex := nodeTemplates[nodeTemplate.Name]
			edge := vertex.NewEdgeTo(nodeTemplateVertex)

			SetMetadata(edge, "RequirementMapping")
			edge.Properties["requirement"] = requirement
		}

		for nodeTemplate, property := range serviceTemplate.Substitution.PropertyMappings {
			nodeTemplateVertex := nodeTemplates[nodeTemplate.Name]
			edge := vertex.NewEdgeTo(nodeTemplateVertex)

			SetMetadata(edge, "PropertyMapping")
			edge.Properties["property"] = property
		}

		for nodeTemplate, attribute := range serviceTemplate.Substitution.AttributeMappings {
			nodeTemplateVertex := nodeTemplates[nodeTemplate.Name]
			edge := vertex.NewEdgeTo(nodeTemplateVertex)

			SetMetadata(edge, "AttributeMapping")
			edge.Properties["attribute"] = attribute
		}

		for nodeTemplate, interface_ := range serviceTemplate.Substitution.InterfaceMappings {
			nodeTemplateVertex := nodeTemplates[nodeTemplate.Name]
			edge := vertex.NewEdgeTo(nodeTemplateVertex)

			SetMetadata(edge, "InterfaceMapping")
			edge.Properties["interface"] = interface_
		}
	}

	// Normalize
	clout, err = clout.Normalize()
	if err != nil {
		return clout, err
	}

	// TODO: call JavaScript plugins

	return clout, nil
}

func SetMetadata(entity cloutpkg.Entity, kind string) {
	metadata := make(ard.StringMap)
	metadata["version"] = VERSION
	metadata["kind"] = kind
	entity.GetMetadata()["puccini"] = metadata
}
