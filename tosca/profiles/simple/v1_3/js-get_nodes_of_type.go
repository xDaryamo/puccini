// This file was auto-generated from a YAML file

package v1_3

func init() {
	Profile["/tosca/simple/1.3/js/get_nodes_of_type.js"] = `

// [TOSCA-Simple-Profile-YAML-v1.3] @ 4.7.1
// [TOSCA-Simple-Profile-YAML-v1.2] @ 4.7.1
// [TOSCA-Simple-Profile-YAML-v1.1] @ 4.7.1

clout.exec('tosca.helpers');

function evaluate(typeName) {
	if (arguments.length !== 1)
		throw 'must have 1 argument';
	var names = [];
	for (var name in clout.vertexes) {
		var vertex = clout.vertexes[name];
		if (tosca.isTosca(vertex))
			names.push(vertex.properties.name);
	}
	return names;
}
`
}
