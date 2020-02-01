// This file was auto-generated from a YAML file

package v1_0

func init() {
	Profile["/tosca/simple/1.0/profile.yaml"] = `
tosca_definitions_version: tosca_simple_yaml_1_0

metadata:
  puccini.scriptlet.import|tosca.resolve: internal:/tosca/common/1.0/js/resolve.js
  puccini.scriptlet.import|tosca.coerce: internal:/tosca/common/1.0/js/coerce.js
  puccini.scriptlet.import|tosca.utils: internal:/tosca/common/1.0/js/utils.js
  puccini.scriptlet.import|tosca.helpers: internal:/tosca/common/1.0/js/helpers.js

imports:
- artifacts.yaml
- groups.yaml
- nodes.yaml
- policies.yaml
`
}
