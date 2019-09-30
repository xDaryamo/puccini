
clout.exec('tosca.utils');

tosca.coerce();

var jQueryVersion = '2.2.4';
var jQueryUiVersion = '1.12.1';
var jQueryLayoutVersion = '1.4.3';
var visJsVersion = '4.21.0';

var colorNode = 'rgb(100,200,255)';
var colorPolicy = 'rgb(255,165,0)';
var colorSubstitution = 'rgb(150,200,255)';
var colorWorkflow = 'rgb(100,255,100)';

var nodes = [];
var edges = [];

if (tosca.isTosca(clout)) {
	var templateName = clout.properties.tosca.metadata.template_name;
	var templateAuthor = clout.properties.tosca.metadata.template_author;
	var templateVersion = clout.properties.tosca.metadata.template_version;
	var description = clout.properties.tosca.description;

	var header = '<h1>Clout from TOSCA Service Template</h1>';
	if (templateName)
		header += '<h2>' + escapeHtml(templateName) + '</h2>';
	if (templateVersion)
		header += '<p><b>Version</b>: ' + escapeHtml(templateVersion) + '</p>';
	if (templateAuthor)
		header += '<p><b>Author</b>: ' + escapeHtml(templateAuthor) + '</p>';
	if (description)
		header += '<p>' + formatDescription(description) + '</p>';
} else {
	header = '<h1>Clout</h1>';
}

for (var id in clout.vertexes) {
	var vertex = clout.vertexes[id];
	addVertex(id, vertex);
}

function formatDescription(description) {
	var r = '';
	var paragraphs = description.split('\n');
	for (var p in paragraphs) {
		var paragraph = paragraphs[p];
		if (paragraph)
			r += '<p>' + escapeHtml(paragraph) + '</p>';
	}
	return r;
}

