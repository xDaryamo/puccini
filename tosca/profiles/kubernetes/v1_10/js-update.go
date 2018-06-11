// This file was auto-generated from YAML files

package v1_10

func init() {
	Profile["/tosca/kubernetes/1.10/js/update.js"] = `

// TODO

function validateClout() {
	var version;
	if (clout.metadata['puccini-tosca']) {
		version = clout.metadata['puccini-tosca'].version;
	}
	if (!version) {
		throw 'Clout is not TOSCA';
	}
	if (version != '1.0') {
		throw sprintf('unsupported puccini-tosca.version: %s', version);
	}
}

function isKubernetes(vertex) {
	return (vertex.properties.TOSCA !== undefined) && ('kubernetes.Service' in vertex.properties.TOSCA.types);
}

validateClout();

specs = []

for (key in clout.Vertexes) {
	v = clout.Vertexes[key];
	if (!isKubernetes(v)) continue;

	d = v.Properties.TOSCA.Properties.management_port;
	specs.push(d);
	d = clout.coerce(d, v);
	specs.push(d);

	if (!v.Properties.TOSCA || !v.Properties.TOSCA.Types || !('kubernetes.Service' in v.Properties.TOSCA.Types)) {
		continue;
	}

	log.Infof('updating service: %s', key);

	d = v.Properties.TOSCA.Capabilities.service.Properties.port.port;
	specs.push(d);
	d = clout.coerce(d);
	specs.push(d);
}

puccini.write(specs);
`
}
