#!/usr/bin/env python

# Note that installing `puccini` will also install `ruamel.yaml` 

import sys, puccini.tosca
from ruamel.yaml import YAML

yaml = YAML()

if len(sys.argv) <= 1:
    sys.stderr.write('no URL provided\n')
    sys.exit(1)

url = sys.argv[1]

clout = puccini.tosca.compile(url)

yaml.dump(clout, sys.stdout)
