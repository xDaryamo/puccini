// This file was auto-generated from YAML files

package v1_1

func init() {
	Profile["/tosca/simple/1.1/js/get_input.js"] = `

// [TOSCA-Simple-Profile-YAML-v1.1] @ 4.4.1

clout.exec('tosca.helpers');

function evaluate(input) {
	if (arguments.length !== 1)
		throw 'must have 1 argument';
	if (!tosca.isTosca(clout))
		throw 'Clout is not TOSCA';
	var inputs = clout.properties.tosca.inputs;
	if (!(input in inputs))
		throw puccini.sprintf('input "%s" not found', input);
	var r = inputs[input];
	r = clout.coerce(r);
	return r;
}
`
}
