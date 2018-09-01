// This file was auto-generated from YAML files

package v1_0

func init() {
	Profile["/tosca/openstack/1.0/js/generate.js"] = `

clout.exec('tosca.utils');

tosca.coerce();

playbook = [];

provision = {
	hosts: 'localhost',
	gather_facts: false,
	tasks: [{
		name: 'Provision servers',
		async: 300, // 5 minutes
		register: 'servers_async',
		with_items: '{{ topology.servers }}',
		os_server: {
		    state: 'present',
		    name: '{{ item.type }}-{{ item.index }}.{{ topology.site_name }}.{{ topology.zone }}',
		    image: '{{ topology.image }}',
		    flavor: '{{ item.flavor }}',
		    key_name: '{{ keypair.key.name }}'
		}
	}]
};

playbook.push(provision);

for (v in clout.vertexes) {
	vertex = clout.vertexes[v];
	if (!tosca.isNodeTemplate(vertex, 'openstack.Nova.Server'))
		continue;
	nodeTemplate = vertex.properties;

	//provision.tasks.push();
}

puccini.write(playbook, 'main.' + puccini.format);
`
}
