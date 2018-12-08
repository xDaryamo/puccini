
clout.exec('tosca.helpers');

tosca.prepare = function() {
	tosca.traverseValues(clout.newCoercible);
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

	for (vertexId in clout.vertexes) {
		vertex = clout.vertexes[vertexId];
		if (tosca.isNodeTemplate(vertex)) {
			nodeTemplate = vertex.properties;

			tosca.traverseObjectValues(traverser, nodeTemplate.properties, vertex);
			tosca.traverseObjectValues(traverser, nodeTemplate.attributes, vertex);
			tosca.traverseInterfaceValues(traverser, nodeTemplate.interfaces, vertex)

			for (capabilityName in nodeTemplate.capabilities) {
				capability = nodeTemplate.capabilities[capabilityName];
				tosca.traverseObjectValues(traverser, capability.properties, vertex);
				tosca.traverseObjectValues(traverser, capability.attributes, vertex);
			}

			for (artifactName in nodeTemplate.artifacts) {
				artifact = nodeTemplate.artifacts[artifactName];
				tosca.traverseObjectValues(traverser, artifact.properties, vertex);
			}

			for (e = 0; e < vertex.edgesOut.length; e++) {
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
	for (interfaceName in interfaces) {
		interface_ = interfaces[interfaceName];
		tosca.traverseObjectValues(traverser, interface_.inputs, site, source, target);
		for (operationName in interface_.operations)
			tosca.traverseObjectValues(traverser, interface_.operations[operationName].Inputs, site, source, target);
	}
};

tosca.traverseObjectValues = function(traverser, o, site, source, target) {
	for (k in o)
		o[k] = traverser(o[k], site, source, target);
};
