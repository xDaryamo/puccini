// This file was auto-generated from a YAML file

package v2_0

func init() {
	Profile["/tosca/implicit/2.0/js/functions/get_attribute.js"] = `

// [TOSCA-Simple-Profile-YAML-v1.3] @ 4.5.1
// [TOSCA-Simple-Profile-YAML-v1.2] @ 4.5.1
// [TOSCA-Simple-Profile-YAML-v1.1] @ 4.5.1
// [TOSCA-Simple-Profile-YAML-v1.0] @ 4.5.1

clout.exec('tosca.lib.utils');

function evaluate(entity, first) {
	return tosca.getNestedValue('attribute', 'attributes', arguments);
}
`
}
