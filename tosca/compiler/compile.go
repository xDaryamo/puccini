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

	for key, nodeTemplate := range s.NodeTemplates {
		v := c.NewVertex(key)

		metadata = make(ard.Map)
		metadata["version"] = "1.0"
		metadata["kind"] = "nodeTemplate"
		v.Metadata["puccini-tosca"] = metadata

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

	for key, nodeTemplate := range s.NodeTemplates {
		v := c.Vertexes[key]

		for _, relationship := range nodeTemplate.Relationships {
			vv := c.Vertexes[relationship.TargetNodeTemplate.Name]
			e := v.NewEdgeTo(vv)

			metadata = make(ard.Map)
			metadata["version"] = "1.0"
			metadata["kind"] = "relationship"
			e.Metadata["puccini-tosca"] = metadata

			e.Properties["name"] = relationship.Name
			e.Properties["description"] = relationship.Description
			e.Properties["types"] = relationship.Types
			e.Properties["properties"] = relationship.Properties
			e.Properties["attributes"] = relationship.Attributes
			e.Properties["interfaces"] = relationship.Interfaces
		}
	}

	// TODO: groups and policies?

	// TODO: call JavaScript plugins

	return c, nil
}
