
// Here we're demonstrating how to generate XML from Puccini, which is more complex than YAML/JSON
// Specifically we're exporting a BPMN version of our TOSCA workflows
// (This is not intended as a full BPMN solution, it's a just a demonstration of XML)

// See endpoints.js for general comments about JavaScript

clout.exec('tosca.utils');

tosca.coerce();

bpmn = puccini.newXmlDocument();

// XML tree generation is handled by a library called etree
// Documentation: https://godoc.org/github.com/beevik/etree

bpmn.createProcInst('xml', 'version="1.0" encoding="UTF-8"');

definitions = bpmn.createElement('bpmn:definitions');
definitions.createAttr('xmlns:bpmn', 'http://www.omg.org/spec/BPMN/20100524/MODEL');
definitions.createAttr('xmlns:bpmndi', "http://www.omg.org/spec/BPMN/20100524/DI");
definitions.createAttr('xmlns:di', 'http://www.omg.org/spec/DD/20100524/DI');
definitions.createAttr('xmlns:dc', 'http://www.omg.org/spec/DD/20100524/DC');
definitions.createAttr('targetNamespace', 'http://bpmn.io/schema/bpmn');
definitions.createAttr('exporter', 'puccini');

processes = [];

for (v in clout.vertexes) {
	vertex = clout.vertexes[v];

	// We'll skip vertexes that are not TOSCA workflows
	if (!tosca.isTosca(vertex, 'workflow'))
		continue;

	workflow = vertex.properties;

	process = definitions.createElement('bpmn:process');
	process.createAttr('id', v);
	process.createAttr('isExecutable', 'true');
	process.createAttr('name', workflow.name);

	// Iterate steps
	for (e in vertex.edgesOut) {
		edge = vertex.edgesOut[e];
		if (!tosca.isTosca(edge, 'workflowStep'))
			continue;

		step = edge.target;

		task = process.createElement('bpmn:scriptTask');
		task.createAttr('id', edge.targetID);
		task.createAttr('scriptFormat', 'javascript');

		code = '\nnodeTemplates = [];\ngroups = [];\nactivities = [];';

		// Iterate activities
		activities = [];
		for (ee in step.edgesOut) {
			edge = step.edgesOut[ee];
			if (tosca.isTosca(edge, 'nodeTemplateTarget'))
				code += puccini.sprintf('\nnodeTemplates.push("%s");', edge.target.properties.name);
			else if (tosca.isTosca(edge, 'groupTarget'))
				code += puccini.sprintf('\ngroups.push("%s");', edge.target.properties.name);
			else if (tosca.isTosca(edge, 'workflowActivity')) {
				// Put activities in the right sequence
				sequence = edge.properties.sequence;
				activities[sequence] = edge.target.properties;
			}
		}

		for (a in activities) {
			activity = activities[a];
			if (activity.setNodeState)
				code += puccini.sprintf('\nsetNodeState(nodeTemplates, groups, "%s");', activity.setNodeState);
			else if (activity.callOperation)
				code += puccini.sprintf('\ncallOperation(nodeTemplates, groups, "%s", "%s");', activity.callOperation.interface, activity.callOperation.operation);
		}

		script = task.createElement('script');
		script.setText(code + '\n');
	}
}

bpmn.indent(2);
bpmn.writeTo(puccini.stdout);
