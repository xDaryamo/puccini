Python Example
==============

The root [`python`](../../python/) directory contains a library for calling Puccini from within
Python. It works by using [ctypes](https://docs.python.org/3/library/ctypes.html) to call a shared
library (.so) built from Puccini's Go code.

Note that we currently only support installation on 64-bit Linux.

To install the latest development version use [pip](https://pip.pypa.io/): 

    pip install git+https://github.com/tliron/puccini#subdirectory=python

Or, if you've cloned the repository locally: 

    pip install python/

You can now run the example:

    examples/python/compile.py examples/tosca/data-types.yaml

For testing it's recommended to install in a [virtualenv](https://virtualenv.pypa.io/):

    virtualenv env
    . env/bin/activate
    pip install git+https://github.com/tliron/puccini#subdirectory=python
