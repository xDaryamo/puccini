// This file was auto-generated from YAML files

package v1_1

func init() {
	Profile["/tosca/simple/1.1/js/resolve.js"] = `

clout.exec('tosca.utils');

for (vertexId in clout.vertexes) {
	vertex = clout.vertexes[vertexId];
	if (!tosca.isNodeTemplate(vertex))
		continue;

	nodeTemplate = vertex.properties;
	requirements = nodeTemplate.requirements;
	for (r = 0; r < requirements.length; r++) {
		requirement = requirements[r];
		resolve(vertex, nodeTemplate, requirement);
	}
}

function resolve(sourceVertex, sourceNodeTemplate, requirement) {
	path = requirement.path;
	name = requirement.name;

	if (isSubstituted(sourceNodeTemplate.name, name)) {
		puccini.log.infof('{resolve} %s: skipping because in substitution mappings', path)
		return;
	}

	nodeTemplateName = requirement.nodeTemplateName;
	nodeTypeName = requirement.nodeTypeName;
	nodeTemplatePropertyConstraints = requirement.nodeTemplatePropertyConstraints;
	capabilityPropertyConstraintsMap = requirement.capabilityPropertyConstraints;

	// Gather candidate target node templates
	candidateTargetVertexes = [];
	for (vertexId in clout.vertexes) {
		vertex = clout.vertexes[vertexId];
		if (!tosca.isNodeTemplate(vertex))
			continue;

		candidateNodeTemplate = vertex.properties;
		candidateNodeTemplateName = candidateNodeTemplate.name;

		if ((nodeTemplateName !== '') && (nodeTemplateName !== candidateNodeTemplateName)) {
			puccini.log.debugf('{resolve} %s: node template "%s" is not named "%s"', path, candidateNodeTemplateName, nodeTemplateName);
			continue;
		}

		if (nodeTypeName !== '') {
			if (!(nodeTypeName in candidateNodeTemplate.types)) {
				puccini.log.debugf('{resolve} %s: node template "%s" is not of type "%s"', path, candidateNodeTemplateName, nodeTypeName);
				continue;
			}
		}

		if (nodeTemplatePropertyConstraints.length !== 0) {
			if (!arePropertiesValid(path, 'node template', candidateNodeTemplateName, candidateNodeTemplate, nodeTemplatePropertyConstraints)) {
				puccini.log.debugf('{resolve} %s: properties of node template "%s" do not match constraints', path, candidateNodeTemplateName);
				continue;
			}
		}

		if (capabilityPropertyConstraintsMap.length !== 0) {
			valid = true;
			candidateCapabilities = candidateNodeTemplate.capabilities;
			for (candidateCapabilityName in candidateCapabilities) {
				candidateCapability = candidateCapabilities[candidateCapabilityName];
				capabilityPropertyConstraints = capabilityPropertyConstraintsMap[candidateCapabilityName];
				if ((capabilityPropertyConstraints !== undefined) && (capabilityPropertyConstraints.length !== 0)) {
					if (!arePropertiesValid(path, 'capability', candidateCapabilityName, candidateCapability, capabilityPropertyConstraints)) {
						puccini.log.debugf('{resolve} %s: properties of capability "%s" in node template "%s" do not match constraints', path, candidateCapabilityName, candidateNodeTemplateName);
						valid = false;
						break;
					}
				}
			}
			if (!valid)
				continue;
		}

		candidateTargetVertexes.push(vertex);
	}

	if (candidateTargetVertexes.length === 0)
		unsatisfied(path, name, 'no candidate node template provides required capability');

	capabilityName = requirement.capabilityName;
	capabilityTypeName = requirement.capabilityTypeName;

	// Find first matching capability in candidate node templates
	satisfied = false;
	for (c in candidateTargetVertexes) {
		candidateTargetVertex = candidateTargetVertexes[c];
		candidateNodeTemplate = candidateTargetVertex.properties;
		candidateCapabilities = candidateNodeTemplate.capabilities;
		for (candidateCapabilityName in candidateCapabilities) {
			for (candidateCapabilityName in candidateCapabilities) {
				candidateCapability = candidateCapabilities[candidateCapabilityName];

				if ((capabilityName !== '') && (capabilityName !== candidateCapabilityName)) {
					puccini.log.debugf('{resolve} %s: capability "%s" in node template "%s" is not named "%s"', path, candidateCapabilityName, candidateNodeTemplateName, capabilityName);
					continue;
				}

				if (capabilityTypeName !== '') {
					if (!(capabilityTypeName in candidateCapability.types)) {
						puccini.log.debugf('{resolve} %s: capability "%s" in node template "%s" is not of type "%s"', path, candidateCapabilityName, candidateNodeTemplateName, capabilityTypeName);
						continue;
					}
				}

				// TODO: check that capability occurrences have not been filled

				puccini.log.infof('{resolve} %s: satisfied "%s" with capability "%s" in node template "%s"', path, name, candidateCapabilityName, candidateNodeTemplateName);
				addRelationship(sourceVertex, requirement, candidateTargetVertex, candidateCapabilityName);
				// return; GOJA: returning from inside for-loop is broken
				satisfied = true;
				break;
			}
			if (satisfied)
				break;
		}
		if (satisfied)
			break;
	}

	if (!satisfied)
		unsatisfied(path, name, 'no candidate node template provides required capability');
}

function addRelationship(sourceVertex, requirement, targetVertex, capabilityName) {
	edge = sourceVertex.newEdgeTo(targetVertex);
	edge.metadata['puccini-tosca'] = {
		version: '1.0',
		kind: 'relationship'
	};
	relationship = requirement.relationship;
	if (relationship)
		edge.properties = {
			name: requirement.name,
			description: relationship.description,
			types: relationship.types,
			properties: relationship.properties,
			attributes: relationship.attributes,
			interfaces: relationship.interfaces,
			capability: capabilityName
		};
	else
		edge.properties = {
			name: requirement.name,
			description: '',
			types: {},
			properties: {},
			attributes: {},
			interfaces: {},
			capability: capabilityName
		};
}

function arePropertiesValid(path, kind, name, entity, constraintsMap) {
	valid = true;

	properties = entity.properties;
	for (propertyName in constraintsMap) {
		puccini.log.debugf('{resolve} %s: applying constraints to property "%s" of %s "%s"', path, propertyName, kind, name);

		property = properties[propertyName];
		if (property === undefined) {
			valid = false;
			break; // return false; GOJA: returning from inside for-loop is broken
		}
		property = clout.newCoercible(property, null, null, null)

		constraints = constraintsMap[propertyName];
		constraints = clout.newConstraints(constraints, null, null, null)
		if (!constraints.validate(property)) {
			valid = false; // return false; GOJA: returning from inside for-loop is broken
			break;
		}
	}

	return valid;
}

function isSubstituted(nodeTemplateName, requirementName) {
	substituted = false;

	for (vertexId in clout.vertexes) {
		vertex = clout.vertexes[vertexId];
		if (!tosca.isTosca(vertex, 'substitution'))
			continue;

		for (e = 0; e < vertex.edgesOut.length; e++) {
			edge = vertex.edgesOut[e];
			if (!tosca.isTosca(edge, 'requirementMapping'))
				continue;

			if ((edge.target.properties.name === nodeTemplateName) && (edge.properties.requirement === requirementName)) {
				substituted = true;
				break;
			}
		}

		// There's only one substitution
		break;
	}

	return substituted;
}

function unsatisfied(path, name, message) {
	if (typeof problems === 'undefined')
		puccini.log.infof('%s: could not satisfy "%s" because %s', path, name, message);
	else
		problems.reportf('%s: could not satisfy "%s" because %s', path, name, message);
}
`
}
