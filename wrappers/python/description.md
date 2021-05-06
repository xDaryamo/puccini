Puccini
=======

Parse and compile [TOSCA](https://www.oasis-open.org/committees/tosca/)
to [Clout](https://puccini.cloud/clout/).

Part of the [Puccini](https://puccini.cloud) project.

This is a work in progress.


Installation
------------

At this time only Linux x86_64 platforms are supported. We will update this page
as more platforms are added. In most cases this should work:

    pip install puccini

If you are using Python 3.9 on most Linuxes then you will get our pre-built binaries
when installing.

For other environments, Puccini will be built from source. Though Puccini is written
in Go, you do not need Go tooling installed, as it will be downloaded on-demand by
our installer. The only requirement is that your operating system have the `curl`
and `tar` tools.


Usage
-----

Example:

    import sys, puccini.tosca, ard

    try:
        clout = puccini.tosca.compile('/path/to/my-tosca-service.csar') # can also be a URL
        ard.write(clout, sys.stdout)
    except puccini.tosca.Problems as e:
        print('Problems:', file=sys.stderr)
        for problem in e.problems:
            ard.write(problem, sys.stderr)
