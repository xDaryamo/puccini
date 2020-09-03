// This file was auto-generated from a YAML file

package v1_0

func init() {
	Profile["/tosca/common/1.0/js/lib/traversal.js"] = `

clout.exec('tosca.lib.utils');

tosca.toCoercibles = function() {
	tosca.traverseValues(function(data) {
		return clout.newCoercible(data.value, data.site, data.source, data.target);
	});
};

tosca.unwrapCoercibles = function() {
	tosca.traverseValues(function(data) {
		return clout.unwrap(data.value);
	});
};

tosca.coerce = function() {
	tosca.toCoercibles();
	tosca.traverseValues(function(data) {
		return clout.coerce(data.value);
	});
};

tosca.getValueInformation = function() {
	var information = {};
	tosca.traverseValues(function(data) {
		if (data.value.$information)
			information[data.path.join('.')] = data.value.$information;
		return data.value;
	});
	return information;
};

tosca.traverseValues = function(traverser) {
	if (tosca.isTosca(clout)) {
		tosca.traverseObjectValues(traverser, ['inputs'], clout.properties.tosca.inputs);
		tosca.traverseObjectValues(traverser, ['outputs'], clout.properties.tosca.outputs);
	}

	for (var vertexId in clout.vertexes) {
		var vertex = clout.vertexes[vertexId];
		if (!tosca.isTosca(vertex))
			continue;

		if (tosca.isNodeTemplate(vertex)) {
			var nodeTemplate = vertex.properties;
			var path = ['nodeTemplates', nodeTemplate.name];

			tosca.traverseObjectValues(traverser, copyAndPush(path, 'properties'), nodeTemplate.properties, vertex);
			tosca.traverseObjectValues(traverser, copyAndPush(path, 'attributes'), nodeTemplate.attributes, vertex);
			tosca.traverseInterfaceValues(traverser, copyAndPush(path, 'interfaces'), nodeTemplate.interfaces, vertex)

			for (var capabilityName in nodeTemplate.capabilities) {
				var capability = nodeTemplate.capabilities[capabilityName];
				var capabilityPath = copyAndPush(path, 'capabilities', capabilityName);
				tosca.traverseObjectValues(traverser, copyAndPush(capabilityPath, 'properties'), capability.properties, vertex);
				tosca.traverseObjectValues(traverser, copyAndPush(capabilityPath, 'attributes'), capability.attributes, vertex);
			}

			for (var artifactName in nodeTemplate.artifacts) {
				var artifact = nodeTemplate.artifacts[artifactName];
				var artifactPath = copyAndPush(path, 'artifacts', artifactName);
				tosca.traverseObjectValues(traverser, copyAndPush(artifactPath, 'properties'), artifact.properties, vertex);
				if (artifact.credential !== null)
					try {
						artifact.credential = traverser({
							path: copyAndPush(artifactPath, 'credential'),
							value: artifact.credential,
							site: vertex
						});
					} catch (x) {
						if ((typeof problems !== 'undefined') && x.value && x.value.error)
							// Unwrap Go error
							problems.reportError(x.value);
						else
							throw x;
					}
			}

			for (var e = 0, l = vertex.edgesOut.length; e < l; e++) {
				var edge = vertex.edgesOut[e];
				if (!tosca.isTosca(edge, 'Relationship'))
					continue;

				var relationship = edge.properties;
				var relationshipPath = copyAndPush(path, 'relationships', relationship.name);
				tosca.traverseObjectValues(traverser, copyAndPush(relationshipPath, 'properties'), relationship.properties, edge, vertex, edge.target);
				tosca.traverseObjectValues(traverser,copyAndPush(relationshipPath, 'attributes'), relationship.attributes, edge, vertex, edge.target);
				tosca.traverseInterfaceValues(traverser, copyAndPush(relationshipPath, 'interfaces'), relationship.interfaces, edge, vertex, edge.target);
			}
		} else if (tosca.isTosca(vertex, 'Group')) {
			var group = vertex.properties;
			var path = ['groups', group.name];

			tosca.traverseObjectValues(traverser, copyAndPush(path, 'properties'), group.properties, vertex);
			tosca.traverseInterfaceValues(traverser, copyAndPush(path, 'attributes'), group.interfaces, vertex)
		} else if (tosca.isTosca(vertex, 'Policy')) {
			var policy = vertex.properties;
			var path = ['policies', policy.name];

			tosca.traverseObjectValues(traverser, copyAndPush(path, 'properties'), policy.properties, vertex);
		}
	}
};

tosca.traverseInterfaceValues = function(traverser, path, interfaces, site, source, target) {
	for (var interfaceName in interfaces) {
		var interface_ = interfaces[interfaceName];
		var interfacePath = copyAndPush(path, interfaceName)
		tosca.traverseObjectValues(traverser, copyAndPush(interfacePath, 'inputs'), interface_.inputs, site, source, target);
		for (var operationName in interface_.operations)
			tosca.traverseObjectValues(traverser, copyAndPush(interfacePath, 'operations', operationName), interface_.operations[operationName].inputs, site, source, target);
	}
};

tosca.traverseObjectValues = function(traverser, path, object, site, source, target) {
	for (var key in object)
		try {
			object[key] = traverser({
				path: copyAndPush(path, key),
				value: object[key],
				site: site,
				source: source,
				target: target
			});
		} catch (x) {
			if ((typeof problems !== 'undefined') && x.value && x.value.error)
				// Unwrap Go error
				problems.reportError(x.value);
			else
				throw x;
		}
};

function copyAndPush(array) {
	var array_ = [];
	for (var i = 0, l = array.length; i < l; i++)
		array_.push(array[i]);
	for (var i = 1, l = arguments.length; i < l; i++)
		array_.push(arguments[i]);
	return array_;
}
`
}
