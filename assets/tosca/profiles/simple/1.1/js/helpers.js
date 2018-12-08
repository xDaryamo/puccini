
var tosca = {};

tosca.isTosca = function(o, kind) {
	if (o.metadata['puccini-tosca'] !== undefined) {
		o = o.metadata['puccini-tosca'];
		if (o.version === '1.0') {
			if (kind !== undefined)
				return kind === o.kind;
			return true;
		}
	}
	return false;
};

tosca.isNodeTemplate = function(vertex, typeName) {
	if (tosca.isTosca(vertex, 'nodeTemplate')) {
		if (typeName !== undefined)
			return typeName in vertex.properties.types;
		return true;
	}
	return false;
};

tosca.getNodeTemplate = function(entity) {
	switch (entity) {
	case 'SELF':
		vertex = site;
		break;
	case 'SOURCE':
		vertex = source;
		break;
	case 'TARGET':
		vertex = target;
		break;
	case 'HOST':
		vertex = tosca.getHost(site);
		break;
	default:
		for (vertexId in clout.vertexes) {
			vertex = clout.vertexes[vertexId];
			if (tosca.isNodeTemplate(vertex) && (vertex.properties.name === entity))
				return vertex.properties;
		}
		throw puccini.sprintf('node template "%s" not found', entity);
	}
	if (!tosca.isTosca(vertex))
		throw puccini.sprintf('node template "%s" not found', entity);
	return vertex.properties;
};

tosca.getHost = function(vertex) {
	for (e = 0; e < vertex.edgesOut.length; e++) {
		edge = vertex.edgesOut[e];
		if (tosca.isTosca(edge, 'relationship')) {
			for (typeName in edge.properties.types) {
				type = edge.properties.types[typeName];
				if (type.metadata.role === 'host')
					return edge.target;
			}
		}
	}
	return {};
};

tosca.getComparable = function(v) {
	if ((v === undefined) || (v === null))
		return null;
	if (v.$number !== undefined)
		return v.$number;
	if (v.$string !== undefined)
		return v.$string;
	return v;
};
