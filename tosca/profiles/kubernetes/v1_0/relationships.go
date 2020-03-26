// This file was auto-generated from a YAML file

package v1_0

func init() {
	Profile["/tosca/kubernetes/1.0/relationships.yaml"] = `
tosca_definitions_version: tosca_simple_yaml_1_3

imports:

- capabilities.yaml

relationship_types:

  Route:
    valid_target_types: [ Service ] # capability
`
}
