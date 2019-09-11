// This file was auto-generated from a YAML file

package v1_3

func init() {
	Profile["/tosca/simple/1.3/js/get_property.js"] = `

// [TOSCA-Simple-Profile-YAML-v1.3] @ 4.4.2
// [TOSCA-Simple-Profile-YAML-v1.2] @ 4.4.2
// [TOSCA-Simple-Profile-YAML-v1.1] @ 4.4.2

clout.exec('tosca.helpers');

function evaluate(entity, property) {
	var length = arguments.length;
	if (length < 2)
		throw 'must have at least 2 arguments';
	var nodeTemplate = tosca.getNodeTemplate(entity);
	var properties = nodeTemplate.properties;
	if (!(property in properties))
		throw puccini.sprintf('property "%s" not found in "%s"', property, nodeTemplate.name);
	var r = clout.coerce(properties[property]);
	for (var i = 2; i < length; i++)
		r = r[arguments[i]];
	return r;
}
`
}
