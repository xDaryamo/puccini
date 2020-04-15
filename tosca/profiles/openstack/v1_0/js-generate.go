// This file was auto-generated from a YAML file

package v1_0

func init() {
	Profile["/tosca/openstack/1.0/js/generate.js"] = `

clout.exec('tosca.lib.traversal');

tosca.coerce();

// Vars
writeTopology();

// Playbooks
writePlaybookInstall();

// Roles
writeRoleKeypair();
writeRoleServers();

// Configuration
writeInventory();
writeClouds();
writeCfg();

function writeTopology() {
	var servers = [];

	for (var vertexId in clout.vertexes) {
		var vertex = clout.vertexes[vertexId];
		if (!tosca.isNodeTemplate(vertex, 'openstack::nova.Server'))
			continue;
		var nodeTemplate = vertex.properties;

		servers.push({
			name: nodeTemplate.name,
			image: nodeTemplate.properties.image,
			flavor: nodeTemplate.properties.flavor
		});
	}

	puccini.write({
		servers: servers
	}, 'vars/topology.' + puccini.format);
}

function writePlaybookInstall() {
	puccini.write([[{
		hosts: 'localhost',
		gather_facts: false,
		tasks: [{
			name: 'Configure OpenStack',
			os_client_config: null,
		}, {
			name: 'Import topology.' + puccini.format,
			include_vars: {
				name: 'topology',
				file: 'topology.' + puccini.format
			}
		}, {
			include_role: {
				name: 'openstack-keypair'
			}
		}, {
			include_role: {
				name: 'openstack-servers'
			}
		}]
	}]], 'install.' + puccini.format);
}

function writeRoleKeypair() {
	puccini.write([[{
		name: 'Provision keypair',
		register: 'keypair',
		os_keypair: {
			state: 'present',
			name: 'topology'
		}
	}, {
		name: 'Write keys',
		when: 'keypair.key.private_key is not none', // will only be available when keypair is created
		block: [{
			file: {
				path: '{{ playbook_dir }}/keys',
				state: 'directory'
			}
		}, {
			copy: {
				content: '{{ keypair.key.private_key }}',
				dest: '{{ playbook_dir }}/keys/{{ keypair.key.name }}'
			}
		}, {
			file: {
				mode: 384, // = octal 0600, required by ssh
				dest: '{{ playbook_dir }}/keys/{{ keypair.key.name }}'
			}
		}, {
			copy: {
				content: '{{ keypair.key.public_key }}',
				dest: '{{ playbook_dir }}/keys/{{ keypair.key.name }}.pub'
			}
		}]
	}]], 'roles/openstack-keypair/tasks/main.' + puccini.format);
}

function writeRoleServers() {
	puccini.write([[{
		name: 'Provision servers',
		async: 300, // 5 minutes
		register: 'servers_async',
		with_items: '{{ topology.servers }}',
		os_server: {
		    state: 'present',
		    name: '{{ item.name }}',
		    image: '{{ item.image }}',
		    flavor: '{{ item.flavor }}',
		    key_name: '{{ keypair.key.name }}'
		}
	}, {
		name: 'Wait for servers to become active',
		retries: 300, // delay is 5 seconds
		register: 'servers',
		until: 'servers.finished',
		with_items: '{{ servers_async.results }}',
		async_status: {
			jid: '{{ item.ansible_job_id }}'
		}
	}, {
		name: 'Add servers to group',
		with_items: '{{ servers.results }}',
		add_host: {
			name: '{{ item.server.name }}',
			groups: 'servers',
			// Custom attributes:
			server: '{{ item.server }}',
			// Ansible attributes:
		    ansible_ssh_host: '{{ item.server.public_v4 }}',
		    ansible_ssh_user: 'root',
		    ansible_ssh_private_key_file: '{{ playbook_dir }}/keys/{{ keypair.key.name }}',
		}
	}]], 'roles/openstack-servers/tasks/main.' + puccini.format);
}

function writeInventory() {
	puccini.write({
		all: {}
	}, 'inventory.' + puccini.format);
}

function writeClouds() {
	puccini.write({
		clouds: {
			'default': {
				auth: {
					auth_url: 'AUTH_URL',
					username: 'USERNAME',
					password: 'PASSWORD'
				}
			}
		}
	}, 'clouds.' + puccini.format, true);
}

function writeCfg() {
	puccini.write('\
[defaults]\n\
ansible_managed=Modified by Ansible on %Y-%m-%d %H:%M:%S %Z\n\
inventory=./inventory.' + puccini.format +'\n\
transport=ssh\n\
command_warnings=false\n\
\n\
# Better concurrency\n\
forks=25\n\
\n\
# Required because in cloud IP addresses might be reused\n\
host_key_checking=false\n\
\n\
[ssh_connection]\n\
# Faster SSH\n\
pipelining=true\n\
', 'ansible.cfg');
}
`
}
