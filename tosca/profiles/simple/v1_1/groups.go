// This file was auto-generated from a YAML file

package v1_1

func init() {
	Profile["/tosca/simple/1.1/groups.yaml"] = `
tosca_definitions_version: tosca_simple_yaml_1_1

imports:

- interfaces.yaml

group_types:

  tosca.groups.Root:
    metadata:
      puccini.normative: 'true'
      specification.citation: '[TOSCA-Simple-Profile-YAML-v1.1]'
      specification.location: 5.10.1
    description: >-
      This is the default (root) TOSCA Group Type definition that all other TOSCA base Group Types
      derive from.
    interfaces:
      Standard:
        type: tosca.interfaces.node.lifecycle.Standard
`
}
