
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.2

clout.exec('tosca.helpers');

function validate(v1, v2) {
	if (arguments.length !== 2)
		throw 'must have 1 argument';
	v1 = tosca.getComparable(v1);
	v2 = tosca.getComparable(v2);
	return v1 >= v2;
}
