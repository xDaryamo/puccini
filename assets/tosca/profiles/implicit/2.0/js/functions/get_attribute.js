
// [TOSCA-Simple-Profile-YAML-v1.3] @ 4.5.1
// [TOSCA-Simple-Profile-YAML-v1.2] @ 4.5.1
// [TOSCA-Simple-Profile-YAML-v1.1] @ 4.5.1
// [TOSCA-Simple-Profile-YAML-v1.0] @ 4.5.1

clout.exec('tosca.lib.utils');

function evaluate(entity, first) {
	return tosca.getNestedValue('attribute', 'attributes', arguments);
}
