// This file was auto-generated from a YAML file

package v1_0

func init() {
	Profile["/tosca/kubernetes/1.0/profile.yaml"] = `
tosca_definitions_version: tosca_simple_yaml_1_3

metadata:
  puccini-js.import.kubernetes.generate: js/generate.js
  puccini-js.import.kubernetes.update: js/update.js

imports:
- artifacts.yaml
- groups.yaml
- policies.yaml
`
}
