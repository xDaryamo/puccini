// This file was auto-generated from a YAML file

package v1_3

func init() {
	Profile["/tosca/simple/1.3/js/get_property.js"] = `

// [TOSCA-Simple-Profile-YAML-v1.3] @ 4.4.2
// [TOSCA-Simple-Profile-YAML-v1.2] @ 4.4.2
// [TOSCA-Simple-Profile-YAML-v1.1] @ 4.4.2
// [TOSCA-Simple-Profile-YAML-v1.0] @ 4.4.2

clout.exec('tosca.helpers');

function evaluate(entity, first) {
	var length = arguments.length;
	if (length < 2)
		throw 'must have at least 2 arguments';
	var nodeTemplate = tosca.getNodeTemplate(entity);
	var p = 1;
	var value = nodeTemplate.properties;
	if (first in nodeTemplate.capabilities) {
		value = nodeTemplate.capabilities[first].properties;
		first = arguments[++p];
	} else for (var r = 0; r < nodeTemplate.requirements.length; r++) {
		var requirement = nodeTemplate.requirements[r];
		if ((requirement.name === first) && (requirement.relationship !== undefined)) {
			values = requirement.relationship.properties;
			first = arguments[++p];
			break;
		}
	}
	if (first in value)
		value = value[first];
	else
		throw puccini.sprintf('property "%s" not found in "%s"', first, nodeTemplate.name);
	value = clout.coerce(value);
	for (var i = p + 1; i < length; i++)
		value = value[arguments[i]];
	return value;
}
`
}
