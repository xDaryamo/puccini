Java Example
============

This relies on the [Java wrapper](../../wrappers/java/), so make sure to install that first.

To build this example using [Maven](https://maven.apache.org/):

    mvn -f examples/java

To run this example using Maven we need to make sure that the JVM process can load the shared
libraries:

    LD_LIBRARY_PATH=$LD_LIBRARY_PATH:dist mvn -f examples/java -q exec:java -Dexec.args=examples/tosca/data-types.yaml
