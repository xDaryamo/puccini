Java Example
============

This relies on the [Java wrapper](../../wrappers/java/), so make sure to install that first.

To build this example using [Maven](https://maven.apache.org/):

    mvn --file examples/java

To run this example using Maven we need to make sure that the JVM process can load the shared
libraries:

    LD_LIBRARY_PATH=$LD_LIBRARY_PATH:dist mvn --quiet --file examples/java exec:java --define exec.args=examples/1.3/data-types.yaml
