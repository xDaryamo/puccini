
clout.exec('tosca.lib.traversal');

// This is a copy of the built-in get_input function source
// Except that we added a "* 2" to the returned result
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
	return r * 2;
}

// The "clout.define" API accepts the scriptlet source code as text
// So our little trick is to define the function above and then "stringify" it here
clout.define('tosca.function.get_input', "clout.exec('tosca.lib.utils');\n" + evaluate);

tosca.coerce();

puccini.write(clout);
