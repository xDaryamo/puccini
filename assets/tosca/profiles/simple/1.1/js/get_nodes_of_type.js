
// [TOSCA-Simple-Profile-YAML-v1.1] @ 4.7.1

clout.exec('tosca.helpers');

function evaluate(typeName) {
	if (arguments.length !== 1)
		throw 'must have 1 argument';
	names = [];
	for (var name in clout.vertexes) {
		vertex = clout.vertexes[name];
		if (tosca.isTosca(vertex))
			names.push(vertex.properties.name);
	}
	return names;
}
