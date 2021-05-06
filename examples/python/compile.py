#!/usr/bin/env python3

# Note that installing `puccini` will also install `ard` 

import sys, puccini.tosca, ard

if len(sys.argv) <= 1:
    sys.stderr.write('no URL provided\n')
    sys.exit(1)

url = sys.argv[1]

try:
    clout = puccini.tosca.compile(url)
    ard.write(clout, sys.stdout)
except puccini.tosca.Problems as e:
    print('Problems:', file=sys.stderr)
    for problem in e.problems:
        ard.write(problem, sys.stderr)
    sys.exit(1)
