// This file was auto-generated from a YAML file

package v1_3

func init() {
	Profile["/tosca/simple/1.3/groups.yaml"] = `
tosca_definitions_version: tosca_simple_yaml_1_3

imports:

- interfaces.yaml

group_types:

  tosca.groups.Root:
    metadata:
      puccini.normative: 'true'
      specification.citation: '[TOSCA-Simple-Profile-YAML-v1.3]'
      specification.location: 5.10.1
    description: >-
      This is the default (root) TOSCA Group Type definition that all other TOSCA base Group Types
      derive from.
`
}
