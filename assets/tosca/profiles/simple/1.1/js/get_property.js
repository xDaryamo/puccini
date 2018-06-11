
clout.exec('tosca.helpers');

function evaluate(entity, property) {
	length = arguments.length;
	if (length < 2)
		throw 'must have at least 2 arguments';
	nodeTemplate = tosca.getNodeTemplate(entity);
	properties = nodeTemplate.properties;
	if (!(property in properties))
		throw sprintf('property "%s" not found in "%s"', property, nodeTemplate.name);
	r = clout.coerce(properties[property]);
	for (i = 2; i < length; i++)
		r = r[arguments[i]];
	return r;
}
