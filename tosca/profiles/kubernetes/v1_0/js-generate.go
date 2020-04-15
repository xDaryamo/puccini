// This file was auto-generated from a YAML file

package v1_0

func init() {
	Profile["/tosca/kubernetes/1.0/js/generate.js"] = `

clout.exec('tosca.lib.traversal');

tosca.coerce();

var specs = [];

for (var vertexId in clout.vertexes) {
	var vertex = clout.vertexes[vertexId];
	if (!tosca.isNodeTemplate(vertex))
		continue;
	var nodeTemplate = vertex.properties;

	// Find metadata
	var metadata = {};
	for (var capabilityName in nodeTemplate.capabilities) {
		var capability = nodeTemplate.capabilities[capabilityName];
		if ('kubernetes::Metadata' in capability.types) {
			metadata = capability.properties;
			break;
		}
	}

	// At least have the "service" label
	if (metadata.labels === undefined)
		metadata.labels = {};
	metadata.labels.service = nodeTemplate.name;

	// Generate specs
	for (var capabilityName in nodeTemplate.capabilities) {
		var capability = nodeTemplate.capabilities[capabilityName];
		if ('kubernetes::Service' in capability.types)
			generateService(capability, metadata);
		else if ('kubernetes::Deployment' in capability.types)
			generateDeployment(capability, metadata);
	}

	// Run plugins
//	plugins = clout.getPlugins('kubernetes:plugins');
//	for (var p in plugins) {
//		plugin = plugins[p];
//		log.debugf('calling plugin: %s', plugin.name);
//		if (plugin.process)
//			entries = plugin.process(clout, vertex, entries);
//	}
}

puccini.write(specs);

function generateService(capability, metadata) {
	var spec = {
		apiVersion: 'v1',
		kind: 'Service',
		metadata: metadata,
		spec: {}
	};

	for (var propertyName in capability.properties) {
		var v = capability.properties[propertyName];
		spec.spec[propertyName] = v;
	}

	// Default selector
	if (spec.spec.selector === undefined)
		spec.spec.selector = metadata.labels;

	specs.push(spec);
}

function generateDeployment(capability, labels) {
	var spec = {
		apiVersion: 'apps/v1',
		kind: 'Deployment',
		metadata: metadata,
		spec: {}
	};

	for (var propertyName in capability.properties) {
		var v = capability.properties[propertyName];
		switch (propertyName) {
		case 'minReadySeconds':
		case 'progressDeadlineSeconds':
			v = convertScalarUnit(v);
			break;
		case 'strategy':
			var s = {
				type: v.type
			};
			if (v.type === 'RollingUpdate') {
				s.rollingUpdate = {
					maxSurge: convertAmount(v.maxSurge),
					maxUnavailable: convertAmount(v.maxUnavailable)
				};
			}
			v = s;
			break;
		case 'template':
			var s = {};
			for (var t in v) {
				var vv = v[t];
				switch (t) {
				case 'activeDeadlineSeconds':
				case 'terminationGracePeriodSeconds':
					vv = convertScalarUnit(vv);
					break;
				}
				s[t] = vv;
			}
			v = {
				metadata: metadata,
				spec: s
			};
		}
		spec.spec[propertyName] = v;
	}

	// Default selector
	if ((spec.spec.selector.matchExpressions == undefined) && (spec.spec.selector.matchLabels === undefined))
		spec.spec.selector.matchLabels = metadata.labels;

	specs.push(spec);
}

function convertScalarUnit(v) {
	return v.$number;
}

function convertAmount(v) {
	if (v.factor !== undefined)
		return (v.factor * 100) + '%';
	return v.count;
}
`
}
