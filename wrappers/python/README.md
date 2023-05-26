Python Wrapper for Puccini
==========================

This is a Python library for calling Puccini. It works by using
[ctypes](https://docs.python.org/3/library/ctypes.html) to call a shared library (.so) built from
Puccini's Go code. This is done in-process, so there's no sub-process forking.

Note that we currently only support installation on 64-bit Linux.

The latest stable version is published on [PyPI](https://pypi.org/project/puccini/). To install:

    pip install puccini

To install from source:

    git clone https://github.com/tliron/puccini.git
    cd puccini
    scripts/build-wrapper-python -e

This will create a [virtualenv](https://virtualenv.pypa.io/). To use it:

    . dist/python-env/bin/activate

Also see: [Python examples](../../examples/python/).
