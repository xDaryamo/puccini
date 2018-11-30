
// [TOSCA-Simple-Profile-YAML-v1.1] @ 3.5.2

clout.exec('tosca.helpers');

function validate(v, lower, upper) {
	if (arguments.length !== 3)
		throw 'must have 2 arguments';
	v = tosca.getComparable(v);
	lower = tosca.getComparable(lower);
	upper = tosca.getComparable(upper);
	return (v >= lower) && (v <= upper);
}
