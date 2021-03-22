Puccini TOSCA Collection for Ansible
====================================

Enables [TOSCA](https://www.oasis-open.org/committees/tosca/) support for Ansible.

Part of the [Puccini](https://puccini.cloud) project.

This is a work in progress.

Installation
------------

Requires the [Puccini Python library](https://pypi.org/project/puccini/).

Often this should be enough to get up and running:

    pip install puccini
    ansible-galaxy collection install puccini.tosca

Usage
-----

We currently enable two ways to consume TOSCA in Ansible:

* Iterate over arbitrary TOSCA nodes using our custom task `puccini.tosca.compile`
* Use TOSCA nodes as an Ansible inventory, via our custom inventory plugin `puccini.tosca.nodes`

See the [examples](https://github.com/tliron/puccini/tree/main/examples/ansible).
