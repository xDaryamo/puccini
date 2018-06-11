
clout.exec('tosca.helpers');

function validate(v, lower, upper) {
	if (arguments.length !== 3)
		throw 'must have 2 arguments';
	v = tosca.getComparable(v);
	lower = tosca.getComparable(lower);
	upper = tosca.getComparable(upper);
	return (v >= lower) && (v <= upper);
}
