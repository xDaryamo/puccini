
clout.exec('tosca.helpers');

function validate(v, bounds) {
	if (arguments.length !== 2)
		throw 'must have 1 arguments';
	if ((bounds.min === undefined) && (bounds.max === undefined))
		throw 'must provide "min" and/or "max"';
	v = tosca.getComparable(v);
	if (bounds.min !== undefined)
		if (v < tosca.getComparable(bounds.min))
			return false;
	if (bounds.max !== undefined)
		if (v > tosca.getComparable(bounds.max))
			return false;
	return true;
}
