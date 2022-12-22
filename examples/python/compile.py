#!/usr/bin/env python3

# Note that installing `puccini` will also install `ard` 

import sys, argparse, puccini.tosca, ard

parser = argparse.ArgumentParser(description='Compile TOSCA')
parser.add_argument('url', help='URL to TOSCA file or CSAR')
parser.add_argument('-i', '--input', dest='inputs', nargs='*', action='extend', help='specify input (format is name=value)')
parser.add_argument('-q', '--quirk', dest='quirks', nargs='*', action='extend', help='specify quirk')

args = parser.parse_args()

if args.inputs:
    inputs = {}
    for i in args.inputs:
        k, v = i.split('=')
        inputs[k] = ard.decode_yaml(v)
    args.inputs = inputs

try:
    clout = puccini.tosca.compile(args.url, args.inputs, args.quirks)
    ard.write(clout, sys.stdout)
except puccini.tosca.Problems as e:
    print('Problems:', file=sys.stderr)
    for problem in e.problems:
        ard.write(problem, sys.stderr)
    sys.exit(1)
