// This file was auto-generated from a YAML file

package v1_3

func init() {
	Profile["/tosca/simple/1.3/js/functions/get_property.js"] = `

// [TOSCA-Simple-Profile-YAML-v1.3] @ 4.4.2
// [TOSCA-Simple-Profile-YAML-v1.2] @ 4.4.2
// [TOSCA-Simple-Profile-YAML-v1.1] @ 4.4.2
// [TOSCA-Simple-Profile-YAML-v1.0] @ 4.4.2

clout.exec('tosca.lib.utils');

function evaluate() {
	return tosca.getNestedValue('property', 'properties', arguments);
}
`
}
