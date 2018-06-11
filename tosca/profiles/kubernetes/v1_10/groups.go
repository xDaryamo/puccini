// This file was auto-generated from YAML files

package v1_10

func init() {
	Profile["/tosca/kubernetes/1.10/groups.yaml"] = `
tosca_definitions_version: tosca_simple_yaml_1_1

imports:
- nodes.yaml

group_types:

  kubernetes.Namespace:
    description: >-
      Will automatically use a "group" label (the name of the group) for all deployment controllers.
    derived_from: tosca.groups.Root
    members:
    - tosca.nodes.Root
    properties:
      namespace:
        type: string
`
}
