
// [TOSCA-Simple-Profile-YAML-v1.2] @ 4.5.1
// [TOSCA-Simple-Profile-YAML-v1.1] @ 4.5.1

clout.exec('tosca.helpers');

function evaluate(entity, attribute) {
	var length = arguments.length;
	if (length < 2)
		throw 'must have at least 2 arguments';
	var nodeTemplate = tosca.getNodeTemplate(entity);
	var attributes = nodeTemplate.attributes;
	if (!(attribute in attributes))
		throw puccini.sprintf('attribute "%s" not found in "%s"', attribute, nodeTemplate.name);
	var r = clout.coerce(attributes[attribute]);
	for (var i = 2; i < length; i++)
		r = r[arguments[i]];
	return r;
}
