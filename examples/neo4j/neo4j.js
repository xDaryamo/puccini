
// See: https://neo4j.com/docs/http-api/3.5/

clout.exec('tosca.utils');

tosca.coerce();

let statements = [];

for (let vertexId in clout.vertexes) {
	let vertex = clout.vertexes[vertexId];

	if (tosca.isTosca(vertex, 'nodeTemplate'))
		createNodeTemplate(vertexId, vertex.properties);
}

for (let vertexId in clout.vertexes) {
	let vertex = clout.vertexes[vertexId];

	for (let e in vertex.edgesOut)
		createEdge(vertexId, vertex.edgesOut[e]);
}

function createEdge(id, edge) {
	if (tosca.isTosca(edge, 'relationship')) 
		createRelationship(id, edge.targetID, edge.properties);
}

function addToscaProperties(entity, properties) {
	properties.name = entity.name;
	properties.description = entity.description;
	properties.type = entity.types;
	properties.property = entity.properties;
	properties.attribute = entity.attributes;
}

function createNodeTemplate(id, nodeTemplate) {
	let properties = {id: id};
	addToscaProperties(nodeTemplate, properties);
	addStatementf('CREATE (n:Clout:TOSCA:NodeTemplate { %s })', formatProperties(properties));

	for (let name in nodeTemplate.capabilities) {
		let capability = nodeTemplate.capabilities[name];
		let capabilityId = id + '.capability.' + name;
		properties = {id: capabilityId};
		addToscaProperties(capability, properties);
		properties.name = name;
		addStatementf('CREATE (n:Clout:TOSCA:Capability { %s })', formatProperties(properties));
		relate('TOSCA:NodeTemplate', id, 'TOSCA:Capability', capabilityId, 'TOSCA_CAPABILITY', properties);
	}
}

function createRelationship(sourceId, targetId, relationship) {
	let properties = {};
	addToscaProperties(relationship, properties);
	relate('TOSCA:NodeTemplate', sourceId, 'TOSCA:NodeTemplate', targetId, 'TOSCA_RELATIONSHIP', properties);
}

function relate(sourceLabel, sourceId, targetLabel, targetId, label, properties) {
	addStatementf('MATCH (s:%s { id: %s }), (t:%s { id: %s }) CREATE (s)-[r:%s { %s }]->(t)', sourceLabel, JSON.stringify(sourceId), targetLabel, JSON.stringify(targetId), label, formatProperties(properties));
}

function formatProperties(properties) {
	let values = {};
	addValues(values, properties);
	let r = [];
	for (let name in values)
		r.push(backtick(name) + ': ' + JSON.stringify(String(values[name])));
	return r.join(', ');
}

function addValues(values, properties, prefix) {
	if (prefix === undefined)
		prefix = '';
	for (let name in properties) {
		let value = properties[name];
		switch (typeof value) {
		case 'object':
			addValues(values, value, prefix + name + ':');
			break;
		default:
			values[prefix + name] = value;
		}
	}
}

function backtick(v) {
	// It looks like there is no escape sequence for backticks
	return '`' + v.replace(/`/g, '_') + '`';
}

function addStatementf() {
	addStatement(puccini.sprintf.apply(null, arguments));
}

function addStatement(statement) {
	statements.push({statement: statement})
}

puccini.format = 'json';
puccini.write({statements: statements});
