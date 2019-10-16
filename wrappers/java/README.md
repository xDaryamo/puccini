Java Wrapper for Puccini
========================

This is a Java library for calling Puccini. It works by using a JNI shim library to call a shared
library (.so) built from Puccini's Go code. This is done in-process, so there's no sub-process
forking.

Note that we currently only support installation on 64-bit Linux.

The build requirements are [Maven](https://maven.apache.org/), gcc, and a full JDK. To install them
on Fedora:

    sudo dnf install mvn gcc java-11-openjdk-devel

To build the libraries:

    scripts/build-library.sh
    mvn -f wrappers/java

The Puccini shared library as well as the JNI shim shared library will both be in the `dist/`
subdirectory. To use them you can either copy them into your operating system's standard library
path, or else set the path before running the JVM, e.g.:

    LD_LIBRARY_PATH=$LD_LIBRARY_PATH:dist java ...

Also see: [Java examples](../../examples/java/).
