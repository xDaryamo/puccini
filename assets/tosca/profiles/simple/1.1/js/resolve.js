
clout.exec('tosca.utils');

nodeTemplateVertexes = [];
for (var vertexId in clout.vertexes) {
	vertex = clout.vertexes[vertexId];
	if (tosca.isNodeTemplate(vertex))
		nodeTemplateVertexes.push(vertex);
}

// For consistent results, we will sort the node templates by name
nodeTemplateVertexes.sort(function(a, b) {
	return a.properties.name < b.properties.name ? -1 : 1;
});

// Resolve all requirements
for (var v = 0; v < nodeTemplateVertexes.length; v++) {
	vertex = nodeTemplateVertexes[v];
	nodeTemplate = vertex.properties;
	requirements = nodeTemplate.requirements;
	for (var r = 0; r < requirements.length; r++) {
		requirement = requirements[r];
		resolve(vertex, nodeTemplate, requirement);
	}
}

// Check that all capabilities have their minimum relationship count
for (var v = 0; v < nodeTemplateVertexes.length; v++) {
	vertex = nodeTemplateVertexes[v];
	nodeTemplate = vertex.properties;
	capabilities = nodeTemplate.capabilities;
	for (var capabilityName in capabilities) {
		capability = capabilities[capabilityName];
		relationshipCount = countRelationships(vertex, capabilityName);
		minRelationshipCount = capability.minRelationshipCount;
		if (relationshipCount < minRelationshipCount)
			notEnoughRelationships(nodeTemplate.name, capabilityName, relationshipCount, minRelationshipCount)
	}
}

function resolve(sourceVertex, sourceNodeTemplate, requirement) {
	path = requirement.path;
	name = requirement.name;

	if (isSubstituted(sourceNodeTemplate.name, name)) {
		puccini.log.infof('{resolve} %s: skipping because in substitution mappings', path)
		return;
	}

	candidates = gatherCandidateNodeTemplates(requirement);
	if (candidates.length === 0) {
		unsatisfied(path, name, 'there are no candidate node templates');
		return;
	}

	candidates = gatherCandidateCapabilities(requirement, candidates);
	if (candidates.length === 0) {
		unsatisfied(path, name, 'no candidate node template provides required capability');
		return;
	}

	// Gather priority candidates: those that have not yet fulfilled their minimum relationship count
	priorityCandidates = [];
	for (var c = 0; c < candidates.length; c++) {
		candidate = candidates[c];
		if ((candidate.capability.minRelationshipCount !== 0) && (countRelationships(candidate.vertex, candidate.capabilityName) < candidate.capability.minRelationshipCount))
			priorityCandidates.push(candidate);
	}

	chosen = null;

	if (priorityCandidates.length !== 0)
		// Of the priority candidates, pick the one with the highest minimum relationship count
		// (needs to be fulfilled soonest)
		for (var c = 0; c < priorityCandidates.length; c++) {
			candidate = priorityCandidates[c];
			if ((chosen === null) || (candidate.capability.minRelationshipCount > chosen.capability.minRelationshipCount))
				chosen = candidate;
		}
	else
		// Of the candidates, pick the one with highest maximum relationship count
		// (has the most room)
		for (var c = 0; c < candidates.length; c++) {
			candidate = candidates[c];
			if ((chosen === null) || isMaxCountGreater(candidate.capability.maxRelationshipCount, chosen.capability.maxRelationshipCount))
				chosen = candidate;
		}

	puccini.log.infof('{resolve} %s: satisfied "%s" with capability "%s" in node template "%s"', path, name, chosen.capabilityName, chosen.nodeTemplateName);
	addRelationship(sourceVertex, requirement, chosen.vertex, chosen.capabilityName);
}

function gatherCandidateNodeTemplates(requirement) {
	path = requirement.path;
	nodeTemplateName = requirement.nodeTemplateName;
	nodeTypeName = requirement.nodeTypeName;
	nodeTemplatePropertyConstraints = requirement.nodeTemplatePropertyConstraints;
	capabilityPropertyConstraintsMap = requirement.capabilityPropertyConstraints;

	candidates = [];
	for (var v = 0; v < nodeTemplateVertexes.length; v++) {
		vertex = nodeTemplateVertexes[v];
		candidateNodeTemplate = vertex.properties;
		candidateNodeTemplateName = candidateNodeTemplate.name;

		if ((nodeTemplateName !== '') && (nodeTemplateName !== candidateNodeTemplateName)) {
			puccini.log.debugf('{resolve} %s: node template "%s" is not named "%s"', path, candidateNodeTemplateName, nodeTemplateName);
			continue;
		}

		if ((nodeTypeName !== '') && !(nodeTypeName in candidateNodeTemplate.types)) {
			puccini.log.debugf('{resolve} %s: node template "%s" is not of type "%s"', path, candidateNodeTemplateName, nodeTypeName);
			continue;
		}

		// Node filter
		if ((nodeTemplatePropertyConstraints.length !== 0) && !arePropertiesValid(path, 'node template', candidateNodeTemplateName, candidateNodeTemplate, nodeTemplatePropertyConstraints)) {
			puccini.log.debugf('{resolve} %s: properties of node template "%s" do not match constraints', path, candidateNodeTemplateName);
			continue;
		}

		candidateCapabilities = candidateNodeTemplate.capabilities;

		// Capability filter
		if (capabilityPropertyConstraintsMap.length !== 0) {
			valid = true;
			for (var candidateCapabilityName in candidateCapabilities) {
				candidateCapability = candidateCapabilities[candidateCapabilityName];
				capabilityPropertyConstraints = capabilityPropertyConstraintsMap[candidateCapabilityName];
				if ((capabilityPropertyConstraints !== undefined) && (capabilityPropertyConstraints.length !== 0) && !arePropertiesValid(path, 'capability', candidateCapabilityName, candidateCapability, capabilityPropertyConstraints)) {
					puccini.log.debugf('{resolve} %s: properties of capability "%s" in node template "%s" do not match constraints', path, candidateCapabilityName, candidateNodeTemplateName);
					valid = false;
					break;
				}
			}
			if (!valid)
				continue;
		}

		candidates.push({
			vertex: vertex,
			nodeTemplateName: candidateNodeTemplateName,
			capabilities: candidateCapabilities
		});
	}

	return candidates;
}

