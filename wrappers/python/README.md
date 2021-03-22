Python Wrapper for Puccini
==========================

This is a Python library for calling Puccini. It works by using
[ctypes](https://docs.python.org/3/library/ctypes.html) to call a shared library (.so) built from
Puccini's Go code. This is done in-process, so there's no sub-process forking.

Note that we currently only support installation on 64-bit Linux.

The latest stable version is published on [PyPI](https://pypi.org/project/puccini/). To install:

    pip install puccini

To install the latest development version from GitHub:

    pip install git+https://github.com/tliron/puccini#subdirectory=wrappers/python

Or, if you've cloned the repository locally: 

    pip install wrappers/python/

For testing it's recommended to install in a [virtualenv](https://virtualenv.pypa.io/), e.g.:

    python -m venv env
    . env/bin/activate
    pip install git+https://github.com/tliron/puccini#subdirectory=wrappers/python

Also see: [Python examples](../../examples/python/).
