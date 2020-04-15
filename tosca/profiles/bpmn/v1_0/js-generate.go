// This file was auto-generated from a YAML file

package v1_0

func init() {
	Profile["/tosca/bpmn/1.0/js/generate.js"] = `

clout.exec('tosca.lib.traversal');

tosca.coerce();

var bpmn = puccini.newXMLDocument();

bpmn.createProcInst('xml', 'version="1.0" encoding="UTF-8"');

var definitions = bpmn.createElement('bpmn:definitions');
definitions.createAttr('id', 'definitions');
definitions.createAttr('xmlns:bpmn', 'http://www.omg.org/spec/BPMN/20100524/MODEL');
definitions.createAttr('xmlns:bpmndi', 'http://www.omg.org/spec/BPMN/20100524/DI');
definitions.createAttr('xmlns:di', 'http://www.omg.org/spec/DD/20100524/DI');
definitions.createAttr('xmlns:dc', 'http://www.omg.org/spec/DD/20100524/DC');
definitions.createAttr('typeLanguage', 'http://www.java.com/javaTypes');
definitions.createAttr('expressionLanguage', 'http://www.mvel.org/2.0');
definitions.createAttr('targetNamespace', 'http://bpmn.io/schema/bpmn');
definitions.createAttr('exporter', 'puccini');

for (var id in clout.vertexes) {
	var vertex = clout.vertexes[id];

	if (tosca.isTosca(vertex, 'Policy') && ('bpmn::Process' in vertex.properties.types))
		createPolicyProcess(id, vertex);

	if (tosca.isTosca(vertex, 'Workflow'))
		createWorkflowProcess(id, vertex);
}

puccini.write(bpmn);

function createPolicyProcess(id, vertex) {
	var policy = vertex.properties;

	var process = createProcess(id, policy.name + ' policy');

	var tasks = [];
	for (var e = 0; e < vertex.edgesOut.length; e++) {
		var edge = vertex.edgesOut[e];
		if (!tosca.isTosca(edge, 'NodeTemplateTarget') && !tosca.isTosca(edge, 'GroupTarget'))
			continue;
		var target = edge.target.properties;

		// Iterate edges
		for (var ee = 0; ee < vertex.edgesOut.length; ee++) {
			edge = vertex.edgesOut[ee];
			if (tosca.isTosca(edge, 'PolicyTriggerOperation')) {
				task = createPolicyTriggerOperationTask(process, target, edge.target.properties);
				tasks.push(task);
			} else if (tosca.isTosca(edge, 'PolicyTriggerWorkflow')) {
				// TODO
			}
		}
	}

	var startGateway, startEvent, endGateway, endEvent, endTask;
	var startGateway = startEvent = createEvent(process);
	var endEvent = createEvent(process, true);

	var code = puccini.sprintf('\nstartProcess("%s");\n', policy.properties.bpmn_process_id);
	endGateway = endTask = createScriptTask(process, clout.newKey(), 'startProcess', code);
	createSequenceFlow(process, endTask, endEvent);

	if (tasks.length > 1) {
		startGateway = createGateway(process);
		endGateway = createGateway(process, true);
		createSequenceFlow(process, startEvent, startGateway);
		createSequenceFlow(process, endGateway, endTask);
	}

	for (var t = 0; t < tasks.length; t++) {
		var task = tasks[t];
		createSequenceFlow(process, startGateway, task);
		createSequenceFlow(process, task, endGateway);
	}
}

function createPolicyTriggerOperationTask(process, target, operation) {
	// TODO: handle inputs and dependencies
	var code = puccini.sprintf('\ncallOperation("%s", "%s");\n', target.name, operation.implementation);
	var task = createScriptTask(process, clout.newKey(), 'operation on ' + target.name, code);
	return task;
}

function createWorkflowProcess(id, vertex) {
	var workflow = vertex.properties;

	var process = createProcess(id, workflow.name + ' workflow');

	// Iterate steps
	var tasks = {};
	for (var e = 0; e < vertex.edgesOut.length; e++) {
		var edge = vertex.edgesOut[e];
		if (!tosca.isTosca(edge, 'WorkflowStep'))
			continue;

		var step = edge.target;
		var stepID = edge.targetID;

		createWorkflowTask(process, step, stepID, tasks);
	}

	// Link previous tasks in graph
	for (var name in tasks) {
		var task = tasks[name];
		for (n in task.next) {
			var nextName = task.next[n];
			tasks[nextName].prev.push(name);
		}
	}

	// Count first tasks
	var first = 0;
	for (var name in tasks) {
		var task = tasks[name];
		if (task.prev.length === 0)
			first++;
	}

	// Count last tasks
	var last = 0;
	for (var name in tasks) {
		var task = tasks[name];
		if (task.next.length === 0)
			last++;
	}

	var startGateway, startEvent, endGateway, endEvent;
	startGateway = startEvent = createEvent(process);
	endGateway = endEvent = createEvent(process, true);

	if (first > 1) {
		startGateway = createGateway(process);
		createSequenceFlow(process, startEvent, startGateway);
	}

	if (last > 1) {
		endGateway = createGateway(process, true);
		createSequenceFlow(process, endGateway, endEvent);
	}

	// Incoming
	for (var name in tasks) {
		var task = tasks[name];
		if (task.prev.length === 0)
			createSequenceFlow(process, startGateway, task.task);
		else if (task.prev.length > 1)
			task.incoming = createGateway(process, true);
	}

	// Outgoing
	for (var name in tasks) {
		var task = tasks[name];
		if (task.next.length === 0)
			createSequenceFlow(process, task.task, endGateway);
		else if (task.next.length > 1)
			task.outgoing = createGateway(process);
	}

	// Link incoming to outgoing
	for (var name in tasks) {
		var task = tasks[name];

		if (getAttr(task.task, 'id') != getAttr(task.incoming, 'id'))
			createSequenceFlow(process, task.incoming, task.task);

		if (getAttr(task.task, 'id') != getAttr(task.outgoing, 'id'))
			createSequenceFlow(process, task.task, task.outgoing);

		for (var n = 0; n < task.next.length; n++) {
			var next = tasks[task.next[n]];
			createSequenceFlow(process, task.outgoing, next.incoming);
		}
	}
}

function createWorkflowTask(process, step, stepID, tasks) {
	var name = step.properties.name;

	var next = [];

	var code = '\nvar nodeTemplates = [];\nvar groups = [];';

	// Iterate edges
	var activities = [];
	for (var ee =- 0; ee < step.edgesOut.length; ee++) {
		var edge = step.edgesOut[ee];
		if (tosca.isTosca(edge, 'NodeTemplateTarget'))
			code += puccini.sprintf('\nnodeTemplates.push("%s");', edge.target.properties.name);
		else if (tosca.isTosca(edge, 'GroupTarget'))
			code += puccini.sprintf('\ngroups.push("%s");', edge.target.properties.name);
		else if (tosca.isTosca(edge, 'WorkflowActivity')) {
			// Put activities in the right sequence
			var sequence = edge.properties.sequence;
			activities[sequence] = edge.target.properties;
		} else if (tosca.isTosca(edge, 'OnSuccess')) {
			next.push(edge.target.properties.name);
		} else if (tosca.isTosca(edge, 'OnFailure')) {
			next.push(edge.target.properties.name);
		}
	}

	// Iterate activities
	for (var a = 0; a < activities.length; a++) {
		var activity = activities[a];
		if (activity.setNodeState)
			code += puccini.sprintf('\nsetNodeState(nodeTemplates, groups, "%s");', activity.setNodeState);
		else if (activity.callOperation)
			code += puccini.sprintf('\ncallOperation(nodeTemplates, groups, "%s", "%s");', activity.callOperation.interface, activity.callOperation.operation);
	}

	var task = createScriptTask(process, stepID, name + ' step', code + '\n');

	tasks[name] = {
		task: task,
		incoming: task,
		outgoing: task,
		prev: [],
		next: next
	};
}

function createProcess(id, name) {
	var process = definitions.createElement('bpmn:process');
	process.createAttr('id', id);
	process.createAttr('name', name);
	process.createAttr('isExecutable', 'true');
	return process;
}

function createEvent(process, end) {
	var event = process.createElement(end ? 'bpmn:endEvent' : 'bpmn:startEvent');
	event.createAttr('id', clout.newKey());
	event.createAttr('name', end ? 'end' : 'start');
	if (end)
		event.createElement('bpmn:terminateEventDefinition');
	return event;
}

function createScriptTask(process, id, name, code) {
	var task = process.createElement('bpmn:scriptTask');
	task.createAttr('id', id);
	task.createAttr('name', name);
	task.createAttr('scriptFormat', 'javascript');
	var script = task.createElement('bpmn:script');
	script.setText(code);
	return task;
}

function createGateway(process, converging) {
	var gw = process.createElement('bpmn:parallelGateway');
	gw.createAttr('id', clout.newKey());
	gw.createAttr('name', converging ? 'converge' : 'diverge');
	gw.createAttr('gatewayDirection', converging ? 'Converging' : 'Diverging');
	return gw;
}

function createSequenceFlow(process, source, target) {
	var flow = process.createElement('bpmn:sequenceFlow');
	flow.createAttr('id', clout.newKey());
	flow.createAttr('sourceRef', getAttr(source, 'id'));
	flow.createAttr('targetRef', getAttr(target, 'id'));
}

function getID(element) {
	return getAttr(element, 'id');
}

function getAttr(element, key) {
	var r = null;
	for (var a = 0; a < element.attr.length; a++) {
		var attr = element.attr[a];
		if (attr.key === key)
			r = attr.value;
	}
	return r;
}
`
}
