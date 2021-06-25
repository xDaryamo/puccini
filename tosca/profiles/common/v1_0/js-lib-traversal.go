// This file was auto-generated from a YAML file

package v1_0

func init() {
	Profile["/tosca/common/1.0/js/lib/traversal.js"] = `

const tosca = require('tosca.lib.utils');

exports.toCoercibles = function(clout_) {
	if (!clout_)
		clout_ = clout;
	exports.traverseValues(clout_, function(data) {
		return clout_.newCoercible(data.value, data.site, data.source, data.target);
	});
};

exports.unwrapCoercibles = function(clout_) {
	if (!clout_)
		clout_ = clout;
	exports.traverseValues(clout_, function(data) {
		return clout_.unwrap(data.value);
	});
};

exports.coerce = function(clout_) {
	if (!clout_)
		clout_ = clout;
	exports.toCoercibles(clout_);
	exports.traverseValues(clout_, function(data) {
		return clout_.coerce(data.value);
	});
};

exports.getValueInformation = function(clout_) {
	if (!clout_)
		clout_ = clout;
	let information = {};
	exports.traverseValues(clout_, function(data) {
		if (data.value.$information)
			information[data.path.join('.')] = data.value.$information;
		return data.value;
	});
	return information;
};

exports.traverseValues = function(clout_, traverser) {
	if (!clout_)
		clout_ = clout;

	if (tosca.isTosca(clout_)) {
		exports.traverseObjectValues(traverser, ['inputs'], clout_.properties.tosca.inputs);
		exports.traverseObjectValues(traverser, ['outputs'], clout_.properties.tosca.outputs);
	}

	for (let vertexId in clout_.vertexes) {
		let vertex = clout_.vertexes[vertexId];
		if (!tosca.isTosca(vertex))
			continue;

		if (tosca.isNodeTemplate(vertex)) {
			let nodeTemplate = vertex.properties;
			let path = ['nodeTemplates', nodeTemplate.name];

			exports.traverseObjectValues(traverser, copyAndPush(path, 'properties'), nodeTemplate.properties, vertex);
			exports.traverseObjectValues(traverser, copyAndPush(path, 'attributes'), nodeTemplate.attributes, vertex);
			exports.traverseInterfaceValues(traverser, copyAndPush(path, 'interfaces'), nodeTemplate.interfaces, vertex)

			for (let capabilityName in nodeTemplate.capabilities) {
				let capability = nodeTemplate.capabilities[capabilityName];
				let capabilityPath = copyAndPush(path, 'capabilities', capabilityName);
				exports.traverseObjectValues(traverser, copyAndPush(capabilityPath, 'properties'), capability.properties, vertex);
				exports.traverseObjectValues(traverser, copyAndPush(capabilityPath, 'attributes'), capability.attributes, vertex);
			}

			for (let artifactName in nodeTemplate.artifacts) {
				let artifact = nodeTemplate.artifacts[artifactName];
				let artifactPath = copyAndPush(path, 'artifacts', artifactName);
				exports.traverseObjectValues(traverser, copyAndPush(artifactPath, 'properties'), artifact.properties, vertex);
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

			for (let e = 0, l = vertex.edgesOut.length; e < l; e++) {
				let edge = vertex.edgesOut[e];
				if (!tosca.isTosca(edge, 'Relationship'))
					continue;

				let relationship = edge.properties;
				let relationshipPath = copyAndPush(path, 'relationships', relationship.name);
				exports.traverseObjectValues(traverser, copyAndPush(relationshipPath, 'properties'), relationship.properties, edge, vertex, edge.target);
				exports.traverseObjectValues(traverser,copyAndPush(relationshipPath, 'attributes'), relationship.attributes, edge, vertex, edge.target);
				exports.traverseInterfaceValues(traverser, copyAndPush(relationshipPath, 'interfaces'), relationship.interfaces, edge, vertex, edge.target);
			}
		} else if (tosca.isTosca(vertex, 'Group')) {
			let group = vertex.properties;
			let path = ['groups', group.name];

			exports.traverseObjectValues(traverser, copyAndPush(path, 'properties'), group.properties, vertex);
			exports.traverseInterfaceValues(traverser, copyAndPush(path, 'attributes'), group.interfaces, vertex)
		} else if (tosca.isTosca(vertex, 'Policy')) {
			let policy = vertex.properties;
			let path = ['policies', policy.name];

			exports.traverseObjectValues(traverser, copyAndPush(path, 'properties'), policy.properties, vertex);
		} else if (tosca.isTosca(vertex, 'Substitution')) {
			let substitution = vertex.properties;
			let path = ['substitution'];

			exports.traverseObjectValues(traverser, copyAndPush(path, 'properties'), substitution.properties, vertex);
		}
	}
};

exports.traverseInterfaceValues = function(traverser, path, interfaces, site, source, target) {
	for (let interfaceName in interfaces) {
		let interface_ = interfaces[interfaceName];
		let interfacePath = copyAndPush(path, interfaceName)
		exports.traverseObjectValues(traverser, copyAndPush(interfacePath, 'inputs'), interface_.inputs, site, source, target);
		for (let operationName in interface_.operations)
			exports.traverseObjectValues(traverser, copyAndPush(interfacePath, 'operations', operationName), interface_.operations[operationName].inputs, site, source, target);
	}
};

exports.traverseObjectValues = function(traverser, path, object, site, source, target) {
	for (let key in object)
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
	let array_ = [];
	for (let i = 0, l = array.length; i < l; i++)
		array_.push(array[i]);
	for (let i = 1, l = arguments.length; i < l; i++)
		array_.push(arguments[i]);
	return array_;
}
`
}
