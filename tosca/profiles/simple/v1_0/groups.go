// This file was auto-generated from a YAML file

package v1_0

func init() {
	Profile["/tosca/simple/1.0/groups.yaml"] = `
# Modified from a file that was distributed with this NOTICE:
#
#   Apache AriaTosca
#   Copyright 2016-2017 The Apache Software Foundation
#
#   This product includes software developed at
#   The Apache Software Foundation (http://www.apache.org/).

tosca_definitions_version: tosca_simple_yaml_1_0

imports:
- interfaces.yaml

group_types:

  tosca.groups.Root:
    metadata:
      normative: 'true'
      citation: '[TOSCA-Simple-Profile-YAML-v1.0]'
      citation_location: 5.9.1
    description: >-
      This is the default (root) TOSCA Group Type definition that all other TOSCA base Group Types
      derive from.
    interfaces:
      Standard:
        type: tosca.interfaces.node.lifecycle.Standard
`
}
