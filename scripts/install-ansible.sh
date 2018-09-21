#!/bin/bash
set -e

# Requirements (Fedora)
# sudo dnf install python3-virtualenv python3-libselinux

virtualenv --system-site-packages env
. env/bin/activate
pip install ansible==2.7.0rc3 os-client-config==1.31.2 rackspaceauth==0.8.1
