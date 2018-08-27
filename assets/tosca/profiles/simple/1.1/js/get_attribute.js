
clout.exec('tosca.helpers');

function evaluate(entity, attribute) {
	length = arguments.length;
	if (length < 2)
		throw 'must have at least 2 arguments';
	nodeTemplate = tosca.getNodeTemplate(entity);
	attributes = nodeTemplate.attributes;
	if (!(attribute in attributes))
		throw puccini.sprintf('attribute "%s" not found in "%s"', attribute, nodeTemplate.name);
	r = clout.coerce(attributes[attribute]);
	for (i = 2; i < length; i++)
		r = r[arguments[i]];
	return r;
}
