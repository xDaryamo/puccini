Ruby Wrapper for Puccini
========================

This is a Ruby library for calling Puccini. It works by using
[Fiddle](https://ruby-doc.org/stdlib-2.0.0/libdoc/fiddle/rdoc/Fiddle.html) to call a shared library 
(.so) built from Puccini's Go code. This is done in-process, so there's no sub-process forking.

Note that we currently only support installation on 64-bit Linux.

To build the library and gem:

    scripts/build-library.sh
    gem build wrappers/ruby/puccini.gemspec -C wrappers/ruby --output ../../dist/puccini.gem

To install the gem:

    gem install dist/puccini.gem

The Puccini shared library will be in the `dist/` subdirectory. To use it you can either copy it into
your operating system's standard library path, or else set the path before running Ruby, e.g.:

    LD_LIBRARY_PATH=$LD_LIBRARY_PATH:dist ruby ...

Also see: [Ruby examples](../../examples/ruby/).
