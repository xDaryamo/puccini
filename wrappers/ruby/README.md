Ruby Wrapper for Puccini
========================

This is a Ruby library for calling Puccini. It works by using
[Fiddle](https://ruby-doc.org/stdlib-2.0.0/libdoc/fiddle/rdoc/Fiddle.html) to call a shared library 
(.so) built from Puccini's Go code. This is done in-process, so there's no sub-process forking.

Note that we currently only support installation on 64-bit Linux.

To build the library and gem:

    scripts/build-library.sh
    cp dist/libpuccini.so wrappers/ruby/lib/
    gem build wrappers/ruby/puccini.gemspec -C wrappers/ruby --output ../../dist/puccini.gem

To install the gem:

    gem install dist/puccini.gem

Also see: [Ruby examples](../../examples/ruby/).
