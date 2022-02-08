// This scriptlet gathers all endpoint capabilities and generates a report

const traversal = require('tosca.lib.traversal');
const tosca = require('tosca.lib.utils');

// "traversal.coerce" calls all intrinsic functions and validates all constraints
traversal.coerce();

let endpoints = [];

for (let v in clout.vertexes) {
	let vertex = clout.vertexes[v];

	// We'll skip vertexes that are not TOSCA node templates
	if (!tosca.isNodeTemplate(vertex))
		continue;

	let nodeTemplate = vertex.properties;

	for (let c in nodeTemplate.capabilities) {
		let capability = nodeTemplate.capabilities[c];

		// We'll skip capabilities that do not inherit from Endpoint
		if (!('tosca::Endpoint' in capability.types))
			continue;

		// Add endpoint to the report
		endpoints.push({
			name: nodeTemplate.name + '.' + c,
			protocol: capability.properties.protocol,
			port: capability.properties.port,
		});
	}
}

// "puccini.write" will use either YAML (the default), JSON, or XML according to the format selected
// in the command line (use --format to change it)
puccini.write(endpoints);
