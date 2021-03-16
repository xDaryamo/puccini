#!/usr/bin/env python3

# Note that installing `puccini` will also install `ruamel.yaml` 

import sys, puccini.tosca
from ruamel.yaml import YAML

yaml = YAML()

if len(sys.argv) <= 1:
    sys.stderr.write('no URL provided\n')
    sys.exit(1)

url = sys.argv[1]

try:
    clout = puccini.tosca.compile(url)
    yaml.dump(clout, sys.stdout)
except puccini.tosca.Problems as e:
    print('Problems:', file=sys.stderr)
    for problem in e.problems:
        yaml.dump(problem, sys.stderr)
    sys.exit(1)
