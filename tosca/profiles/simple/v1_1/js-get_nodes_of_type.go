// This file was auto-generated from YAML files

package v1_1

func init() {
	Profile["/tosca/simple/1.1/js/get_nodes_of_type.js"] = `

clout.exec('tosca.helpers');

function evaluate(typeName) {
	if (arguments.length !== 1)
		throw 'must have 1 argument';
	names = [];
	for (name in clout.vertexes) {
		vertex = clout.vertexes[name];
		if (tosca.isTosca(vertex))
			names.push(vertex.properties.name);
	}
	return names;
}
`
}
