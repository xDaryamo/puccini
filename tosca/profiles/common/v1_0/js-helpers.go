// This file was auto-generated from a YAML file

package v1_0

func init() {
	Profile["/tosca/common/1.0/js/helpers.js"] = `

var tosca = {};

tosca.isTosca = function(o, kind) {
	if (o.metadata === undefined)
		return false;
	o = o.metadata['puccini-tosca'];
	if (o === undefined)
		return false;
	if (o.version !== '1.0')
		return false;
	if (kind !== undefined)
		return kind === o.kind;
	return true;
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
	var vertex;
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
		for (var vertexId in clout.vertexes) {
			var vertex = clout.vertexes[vertexId];
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
	for (var e = 0; e < vertex.edgesOut.length; e++) {
		var edge = vertex.edgesOut[e];
		if (tosca.isTosca(edge, 'relationship')) {
			for (var typeName in edge.properties.types) {
				var type = edge.properties.types[typeName];
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
	var c = v.$number;
	if (c !== undefined)
		return c;
	c = v.$string;
	if (c !== undefined)
		return c;
	return v;
};
`
}
