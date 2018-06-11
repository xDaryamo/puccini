// This file was auto-generated from YAML files

package v1_1

func init() {
	Profile["/tosca/simple/1.1/js/utils.js"] = `

clout.exec('tosca.helpers');

tosca.prepare = function() {
	tosca.traverseValues(clout.prepare);
};

tosca.coerce = function() {
	tosca.prepare();
	tosca.traverseValues(clout.coerce);
};

tosca.traverseValues = function(traverser) {
	if (tosca.isTosca(clout)) {
		tosca.traverseObjectValues(traverser, clout.properties.tosca.inputs);
		tosca.traverseObjectValues(traverser, clout.properties.tosca.outputs);
	}

	for (name in clout.vertexes) {
		vertex = clout.vertexes[name];
		if (!tosca.isNodeTemplate(vertex))
			continue;
		nodeTemplate = vertex.properties;

		tosca.traverseObjectValues(traverser, nodeTemplate.properties, vertex);
		tosca.traverseObjectValues(traverser, nodeTemplate.attributes, vertex);
		tosca.traverseInterfaceValues(traverser, nodeTemplate.interfaces, vertex)

		for (c in nodeTemplate.capabilities) {
			capability = nodeTemplate.capabilities[c];
			tosca.traverseObjectValues(traverser, capability.properties, vertex);
			tosca.traverseObjectValues(traverser, capability.attributes, vertex);
		}

		for (a in nodeTemplate.artifacts) {
			artifact = nodeTemplate.artifacts[a];
			tosca.traverseObjectValues(traverser, artifact.properties, vertex);
		}

		for (e in vertex.edgesOut) {
			edge = vertex.edgesOut[e];
			if (!tosca.isTosca(edge, 'relationship'))
				continue;
			relationship = edge.properties;
			tosca.traverseObjectValues(traverser, relationship.properties, edge, vertex, edge.target);
			tosca.traverseObjectValues(traverser, relationship.attributes, edge, vertex, edge.target);
			tosca.traverseInterfaceValues(traverser, relationship.interfaces, edge, vertex, edge.target);
		}
	}
};

tosca.traverseInterfaceValues = function(interfaces, site, source, target) {
	for (i in interfaces) {
		intr = interfaces[i];
		tosca.traverseObjectValues(traverser, intr.inputs, site, source, target);
		for (o in intr.operations)
			tosca.traverseObjectValues(traverser, intr.operations[o].Inputs, site, source, target);
	}
};

tosca.traverseObjectValues = function(traverser, o, site, source, target) {
	for (k in o)
		o[k] = traverser(o[k], site, source, target);
};
`
}
