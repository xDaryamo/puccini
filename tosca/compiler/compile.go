package compiler

import (
	"github.com/tliron/puccini/ard"
	"github.com/tliron/puccini/clout"
	"github.com/tliron/puccini/common"
	"github.com/tliron/puccini/js"
	"github.com/tliron/puccini/tosca/normal"
)

func Compile(s *normal.ServiceTemplate) (*clout.Clout, error) {
	c := clout.NewClout()

	timestamp, err := common.Timestamp()
	if err != nil {
		return nil, err
	}

	metadata := make(ard.Map)
	for name, jsEntry := range s.ScriptNamespace {
		sourceCode, err := jsEntry.GetSourceCode()
		if err != nil {
			return nil, err
		}
		err = js.SetMapNested(metadata, name, sourceCode)
		if err != nil {
			return nil, err
		}
	}
	c.Metadata["puccini-js"] = metadata

	metadata = make(ard.Map)
	metadata["version"] = "1.0"
	metadata["history"] = []string{timestamp}
	c.Metadata["puccini-tosca"] = metadata

	tosca := make(ard.Map)
	tosca["metadata"] = s.Metadata
	tosca["inputs"] = s.Inputs
	tosca["outputs"] = s.Outputs
	c.Properties["tosca"] = tosca

	nodeTemplates := make(map[string]*clout.Vertex)

	// Node templates
	for _, nodeTemplate := range s.NodeTemplates {
		v := c.NewVertex(clout.NewKey())

		nodeTemplates[nodeTemplate.Name] = v

		SetMetadata(v, "nodeTemplate")
		v.Properties["name"] = nodeTemplate.Name
		v.Properties["description"] = nodeTemplate.Description
		v.Properties["types"] = nodeTemplate.Types
		v.Properties["directives"] = nodeTemplate.Directives
		v.Properties["properties"] = nodeTemplate.Properties
		v.Properties["attributes"] = nodeTemplate.Attributes
		v.Properties["capabilities"] = nodeTemplate.Capabilities
		v.Properties["interfaces"] = nodeTemplate.Interfaces
		v.Properties["artifacts"] = nodeTemplate.Artifacts
	}

	// Relationships
	for name, nodeTemplate := range s.NodeTemplates {
		v := nodeTemplates[name]

		for _, relationship := range nodeTemplate.Relationships {
			vv := nodeTemplates[relationship.TargetNodeTemplate.Name]
			e := v.NewEdgeTo(vv)

			SetMetadata(e, "relationship")
			e.Properties["name"] = relationship.Name
			e.Properties["description"] = relationship.Description
			e.Properties["types"] = relationship.Types
			e.Properties["properties"] = relationship.Properties
			e.Properties["attributes"] = relationship.Attributes
			e.Properties["interfaces"] = relationship.Interfaces
		}
	}

	// TODO: groups and policies?

	// Substitution
	if s.Substitution != nil {
		v := c.NewVertex(clout.NewKey())

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

	// TODO: workflows

	// TODO: call JavaScript plugins

	return c, nil
}

func SetMetadata(hasMetadata clout.HasMetadata, kind string) {
	metadata := make(ard.Map)
	metadata["version"] = "1.0"
	metadata["kind"] = kind
	hasMetadata.GetMetadata()["puccini-tosca"] = metadata
}
