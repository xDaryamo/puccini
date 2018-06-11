// This file was auto-generated from YAML files

package v1_10

func init() {
	Profile["/tosca/kubernetes/1.10/profile.yaml"] = `
tosca_definitions_version: tosca_simple_yaml_1_1

metadata:
  puccini-js.import.kubernetes.generate: js/generate.js
  puccini-js.import.kubernetes.update: js/update.js

imports:
- artifacts.yaml
- groups.yaml
- policies.yaml
`
}
