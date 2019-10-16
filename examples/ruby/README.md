Ruby Example
============

This relies on the [Ruby wrapper](../../wrappers/ruby/), so make sure to install that first.

To run this example we need to make sure that the Ruby process can load the shared library:

    LD_LIBRARY_PATH=$LD_LIBRARY_PATH:dist examples/ruby/compile.rb examples/tosca/data-types.yaml
