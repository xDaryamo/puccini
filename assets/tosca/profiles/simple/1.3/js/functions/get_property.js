
// [TOSCA-Simple-Profile-YAML-v1.3] @ 4.4.2
// [TOSCA-Simple-Profile-YAML-v1.2] @ 4.4.2
// [TOSCA-Simple-Profile-YAML-v1.1] @ 4.4.2
// [TOSCA-Simple-Profile-YAML-v1.0] @ 4.4.2

clout.exec('tosca.lib.utils');

function evaluate() {
	return tosca.getNestedValue('property', 'properties', arguments);
}
