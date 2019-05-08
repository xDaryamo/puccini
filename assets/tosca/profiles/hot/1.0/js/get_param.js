
// [https://docs.openstack.org/heat/stein/template_guide/hot_spec.html#get_param]

clout.exec('tosca.helpers');

function evaluate(input) {
	if (arguments.length !== 1)
		throw 'must have 1 argument';
	if (!tosca.isTosca(clout))
		throw 'Clout is not TOSCA';
	var inputs = clout.properties.tosca.inputs;
	if (!(input in inputs))
		throw puccini.sprintf('parameter "%s" not found', input);
	var r = inputs[input];
	r = clout.coerce(r);
	return r;
}
