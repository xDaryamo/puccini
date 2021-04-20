TOSCA OpenStack Profile Examples
================================

For [OpenStack](https://www.openstack.org/) Puccini can generate
[Ansible](https://www.ansible.com/) playbooks that rely on the
[Ansible OpenStack collection](https://docs.ansible.com/ansible/latest/collections/openstack/cloud/index.html).
Custom operation artifacts, if included, are deployed to the virtual machines and executed.

Effectively, the combination of TOSCA + Ansible provides an equivalent set of features to
[Heat](https://docs.openstack.org/heat/wallaby/)/[Mistral](https://docs.openstack.org/mistral/wallaby/).
However, Ansible is a general-purpose orchestrator that can do a lot more than Heat. The generated
playbooks comprise roles that can be imported and used in other playbooks, allowing for custom
orchestration integrations.

Note that though Puccini can compile HOT directly, we recommend TOSCA because of its much
richer grammar and features. See the [HOT examples](../hot/).

* [Hello World](hello-world.yaml)

If you have [Ansible](https://www.ansible.com/) installed and configured then you can run something
like this to deploy the example: 

    puccini-tosca compile examples/openstack/hello-world.yaml --exec=openstack.generate --output=dist/openstack
    cd dist/openstack
    ansible-playbook install.yaml

When run for the first time it will provision keys for your deployment. The public and private keys
will be under the `keys` subdirectory. Note that the private key cannot be retrieved after creation,
so make sure not to lose it!

You can now use the private key to login to servers, e.g.:

    ssh -i keys/topology -o StrictHostKeyChecking=no -o UserKnownHostsFile=/dev/null root@192.237.176.164

(Note the `ssh` options to avoid storing certificates for the IP address. It's good practice because
IP addresses in the cloud may be reused.)

You can run the playbook multiple times. If the servers are already running, they will *not* be
recreated.


Installing Ansible
------------------

Many operating systems have Ansible as a package, but you can install a specific version manually
in a Python virtual environment. [Here's](scripts/install-ansible) our example script.


Configuring for Your OpenStack
------------------------------

The `openstack.generate` scriptlet will generate a template `clouds.yaml` skeleton for you if the
file does not exist. You will need to edit it with the proper credentials for accessing your
OpenStack instance.
See the [documentation](https://docs.openstack.org/python-openstackclient/pike/configuration/).


Testing with Rackspace
----------------------

[Rackspace](https://www.rackspace.com/) provides a public OpenStack cloud.

You will also need to install Rackspace's authentication plugin:

    pip install rackspaceauth

Edit your `clouds.yaml` to look something like this:

    clouds:
      rackspace:
        region_name: ORD
        auth_type: rackspace_apikey
        auth:
          username: USERNAME
          api_key: API_KEY
          auth_url: https://identity.api.rackspacecloud.com/v2.0/

Rackspace uses non-standard image and flavor names, so you will need to provide inputs to change
the defaults:

    puccini-tosca compile examples/openstack/hello-world.yaml -i image_id="CentOS 7 (PVHVM)" -i flavor="512MB Standard Instance"