function gatherCandidateCapabilities(requirement, candidateNodeTemplates) {
	capabilityName = requirement.capabilityName;
	capabilityTypeName = requirement.capabilityTypeName;

	candidates = [];
	for (var c = 0; c < candidateNodeTemplates.length; c++) {
		candidate = candidateNodeTemplates[c];
		candidateVertex = candidate.vertex;
		candidateNodeTemplateName = candidate.nodeTemplateName;

		candidateCapabilities = [];
		for (var candidateCapabilityName in candidate.capabilities) {
			candidateCapabilities.push({
				name: candidateCapabilityName,
				capability: candidate.capabilities[candidateCapabilityName]
			});
		}

		// For consistent results, we will sort the candidate capabilities by name
		candidateCapabilities.sort(function(a, b) {
			return a.name < b.name ? -1 : 1;
		});

		for (var cc = 0; cc < candidateCapabilities.length; cc++) {
			candidateCapabilityName = candidateCapabilities[cc].name;

			if ((capabilityName !== '') && (capabilityName !== candidateCapabilityName)) {
				puccini.log.debugf('{resolve} %s: capability "%s" in node template "%s" is not named "%s"', path, candidateCapabilityName, candidateNodeTemplateName, capabilityName);
				continue;
			}

			candidateCapability = candidateCapabilities[cc].capability;

			if ((capabilityTypeName !== '') && !(capabilityTypeName in candidateCapability.types)) {
				puccini.log.debugf('{resolve} %s: capability "%s" in node template "%s" is not of type "%s"', path, candidateCapabilityName, candidateNodeTemplateName, capabilityTypeName);
				continue;
			}

			maxRelationshipCount = candidateCapability.maxRelationshipCount;
			if ((maxRelationshipCount !== -1) && (countRelationships(candidateVertex, candidateCapabilityName) === maxRelationshipCount)) {
				puccini.log.debugf('{resolve} %s: capability "%s" in node template "%s" already has %d relationships, the maximum allowed', path, candidateCapabilityName, candidateNodeTemplateName, maxRelationshipCount);
				continue;
			}

			candidates.push({
				vertex: candidateVertex,
				nodeTemplateName: candidateNodeTemplateName,
				capability: candidateCapability,
				capabilityName: candidateCapabilityName
			});
		}
	}

	return candidates;
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

function countRelationships(vertex, capabilityName) {
	count = 0;
	for (var e = 0; e < vertex.edgesIn.length; e++) {
		edge = vertex.edgesIn[e];
		if (tosca.isTosca(edge, 'relationship') && (edge.properties.capability === capabilityName))
			count++;
	}
	return count;
}

function arePropertiesValid(path, kind, name, entity, constraintsMap) {
	valid = true;

	properties = entity.properties;
	for (var propertyName in constraintsMap) {
		puccini.log.debugf('{resolve} %s: applying constraints to property "%s" of %s "%s"', path, propertyName, kind, name);

		property = properties[propertyName];
		if (property === undefined) {
			// return false; GOJA: returning from inside for-loop is broken
			valid = false;
			break;
		}
		property = clout.newCoercible(property, null, null, null)

		constraints = constraintsMap[propertyName];
		constraints = clout.newConstraints(constraints, null, null, null)
		if (!constraints.validate(property)) {
			// return false; GOJA: returning from inside for-loop is broken
			valid = false;
			break;
		}
	}

	return valid;
}

function isSubstituted(nodeTemplateName, requirementName) {
	for (var vertexId in clout.vertexes) {
		vertex = clout.vertexes[vertexId];
		if (tosca.isTosca(vertex, 'substitution')) {
			for (var e = 0; e < vertex.edgesOut.length; e++) {
				edge = vertex.edgesOut[e];
				if (!tosca.isTosca(edge, 'requirementMapping'))
					continue;

				if ((edge.target.properties.name === nodeTemplateName) && (edge.properties.requirement === requirementName))
					return true;
			}

			// There's only ever one substitution
			return false;
		}
	}

	return false;
}

function isMaxCountGreater(a, b) {
	if (a == -1)
		return b !== -1;
	else if (b == -1)
		return false;
	return a > b;
}

function unsatisfied(path, name, message) {
	if (typeof problems === 'undefined')
		puccini.log.infof('%s: could not satisfy "%s" because %s', path, name, message);
	else
		problems.reportf('%s: could not satisfy "%s" because %s', path, name, message);
}

function notEnoughRelationships(nodeTemplateName, capabilityName, relationshipCount, minRelationshipCount) {
	if (typeof problems === 'undefined')
		puccini.log.infof('capability "%s" of node template "%s" does not have enough relationships: %d < %d', capabilityName, nodeTemplateName, relationshipCount, minRelationshipCount);
	else
		problems.reportf('capability "%s" of node template "%s" does not have enough relationships: %d < %d', capabilityName, nodeTemplateName, relationshipCount, minRelationshipCount);
}
