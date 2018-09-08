TOSCA OpenStack Profile Example
===============================

* [Hello World](hello-world.yaml)

Installing Ansible
------------------

Many operating systems have Ansible as a package, but you can install the latest version manually
in a Python virtual environment. Here's how to do it on Fedora:

    sudo dnf install python3-virtualenv libselinux-python3
    virtualenv --system-site-packages env
    . env/bin/activate
    pip install ansible==2.6.4 os-client-config openstacksdk

If you're deploying to Rackspace you will also need:

    pip install rackspaceauth

