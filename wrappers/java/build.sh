#!/bin/bash
set -e

HERE=$(dirname "$(readlink -f "$0")")

JAVA_HOME=/usr/lib/jvm/jre-11-openjdk

javah puccini.TOSCA
gcc -fPIC -c puccini_TOSCA.c -I "$JAVA_HOME/include" -I "$JAVA_HOME/include/linux" -I "$HERE/../dist"
gcc puccini_TOSCA.o -shared -o libpuccinijni.so -Wl,-lpuccini -L "$HERE/../dist"
javac $(find -name \*.java)
LD_LIBRARY_PATH=".:$HERE/../dist" java puccini.Compile
