
clout.exec('tosca.utils');

for (v in clout.vertexes) {
	vertex = clout.vertexes[v];
	if (tosca.isNodeTemplate(vertex)) {
		nodeTemplate = vertex.properties;
	}
}


