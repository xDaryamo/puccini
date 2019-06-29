
clout.exec('tosca.helpers');

tosca.toCoercibles = function() {
	tosca.traverseValues(clout.newCoercible);
};

tosca.unwrapCoercibles = function() {
	tosca.traverseValues(clout.unwrap);
};

tosca.coerce = function() {
	tosca.toCoercibles();
	tosca.traverseValues(clout.coerce);
};

tosca.traverseValues = function(traverser) {
	if (tosca.isTosca(clout)) {
		tosca.traverseObjectValues(traverser, clout.properties.tosca.inputs);
		tosca.traverseObjectValues(traverser, clout.properties.tosca.outputs);
	}

	for (var vertexId in clout.vertexes) {
		var vertex = clout.vertexes[vertexId];
		if (tosca.isNodeTemplate(vertex)) {
			var nodeTemplate = vertex.properties;

			tosca.traverseObjectValues(traverser, nodeTemplate.properties, vertex);
			tosca.traverseObjectValues(traverser, nodeTemplate.attributes, vertex);
			tosca.traverseInterfaceValues(traverser, nodeTemplate.interfaces, vertex)

			for (var capabilityName in nodeTemplate.capabilities) {
				var capability = nodeTemplate.capabilities[capabilityName];
				tosca.traverseObjectValues(traverser, capability.properties, vertex);
				tosca.traverseObjectValues(traverser, capability.attributes, vertex);
			}

			for (var artifactName in nodeTemplate.artifacts) {
				var artifact = nodeTemplate.artifacts[artifactName];
				tosca.traverseObjectValues(traverser, artifact.properties, vertex);
			}

			for (var e = 0; e < vertex.edgesOut.length; e++) {
				var edge = vertex.edgesOut[e];
				if (!tosca.isTosca(edge, 'relationship'))
					continue;

				var relationship = edge.properties;
				tosca.traverseObjectValues(traverser, relationship.properties, edge, vertex, edge.target);
				tosca.traverseObjectValues(traverser, relationship.attributes, edge, vertex, edge.target);
				tosca.traverseInterfaceValues(traverser, relationship.interfaces, edge, vertex, edge.target);
			}
		} else if (tosca.isTosca(vertex, 'group')) {
			var group = vertex.properties;

			tosca.traverseObjectValues(traverser, group.properties, vertex);
			tosca.traverseInterfaceValues(traverser, group.interfaces, vertex)
		} else if (tosca.isTosca(vertex, 'policy')) {
			var policy = vertex.properties;

			tosca.traverseObjectValues(traverser, policy.properties, vertex);
		}
	}
};

tosca.traverseInterfaceValues = function(traverser, interfaces, site, source, target) {
	for (var interfaceName in interfaces) {
		var interface_ = interfaces[interfaceName];
		tosca.traverseObjectValues(traverser, interface_.inputs, site, source, target);
		for (var operationName in interface_.operations)
			tosca.traverseObjectValues(traverser, interface_.operations[operationName].inputs, site, source, target);
	}
};

tosca.traverseObjectValues = function(traverser, o, site, source, target) {
	for (var k in o)
		o[k] = traverser(o[k], site, source, target);
};
