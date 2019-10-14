#!/usr/bin/env python

# Note that installing `puccini` will also install `ruamel.yaml` 

import sys, puccini.tosca
from ruamel.yaml import YAML

yaml = YAML()

url = sys.argv[1]

clout = puccini.tosca.compile(url)

yaml.dump(clout, sys.stdout)
