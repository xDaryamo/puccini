#!/usr/bin/env python

import sys, puccini_tosca
from ruamel.yaml import YAML

yaml = YAML()

url = sys.argv[1]
clout = puccini_tosca.compile(url)
yaml.dump(clout, sys.stdout)
