// This file was auto-generated from a YAML file

package v1_0

func init() {
	Profile["/tosca/bpmn/1.0/policies.yaml"] = `
tosca_definitions_version: tosca_simple_yaml_1_3

policy_types:

  Base:
    description: >-
      Root for policies implemented by BPM software.

  Process:
    description: >-
      Policy implemented by a process defined in BPMN.
    derived_from: Base
    properties:
      bpmn_process_id:
        description: >-
          Execute this BPMN process when triggered.
        type: string
    targets:
    - tosca:Root
`
}
