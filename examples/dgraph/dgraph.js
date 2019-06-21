
// See: https://docs.dgraph.io/mutations/#json-mutation-format

clout.exec('tosca.utils');

tosca.coerce();

var vertexItems = [];
var cloutItem = {'clout:vertex': vertexItems};
var items = [cloutItem];

for (var vertexId in clout.vertexes) {
	var vertex = clout.vertexes[vertexId];

	var vertexItem = {uid: '_:clout.vertex.' + vertexId, 'clout:edge': []};

	if (tosca.isTosca(vertex, 'nodeTemplate'))
		fillNodeTemplate(vertexItem, vertex.properties);

	for (var e in vertex.edgesOut)
		fillEdge(vertexItem, vertex.edgesOut[e]);

	vertexItems.push(vertexItem);
}

function fillEdge(item, edge) {
	var edgeItem = {uid: '_:clout.vertex.' + edge.targetID};

	if (tosca.isTosca(edge, 'relationship'))
		fillRelationship(edgeItem, edge.properties);

	item['clout:edge'].push(edgeItem);
}

function fillTosca(item, entity, type_, prefix) {
	if (prefix === undefined)
		prefix = '';
	item[prefix + 'tosca:entity'] = type_;
	item[prefix + 'tosca:name'] = entity.name;
	item[prefix + 'tosca:description'] = entity.description;
	item[prefix + 'tosca:types'] = JSON.stringify(entity.types);
	item[prefix + 'tosca:properties'] = JSON.stringify(entity.properties);
	item[prefix + 'tosca:attributes'] = JSON.stringify(entity.attributes);
}

function fillNodeTemplate(item, nodeTemplate) {
	fillTosca(item, nodeTemplate, 'nodeTemplate');

	item.capabilities = [];
	for (var name in nodeTemplate.capabilities) {
		var capability = nodeTemplate.capabilities[name];
		var capabilityItem = {};
		fillTosca(capabilityItem, capability, 'capability');
		item.capabilities.push(capabilityItem);
	}
}

function fillRelationship(item, relationship) {
	// As facets
	fillTosca(item, relationship, 'relationship', 'clout:edge|');
}

puccini.format = 'json';
puccini.write({set: items});
