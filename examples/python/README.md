Python Example
==============

First, build the Puccini shared library:

    scripts/build-libraries.sh

Then, install the excellent [ruamel.yaml](https://yaml.readthedocs.io/en/latest/) parser:

    pip install ruamel.yaml

You can now run the example:

    examples/python/compile.py examples/tosca/data-types.yaml
