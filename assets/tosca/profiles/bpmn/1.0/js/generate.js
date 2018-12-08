clout.exec('tosca.utils');

tosca.coerce();

bpmn = puccini.newXmlDocument();

bpmn.createProcInst('xml', 'version="1.0" encoding="UTF-8"');

definitions = bpmn.createElement('bpmn:definitions');
definitions.createAttr('id', 'definitions');
definitions.createAttr('xmlns:bpmn', 'http://www.omg.org/spec/BPMN/20100524/MODEL');
definitions.createAttr('xmlns:bpmndi', 'http://www.omg.org/spec/BPMN/20100524/DI');
definitions.createAttr('xmlns:di', 'http://www.omg.org/spec/DD/20100524/DI');
definitions.createAttr('xmlns:dc', 'http://www.omg.org/spec/DD/20100524/DC');
definitions.createAttr('typeLanguage', 'http://www.java.com/javaTypes');
definitions.createAttr('expressionLanguage', 'http://www.mvel.org/2.0');
definitions.createAttr('targetNamespace', 'http://bpmn.io/schema/bpmn');
definitions.createAttr('exporter', 'puccini');

for (id in clout.vertexes) {
	vertex = clout.vertexes[id];

	if (tosca.isTosca(vertex, 'policy') && ('bpmn.Process' in vertex.properties.types))
		createPolicyProcess(id, vertex);

	if (tosca.isTosca(vertex, 'workflow'))
		createWorkflowProcess(id, vertex);
}

puccini.write(bpmn);

function createPolicyProcess(id, vertex) {
	policy = vertex.properties;

	process = createProcess(id, policy.name + ' policy');

	tasks = [];
	for (e = 0; e < vertex.edgesOut.length; e++) {
		edge = vertex.edgesOut[e];
		if (!tosca.isTosca(edge, 'nodeTemplateTarget') && !tosca.isTosca(edge, 'groupTarget'))
			continue;
		target = edge.target.properties;

		// Iterate edges
		for (ee = 0; ee < vertex.edgesOut.length; ee++) {
			edge = vertex.edgesOut[ee];
			if (tosca.isTosca(edge, 'policyTriggerOperation')) {
				task = createPolicyTriggerOperationTask(target, edge.target.properties);
				tasks.push(task);
			} else if (tosca.isTosca(edge, 'policyTriggerWorkflow')) {
				// TODO
			}
		}
	}

	startGateway = startEvent = createEvent(process);
	endEvent = createEvent(process, true);

	code = puccini.sprintf('\nstartProcess("%s");\n', policy.properties.bpmn_process_id);
	endGateway = endTask = createScriptTask(clout.newKey(), 'startProcess', code);
	createSequenceFlow(process, endTask, endEvent);

	if (tasks.length > 1) {
		startGateway = createGateway(process);
		endGateway = createGateway(process, true);
		createSequenceFlow(process, startEvent, startGateway);
		createSequenceFlow(process, endGateway, endTask);
	}

	for (t in tasks) {
		task = tasks[t];
		createSequenceFlow(process, startGateway, task);
		createSequenceFlow(process, task, endGateway);
	}
}

function createPolicyTriggerOperationTask(target, operation) {
	// TODO: handle inputs and dependencies
	code = puccini.sprintf('\ncallOperation("%s", "%s");\n', target.name, operation.implementation);
	task = createScriptTask(clout.newKey(), 'operation on ' + target.name, code);
	return task
}

