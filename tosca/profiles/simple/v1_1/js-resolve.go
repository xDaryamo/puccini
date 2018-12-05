// This file was auto-generated from YAML files

package v1_1

func init() {
	Profile["/tosca/simple/1.1/js/resolve.js"] = `

clout.exec('tosca.utils');

for (v in clout.vertexes) {
	vertex = clout.vertexes[v];
	if (tosca.isNodeTemplate(vertex)) {
		nodeTemplate = vertex.properties;
	}
}


`
}
