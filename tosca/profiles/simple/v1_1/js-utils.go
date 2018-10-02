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

	for (v in clout.vertexes) {
		vertex = clout.vertexes[v];
		if (tosca.isNodeTemplate(vertex)) {
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
		} else if (tosca.isTosca(vertex, 'group')) {
			group = vertex.properties;

			tosca.traverseObjectValues(traverser, group.properties, vertex);
			tosca.traverseInterfaceValues(traverser, group.interfaces, vertex)
		} else if (tosca.isTosca(vertex, 'policy')) {
			policy = vertex.properties;

			tosca.traverseObjectValues(traverser, policy.properties, vertex);
		}
	}
};

tosca.traverseInterfaceValues = function(interfaces, site, source, target) {
	for (i in interfaces) {
		interface_ = interfaces[i];
		tosca.traverseObjectValues(traverser, interface_.inputs, site, source, target);
		for (o in interface_.operations)
			tosca.traverseObjectValues(traverser, interface_.operations[o].Inputs, site, source, target);
	}
};

tosca.traverseObjectValues = function(traverser, o, site, source, target) {
	for (k in o)
		o[k] = traverser(o[k], site, source, target);
};
`
}