function createWorkflowProcess(id, vertex) {
	workflow = vertex.properties;

	process = createProcess(id, workflow.name + ' workflow');

	// Iterate steps
	tasks = {};
	for (e = 0; e < vertex.edgesOut.length; e++) {
		edge = vertex.edgesOut[e];
		if (!tosca.isTosca(edge, 'workflowStep'))
			continue;

		step = edge.target;
		stepID = edge.targetID;

		createWorkflowTask(step, stepID, tasks);
	}

	// Link previous tasks in graph
	for (name in tasks) {
		task = tasks[name];
		for (n in task.next) {
			nextName = task.next[n];
			tasks[nextName].prev.push(name);
		}
	}

	// Count first tasks
	first = 0;
	for (t in tasks) {
		task = tasks[t];
		if (task.prev.length === 0)
			first++;
	}

	// Count last tasks
	last = 0;
	for (t in tasks) {
		task = tasks[t];
		if (task.next.length === 0)
			last++;
	}

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
	for (t in tasks) {
		task = tasks[t];
		if (task.prev.length === 0)
			createSequenceFlow(process, startGateway, task.task);
		else if (task.prev.length > 1)
			task.incoming = createGateway(process, true);
	}

	// Outgoing
	for (t in tasks) {
		task = tasks[t];
		if (task.next.length === 0)
			createSequenceFlow(process, task.task, endGateway);
		else if (task.next.length > 1)
			task.outgoing = createGateway(process);
	}

	// Link incoming to outgoing
	for (t in tasks) {
		task = tasks[t];

		if (getAttr(task.task, 'id') != getAttr(task.incoming, 'id'))
			createSequenceFlow(process, task.incoming, task.task);

		if (getAttr(task.task, 'id') != getAttr(task.outgoing, 'id'))
			createSequenceFlow(process, task.task, task.outgoing);

		for (n in task.next) {
			next = tasks[task.next[n]];
			createSequenceFlow(process, task.outgoing, next.incoming);
		}
	}
}

function createWorkflowTask(step, stepID, tasks) {
	name = step.properties.name;

	next = [];

	code = '\nnodeTemplates = [];\ngroups = [];';

	// Iterate edges
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
		} else if (tosca.isTosca(edge, 'onSuccess')) {
			next.push(edge.target.properties.name);
		} else if (tosca.isTosca(edge, 'onFailure')) {
			next.push(edge.target.properties.name);
		}
	}

	// Iterate activities
	for (a in activities) {
		activity = activities[a];
		if (activity.setNodeState)
			code += puccini.sprintf('\nsetNodeState(nodeTemplates, groups, "%s");', activity.setNodeState);
		else if (activity.callOperation)
			code += puccini.sprintf('\ncallOperation(nodeTemplates, groups, "%s", "%s");', activity.callOperation.interface, activity.callOperation.operation);
	}

	task = createScriptTask(stepID, name + ' step', code + '\n');

	tasks[name] = {
		task: task,
		incoming: task,
		outgoing: task,
		prev: [],
		next: next
	};
}

function createProcess(id, name) {
	process = definitions.createElement('bpmn:process');
	process.createAttr('id', id);
	process.createAttr('name', name);
	process.createAttr('isExecutable', 'true');
	return process;
}

function createEvent(process, end) {
	event = process.createElement(end ? 'bpmn:endEvent' : 'bpmn:startEvent');
	event.createAttr('id', clout.newKey());
	event.createAttr('name', end ? 'end' : 'start');
	if (end)
		event.createElement('bpmn:terminateEventDefinition');
	return event;
}

function createScriptTask(id, name, code) {
	task = process.createElement('bpmn:scriptTask');
	task.createAttr('id', id);
	task.createAttr('name', name);
	task.createAttr('scriptFormat', 'javascript');
	script = task.createElement('bpmn:script');
	script.setText(code);
	return task;
}

function createGateway(process, converging) {
	gw = process.createElement('bpmn:parallelGateway');
	gw.createAttr('id', clout.newKey());
	gw.createAttr('name', converging ? 'converge' : 'diverge');
	gw.createAttr('gatewayDirection', converging ? 'Converging' : 'Diverging');
	return gw
}

function createSequenceFlow(process, source, target) {
	flow = process.createElement('bpmn:sequenceFlow');
	flow.createAttr('id', clout.newKey());
	flow.createAttr('sourceRef', getAttr(source, 'id'));
	flow.createAttr('targetRef', getAttr(target, 'id'));
}

function getID(element) {
	return getAttr(element, 'id');
}

function getAttr(element, key) {
	r = null;
	for (a in element.attr) {
		attr = element.attr[a];
		if (attr.key === key)
			r = attr.value;
	}
	return r;
}
