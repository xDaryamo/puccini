Python Wrapper for Puccini
==========================

This is a Python library for calling Puccini. It works by using
[ctypes](https://docs.python.org/3/library/ctypes.html) to call a shared library (.so) built from
Puccini's Go code. This is done in-process, so there's no sub-process forking.

Note that we currently only support installation on 64-bit Linux.

The [`setup.py`](setup.py) is designed to fetch what it needs from the Internet: it will fetch a Go
compiler, fetch the Puccini sources, and then build the shared library. So there are no special
prerequisites to install it.

For example, to install the latest development version via [pip](https://pip.pypa.io/): 

    pip install git+https://github.com/tliron/puccini#subdirectory=wrappers/python

Or, if you've cloned the repository locally: 

    pip install wrappers/python/

For testing it's recommended to install in a [virtualenv](https://virtualenv.pypa.io/):

    virtualenv env
    . env/bin/activate
    pip install git+https://github.com/tliron/puccini#subdirectory=wrappers/python

To build a distribution or an egg:

    cd wrappers/python
    ./setup.py bdist
    ./setup.py bdist_egg

The output will be in the `dist/` subdirectory here.

Also see: [Python examples](../../examples/python/).
