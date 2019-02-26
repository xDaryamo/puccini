TOSCA OpenStack Profile Examples
================================

* [Hello World](hello-world.yaml)

If you have [Ansible](https://www.ansible.com/) installed and configured then you can run something
like this to deploy: 

    puccini-tosca compile examples/openstack/hello-world.yaml | puccini-js exec openstack.generate -o test
    cd test
    ansible-playbook install.yaml

When run for the first time it will provision keys for your deployment. The public and private keys
will be under the `keys` directory. Note that the private key cannot be retrieved after creation,
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
in a Python virtual environment. Here's how to do it on Fedora:

    sudo dnf install python3-virtualenv libselinux-python3
    virtualenv --system-site-packages env
    . env/bin/activate
    pip install ansible==2.7.8 os-client-config==1.31.2

In the above we specify versions that we used for testing, but feel free to omit the versions and
try the latest and greatest.


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

    pip install rackspaceauth==0.8.1

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
