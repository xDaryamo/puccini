
// [TOSCA-Simple-Profile-YAML-v1.2] @ 3.6.3
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.2

clout.exec('tosca.helpers');

function validate(v1, v2) {
	if (arguments.length !== 2)
		throw 'must have 1 argument';
	return tosca.getComparable(v1) > tosca.getComparable(v2);
}