function escapeHtml(unsafe) {
	// See: https://stackoverflow.com/a/6234804/849021
	return unsafe
		.replace(/&/g, "&amp;")
		.replace(/</g, "&lt;")
		.replace(/>/g, "&gt;")
		.replace(/"/g, "&quot;")
		.replace(/'/g, "&#039;");
}

function jsonify(data) {
	return JSON.stringify(data, null, '\t').replace(/^/mg, '\t').substr(1);
}

function addVertex(id, vertex) {
	var node = {
		id: id,
		label: id,
		data: tosca.isTosca(vertex) ? vertex.properties : vertex
	};

	if (tosca.isTosca(vertex, 'nodeTemplate'))
		addNodeTemplate(node);
	else if (tosca.isTosca(vertex, 'group'))
		addGroup(node);
	else if (tosca.isTosca(vertex, 'policy'))
		addPolicy(node);
	else if (tosca.isTosca(vertex, 'substitution'))
		addSubstitution(node);
	else if (tosca.isTosca(vertex, 'workflow'))
		addWorkflow(node);
	else if (tosca.isTosca(vertex, 'workflowStep'))
		addWorkflowStep(node);
	else if (tosca.isTosca(vertex, 'workflowActivity'))
		addWorkflowActivity(node);
	else
		node.data = vertex;

	nodes.push(node);

	for (var e in vertex.edgesOut)
		addEdge(id, vertex.edgesOut[e]);
}

function addEdge(id, e) {
	var edge = {
		from: id,
		to: e.targetID,
		arrows: {
			to: true
		},
		font: {
			align: 'middle'
		},
		smooth: {type: 'dynamic'},
		length: 300,
		data: tosca.isTosca(e) ? e.properties : e
	};

	if (tosca.isTosca(e, 'relationship'))
		addRelationship(edge);
	else if (tosca.isTosca(e, 'requirementMapping'))
		addRequirementMapping(edge);
	else if (tosca.isTosca(e, 'capabilityMapping'))
		addCapabilityMapping(edge);
	else if (tosca.isTosca(e, 'propertyMapping'))
		addPropertyMapping(edge);
	else if (tosca.isTosca(e, 'interfaceMapping'))
		addInterfaceMapping(edge);
	else if (tosca.isTosca(e, 'onFailure'))
		addOnFailure(edge);
	else
		edge.data = e;

	edges.push(edge);
}

function addNodeTemplate(node) {
	node.label = node.data.name;
	node.shape = 'box';
	node.color = colorNode;
}

function addGroup(node) {
	node.label = node.data.name;
	node.shape = 'circle';
	node.color = colorNode;
}

function addPolicy(node) {
	node.label = node.data.name;
	node.shape = 'circle';
	node.color = colorPolicy;
}

function addRelationship(edge) {
	edge.label = edge.data.name;
	edge.color = {color: colorNode};
}

function addSubstitution(node) {
	node.label = 'substitution';
	node.shape = 'box';
	node.color = colorSubstitution;
	node.shapeProperties = {borderDashes: true};
}

function addRequirementMapping(edge) {
	edge.label = 'requirement\n' + edge.data.requirement;
	edge.color = {color: colorSubstitution};
	edge.dashes = true;
}

function addCapabilityMapping(edge) {
	edge.label = 'capability\n' + edge.data.capability;
	edge.color = {color: colorSubstitution};
	edge.dashes = true;
}

function addPropertyMapping(edge) {
	edge.label = 'property\n' + edge.data.property;
	edge.color = {color: colorSubstitution};
	edge.dashes = true;
}

function addInterfaceMapping(edge) {
	edge.label = 'interface\n' + edge.data.interface;
	edge.color = {color: colorSubstitution};
	edge.dashes = true;
}

function addWorkflow(node) {
	node.label = node.data.name;
	node.shape = 'diamond';
	node.color = colorWorkflow;
}

function addWorkflowStep(node) {
	node.label = node.data.name;
	node.shape ='diamond';
	node.color = colorWorkflow;
}

function addOnFailure(edge) {
	edge.label = 'onFailure';
	edge.color = {color: 'rgb(255,100,100)'};
}

function addWorkflowActivity(node) {
	node.label = node.data.name;
	node.shape = 'triangle';
	node.color = colorWorkflow;
}

var template = '\
<!doctype html>\n\
<html>\n\
<head>\n\
	<title>Clout from TOSCA Service Template</title>\n\
	<meta charset="utf-8"/>\n\
	<script src="https://cdnjs.cloudflare.com/ajax/libs/jquery/%s/jquery.min.js"></script>\n\
	<script src="https://cdnjs.cloudflare.com/ajax/libs/jqueryui/%s/jquery-ui.min.js"></script>\n\
	<link href="https://cdnjs.cloudflare.com/ajax/libs/jqueryui/%s/jquery-ui.min.css" rel="stylesheet" type="text/css" />\n\
	<script src="https://cdnjs.cloudflare.com/ajax/libs/jquery-layout/%s/jquery.layout.min.js"></script>\n\
	<link href="https://cdnjs.cloudflare.com/ajax/libs/jquery-layout/%s/layout-default.min.css" rel="stylesheet" type="text/css" />\n\
	<script type="text/javascript" src="https://cdnjs.cloudflare.com/ajax/libs/vis/%s/vis.min.js"></script>\n\
	<link href="https://cdnjs.cloudflare.com/ajax/libs/vis/%s/vis.min.css" rel="stylesheet" type="text/css" />\n\
	<link href="https://fonts.googleapis.com/css?family=Roboto" rel="stylesheet" />\n\
	<link href="https://fonts.googleapis.com/css?family=Roboto+Mono" rel="stylesheet" />\n\
	<style type="text/css">\n\
		body {\n\
			font-family: \'Roboto\', sans-serif;\n\
		}\n\
		.ui-layout-pane {\n\
			background-color: #E9E9E9 !important;\n\
		}\n\
		h1, h2 {\n\
			margin-top: 0;\n\
			margin-bottom: 0;\n\
		}\n\
		p {\n\
			margin-top: 6px;\n\
			margin-bottom: 6px;\n\
		}\n\
		#network {\n\
			width: 100%%;\n\
			height: 100%%;\n\
		}\n\
		#info {\n\
			font-family: \'Roboto Mono\', monospace;\n\
			white-space: pre-wrap;\n\
		}\n\
		#corner {\n\
			float: right;\n\
			font-size: small;\n\
			text-align: right;\n\
		}\n\
	</style>\n\
	<script type="text/javascript">\n\
$(document).ready(function () {\n\
	$(\'body\').layout({\n\
		applyDefaultStyles: true,\n\
		east__size: \'25%%\',\n\
	});\n\
	var nodes = new vis.DataSet(%s);\n\
	var edges = new vis.DataSet(%s);\n\
	var network = new vis.Network(\n\
		document.getElementById(\'network\'),\n\
		{\n\
			nodes: nodes,\n\
			edges: edges\n\
		},\n\
		{\n\
			layout: {\n\
				randomSeed: 1,\n\
				improvedLayout: true\n\
			}\n\
		}\n\
	);\n\
	network.on("click", function (params) {\n\
		if (params.nodes.length === 1) {\n\
			node = nodes.get(params.nodes[0]).data;\n\
			$(\'#info\').text(JSON.stringify(node, null, 4));\n\
		} else if (params.edges.length === 1) {\n\
			edge = edges.get(params.edges[0]).data;\n\
			$(\'#info\').text(JSON.stringify(edge, null, 4));\n\
		}\n\
	});\n\
});\n\
	</script>\n\
</head>\n\
<body>\n\
	<div class="ui-layout-north">\n\
		<div id="corner">Generated by <a href="https://github.com/tliron/puccini">Puccini</a></div>\n\
		%s\n\
	</div>\n\
	<div class="ui-layout-center">\n\
		<div id="network"></div>\n\
	</div>\n\
	<div class="ui-layout-east">\n\
		<div id="info"></div>\n\
	</div>\n\
</body>\n\
</html>';

var html = puccini.sprintf(
	template,
	jQueryVersion,
	jQueryUiVersion, jQueryUiVersion,
	jQueryLayoutVersion, jQueryLayoutVersion,
	visJsVersion, visJsVersion,
	jsonify(nodes),
	jsonify(edges),
	header
);

puccini.write(html);
