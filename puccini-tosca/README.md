puccini-tosca
=============

`compile`
---------

The most common command is `compile`. The optional input is a filesystem path or URL to a TOSCA
service template YAML file or a CSAR file. If no input is provided will attempt to read YAML
from stdin. By default the compiled Clout will be output to stdout, but you can use the
`--output/-o` switch to specify a file (or direct to a file in the shell via `>`).

If your TOSCA service template YAML imports other units, and these imported paths are relative,
then Puccini will assume that they are relative to the base URL of the main service template file.
This allows you to move your files around, and even host them at a URL, without changing the file
contents. This works as expected even within a CSAR file because Puccini uses `zip:` URLs to locate
files within it.

Also useful (or necessary) is the `--input/-i` switch, which lets you set the inputs for the TOSCA
topology template. You may specify this switch multiple times for as many inputs as necessary.
The format is `name=value`, where `value` is JSON-encoded. For example, to set an integer:

    --input cores=4

To set a string:

    --input ram=4gb

To set a complex data type:

    --input 'port={"type":"tcp","number":8080}'

If TOSCA parsing/compilation fails will emit a colorful problem report to stderr and exit with 1.
Thus you can use the `--quiet/-q` switch if all you want to do is check for successful compilation.

The default format for input/output is YAML, but you can switch to JSON using `--format/-f`. Note
that Clout in JSON may lose some type information (e.g. JSON doesn't distinguish between an integer
and a float).

Note that TOSCA functions and constraints are not called during compilation. They are embedded in
the Clout and are intended to be executed when necessary. To validate them you can use the
`tosca.coerce` JavaScript with **puccini-js**:

    puccini-tosca compile tosca.yaml | puccini-js exec tosca.coerce

Alternatively, you can use the `parse` command (see below).

`parse`
-------

If you need more diagnostics for TOSCA parsing use the `parse` command. It works similarly to
`compile` but does not emit Clout. Instead, it provides you various switches for examining the
internal working of Puccini's TOSCA parser.

Use `--stop/-s` to specify a [phase](../tosca/parser/README.md) (1-6) at which you wish the parser
to stop. This could be useful if you're getting too many problems in your report and wish to
minimize them to a more manageable list. (`-s 0` will skip the parser entirely and just check that
the input is readable.)

`--print/-p` is used to print out the results of each phase. You may specify multiple phases to
print out using ",", e.g. `-p 2,3,4`. Per phase you will see:

* Phase 1: Read. The hierarchy of imported units starting at the service template URL.
* Phase 2: Namespaces. All names per type. Each imported unit has its own namespace.
* Phase 3: Hierarchies. Tree of all tapes by fully qualified name. Each imported unit has its own
  type hierarchy.
* Phase 4: Inheritance. A tree of all inheritance tasks and their dependencies by path.  
* Phase 5: Rendering. Dumps the rendered entities.
  More useful, perhaps, would be the `--examine/-e` switch (see below).
* Phase 6: Topology. Tree representation of node template relationships.

The `--examine/-e` switch can be used to dump a specific parsed entity. Each entity is given a
path that more-or-less follows JSON. For example, a path can be:

    topology_template.node_templates['store'].properties['name']

The switch will search for all paths that contains your string, e.g. `-e properties`. You can even
include one or more wildcards, e.g. `-e 'node\*properties\*data`.

One more useful switch is `--coerce/-c` which works similarly to `puccini-js exec tosca.coerce`.
This is provided as a convenience to allow **puccini-tosca** to be more self-contained. Note that it
only tests that coercion is successful and does not emit the coerced values. For that, use
**puccini-js**.
