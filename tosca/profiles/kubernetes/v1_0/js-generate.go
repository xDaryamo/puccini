// This file was auto-generated from YAML files

package v1_0

func init() {
	Profile["/tosca/kubernetes/1.0/js/generate.js"] = `

clout.exec('tosca.utils');

tosca.coerce();

specs = [];

for (vertexId in clout.vertexes) {
	vertex = clout.vertexes[vertexId];
	if (!tosca.isNodeTemplate(vertex))
		continue;
	nodeTemplate = vertex.properties;

	// Find metadata
	metadata = {};
	for (capabilityName in nodeTemplate.capabilities) {
		capability = nodeTemplate.capabilities[capabilityName];
		if ('kubernetes.Metadata' in capability.types) {
			metadata = capability.properties;
			break;
		}
	}

	// At least have the "service" label
	if (metadata.labels === undefined)
		metadata.labels = {};
	metadata.labels.service = nodeTemplate.name;

	// Generate specs
	for (capabilityName in nodeTemplate.capabilities) {
		capability = nodeTemplate.capabilities[capabilityName];
		if ('kubernetes.Service' in capability.types)
			generateService(capability, metadata);
		else if ('kubernetes.Deployment' in capability.types)
			generateDeployment(capability, metadata);
	}

	// Run plugins
//	plugins = clout.getPlugins('kubernetes.plugins');
//	for (p in plugins) {
//		plugin = plugins[p];
//		log.debugf('calling plugin: %s', plugin.name);
//		if (plugin.process)
//			entries = plugin.process(clout, vertex, entries);
//	}
}

puccini.write(specs);

function generateService(capability, metadata) {
	spec = {
		apiVersion: 'v1',
		kind: 'Service',
		metadata: metadata,
		spec: {}
	};

	for (propertyName in capability.properties) {
		v = capability.properties[propertyName];
		spec.spec[propertyName] = v;
	}

	// Default selector
	if (spec.spec.selector === undefined)
		spec.spec.selector = metadata.labels;

	specs.push(spec);
}

function generateDeployment(capability, labels) {
	spec = {
		apiVersion: 'apps/v1',
		kind: 'Deployment',
		metadata: metadata,
		spec: {}
	};

	for (propertyName in capability.properties) {
		v = capability.properties[propertyName];
		switch (propertyName) {
		case 'minReadySeconds':
		case 'progressDeadlineSeconds':
			v = convertScalarUnit(v);
			break;
		case 'strategy':
			s = {
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
			s = {};
			for (t in v) {
				vv = v[t];
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
